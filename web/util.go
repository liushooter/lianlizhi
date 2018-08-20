package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

type errRes struct {
	// {"error_code":"TEMPLATE_NOT_FOUND","message":"Template not found"}
	// {"error_code":"TEMPLATE_AUDITING","message":"模板正在审核，无法使用"}
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type okRes struct {
	// {"message_ids":[{"message_id":8454155,"mobile":"14612345678"}]}

	MessageIds []struct {
		MessageID int    `json:"message_id"`
		Mobile    string `json:"mobile"`
	} `json:"message_ids"`
}

func GenSmsCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return vcode
}

func SendCodeSms(mobile string) (string, error) {
	url := "https://sms-api.upyun.com/api/messages"

	template_id := "2038"
	code := GenSmsCode()
	payload := fmt.Sprintf("{\"template_id\": %s, \"mobile\": %s, \"vars\": %s}", template_id, mobile, code)

	var jsonStr = []byte(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", viper.GetString("sms.auth"))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	errorCode := gjson.Get(string(body), "error_code").String()

	if "" == errorCode {
		return code, nil
	} else {
		err := errors.New(errorCode)
		return "", err
	}
}