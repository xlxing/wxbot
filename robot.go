package wxbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sjson "github.com/bitly/go-simplejson"
)

type ReqPara struct {
	ReqType    int32 `json:"reqType"`
	Perception struct {
		InputText struct {
			Text string `json:"text"`
		} `json:"inputText"`
	} `json:"perception"`
	UserInfo struct {
		ApiKey string `json:"apiKey"`
		UserId string `json:"userId"`
	} `json:"userInfo"`
}

type Robot struct {
	Apikey string
	ApiUrl string
	rnet   *http.Client
}

func NewRobot(apikey string) *Robot {
	return &Robot{
		Apikey: apikey,
		ApiUrl: "http://openapi.tuling123.com/openapi/api/v2",
		rnet:   &http.Client{},
	}
}

func (r *Robot) GetReplyMsg(text string, userid string) string {
	var para ReqPara
	para.ReqType = 0
	para.Perception.InputText.Text = text
	para.UserInfo.ApiKey = r.Apikey
	para.UserInfo.UserId = userid
	jsondata, _ := json.Marshal(para)

	body, err := r.post(r.ApiUrl, jsondata)
	if err != nil {
		fmt.Println("get tuling reply err:", err)
		return ""
	}

	info, err := sjson.NewJson(body)
	if err != nil {
		fmt.Println("get tuling reply parse para err:", err)
		return ""
	}

	msg := info.Get("results").GetIndex(0).Get("values").Get("text").MustString("")
	return msg
}

func (r *Robot) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.rnet.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (r *Robot) post(url string, para []byte) ([]byte, error) {
	body := bytes.NewBuffer(para)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "appliction/json;charset=utf-8")

	resp, err := r.rnet.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
