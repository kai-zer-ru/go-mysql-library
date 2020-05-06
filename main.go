package GoMySqlLibrary

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/kaizer666/gologger"
)

type MySQL struct {
	Address, DbName, User, Password, Charset string
	isConnected                              bool
	conn                                     *sql.DB
	logger                                   *gologger.Logger
	loggerConnected                          bool
	nextQuerySkipped                         bool
	debugMode                                bool
	telegram                                 telegram
	telegramConnected                        bool
}

func (v *MySQL) EnableDebug() {
	v.debugMode = true
}

func (v *MySQL) EnableTelegram(botToken string, channel int64) {
	v.telegramConnected = true
	v.telegram = telegram{
		botToken: botToken,
		channel:  channel,
	}
}

func (v *MySQL) Ping() error {
	return v.conn.Ping()
}

func (v *MySQL) Connect() error {
	if v.Charset == "" {
		v.Charset = "utf8mb4"
	}
	address := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", v.User, v.Password, v.Address, v.DbName, v.Charset)
	conn, err := sql.Open("mysql", address)
	if err != nil {
		return err
	}
	v.conn = conn
	v.isConnected = true
	return nil
}

func (v *MySQL) Close() error {
	v.isConnected = false
	return v.conn.Close()
}

func (v *MySQL) AddLogger(logger *gologger.Logger) {
	v.logger = logger
	v.loggerConnected = true
}

func (v *MySQL) GetOneField(query, field string, params ...interface{}) (interface{}, error) {
	data, err := v.GetArray(query, params...)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, sql.ErrNoRows
	}
	d, ok := data[0][field]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return d, nil
}

func (v *MySQL) GetOne(query string, params ...interface{}) (map[string]interface{}, error) {
	data, err := v.GetArray(query, params...)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, sql.ErrNoRows
	}
	return data[0], nil
}

func getRecord(name string, col interface{}, record *map[string]interface{}) {
	rec := *record
	switch col.(type) {
	case int:
		rec[name] = int64(col.(int))
	case int8:
		rec[name] = int64(col.(int8))
	case int16:
		rec[name] = int64(col.(int16))
	case int32:
		rec[name] = int64(col.(int32))
	case int64:
		rec[name] = col.(int64)
	case uint:
		rec[name] = int64(col.(uint))
	case uint8:
		rec[name] = int64(col.(uint8))
	case uint16:
		rec[name] = int64(col.(uint16))
	case uint32:
		rec[name] = int64(col.(uint32))
	case uint64:
		rec[name] = int64(col.(uint64))
	case float32:
		rec[name] = float64(col.(float32))
	case float64:
		rec[name] = col.(float64)
	default:
		rec[name] = string(col.([]uint8))
	}
	record = &rec
}

func (v *MySQL) getRecords(query string, params ...interface{}) ([]map[string]interface{}, error) {
	frame := Caller(3)
	v.SaveQueryToLog(frame, query, params...)
	rows, err := v.conn.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	records := make([]map[string]interface{}, 0)
	columns, _ := rows.Columns()
	ct, _ := rows.ColumnTypes()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		record := make(map[string]interface{})
		for i, col := range values {
			if col != "" {
				colType := ct[i]
				switch colType.ScanType().String() {
				case "sql.NullInt64":
					if col == nil {
						record[columns[i]] = int64(0)
					} else {
						switch col.(type) {
						case int64:
							record[columns[i]] = col.(int64)
						case int32:
							record[columns[i]] = int64(col.(int32))
						default:
							intCol, _ := strconv.Atoi(string(col.([]uint8)))
							record[columns[i]] = int64(intCol)
						}
					}
				case "sql.NullFloat64":
					if col == nil {
						record[columns[i]] = float64(0.0)
					} else {
						switch col.(type) {
						case float64:
							record[columns[i]] = col.(float64)
						case float32:
							record[columns[i]] = float64(col.(float32))
						default:
							colFloat, _ := strconv.ParseFloat(string(col.([]uint8)), 10)
							record[columns[i]] = colFloat
						}
					}
				case "mysql.NullTime":
					fallthrough
				case "sql.RawBytes":
					if col == nil {
						record[columns[i]] = ""
					} else {
						record[columns[i]] = string(col.([]uint8))
					}
				default:
					if v.loggerConnected && v.debugMode {
						v.logger.Debug("type of collumn %s = %v (%T)", columns[i], colType.ScanType().String(), col)
					}
					getRecord(columns[i], col, &record)
				}
			}
		}
		records = append(records, record)
	}
	return records, nil
}

func (v *MySQL) Call(query string, params ...interface{}) (map[string]interface{}, error) {
	records, err := v.getRecords(query, params...)
	if err != nil {
		return nil, err
	}
	return records[0], nil
}

func (v *MySQL) GetArray(query string, params ...interface{}) ([]map[string]interface{}, error) {
	records, err := v.getRecords(query, params...)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (v *MySQL) Execute(query string, params ...interface{}) (sql.Result, error) {
	frame := Caller(2)
	v.SaveQueryToLog(frame, query, params...)
	return v.conn.Exec(query, params...)
}

func (v *MySQL) SaveQueryToLog(frame runtime.Frame, query string, params ...interface{}) {
	src := fmt.Sprintf("%s:%d (%s)", filepath.Base(frame.File), frame.Line, frame.Function)
	if v.loggerConnected {
		if v.nextQuerySkipped {
			v.nextQuerySkipped = false
			return
		}
		for _, param := range params {
			query = strings.Replace(query, "?", fmt.Sprintf("%v", param), 1)
		}
		v.logger.Debug(src)
		v.logger.Debug(query)
	}
}

func Caller(skip int) (frame runtime.Frame) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(skip+2, rpc[:])
	if n < 1 {
		return
	}
	frame, _ = runtime.CallersFrames(rpc).Next()
	return frame
}

func (v *MySQL) SkipNextQueryLog() *MySQL {
	v.nextQuerySkipped = true
	return v
}

func (v *MySQL) NotifyDbDisconnected(err error) {
	if v.loggerConnected {
		v.logger.Error("mysql ping error: %v", err)
	}
	if v.telegramConnected {
		v.telegram.message = fmt.Sprintf("Erorr ping db: %v", err)
		err = v.telegram.SendMessageToTelegram()
		if err != nil {
			if v.loggerConnected {
				v.logger.Error("error SendMessageToTelegram: %v", err)
			}
		}
	}
}

func (v *MySQL) StartPing() {
	for {
		<-time.After(time.Minute)
		err := v.Ping()
		if err != nil {
			v.NotifyDbDisconnected(err)
		}
	}
}
