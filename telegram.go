package GoMySqlLibrary

import (
	`bytes`
	`encoding/json`
	`io/ioutil`
	`net`
	`net/http`
	`time`
)


type telegram struct {
	channel      int64
	botToken     string
	message      string
	rawData      []byte
	photoURL     string
	photoCaption string
}


func httpDial (network, addr string) (net.Conn, error) {
	return net.DialTimeout (network, addr, timeout)
}

var transport = http.Transport{Dial: httpDial, DisableKeepAlives: true }
var timeout = time.Duration(2) * time.Second
var httpClient = http.Client{Transport: &transport, Timeout: timeout}

func (t *telegram) SendMessageToTelegram() error {
	payloads := map[string]interface{}{
		"chat_id":    t.channel,
		"text":       t.message,
		"parse_mode": "Markdown",
	}
	// send message to telegram
	req, err2 := http.NewRequest("POST", "https://api.telegram.org/bot"+t.botToken+"/sendMessage", nil)
	if err2 != nil {
		return err2
	}
	req.Header.Add("Content-Type", `application/json`)
	b, marshalErr := json.Marshal(payloads)
	if marshalErr != nil {
		return marshalErr
	}
	req.Body = ioutil.NopCloser(bytes.NewBufferString(string(b)))
	if _, e := httpClient.Do(req); e != nil {
		return e
	}
	return nil
}
