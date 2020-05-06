package GoMySqlLibrary

import (
	"database/sql"
	"os"
	"testing"

	"gitlab.sarafann.com/kaizer/gologger"
)

func init() {
	_ = os.Remove("main.log")
}
func connectDb() (MySQL, error) {
	db := MySQL{
		Address:  "192.168.0.1:3306",
		User:     "test",
		Password: "test",
		DbName:   "test_database",
	}
	fileName := "main.log"
	log := gologger.Logger{}
	log.SetLogFileName(fileName)
	_ = log.SetLogLevel(1)
	err := log.Init()
	if err != nil {
		return db, err
	}
	db.AddLogger(&log)
	err = db.Connect()
	if err != nil {
		return db, err
	}
	return db, err
}

func TestGetOne(t *testing.T) {
	db, err := connectDb()
	if err != nil {
		t.Error(err)
	}
	_, err = db.GetOne("SELECT NOW()")
	if err != nil {
		t.Error(err)
	}
	_, err = db.GetOne("SELECT * FROM config")
	if err == nil {
		t.Errorf("no error, but table not exist")
	}
	_, err = db.GetOne("SELECT * FROM test_empty")
	if err == nil {
		t.Errorf("no error, but table is empty")
	} else if err != sql.ErrNoRows {
		t.Errorf("no error, but table is empty 2")
	}

	_, err = db.GetOne("SELECT * FROM test_one")
	if err != nil {
		t.Error(err)
	}

	_, err = db.GetOne("SELECT * FROM test_multi")
	if err != nil {
		t.Error(err)
	}
}

func TestConnect(t *testing.T) {
	db, err := connectDb()
	if err != nil {
		t.Error(err)
	}

	_, err = db.GetOneField("SELECT field FROM test_fireld", "field")
	if err != nil {
		t.Error(err)
	}
	_, err = db.GetArray("SELECT NOW()")
	if err != nil {
		t.Error(err)
	}
	_, err = db.Execute("SELECT NOW()")
	if err != nil {
		t.Error(err)
	}
	_, err = db.Call("SELECT NOW()")
	if err != nil {
		t.Error(err)
	}
	err = db.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestWithParams(t *testing.T) {
	db, err := connectDb()
	if err != nil {
		t.Error(err)
	}
	_, err = db.GetOne("SELECT * FROM test_one WHERE id=?", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPing(t *testing.T) {
	db, err := connectDb()
	if err != nil {
		t.Error(err)
	}
	err = db.Ping()
	if err != nil {
		t.Error(err)
	}
}

func TestBids(t *testing.T) {
	db, err := connectDb()
	if err != nil {
		t.Error(err)
	}
	row, err := db.GetOneField("SELECT SUM(money) as summ FROM sarafann_test.product_bids WHERE product_id=306516 and is_finished=0", "summ")
	if err != nil {
		t.Error(err)
	}
	switch row.(type) {
	case float64:
		if row.(float64) == 0.0 {
			t.Log("Fine")
		} else {
			t.Errorf("Error! row != 0.0 but row = %v", row)
		}
	default:
		t.Errorf("error! type of row = %T", row)
	}
}

func TestProduct(t *testing.T) {
	db, err := connectDb()
	if err != nil {
		t.Error(err)
	}
	row, err := db.GetOne("SELECT * FROM sarafann_test.products WHERE id = ?", 8015)
	if err != nil {
		t.Error(err)
	}
	switch row["district_id"].(type) {
	case int64:
		t.Log("Type success!")
	default:
		t.Errorf("error! type of row = %T (%v)", row["district_id"], row["district_id"])
	}
}
