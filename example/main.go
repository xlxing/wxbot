package main

import (
	"fmt"
	"io/ioutil"
	"os"

	sjson "github.com/bitly/go-simplejson"
	"github.com/xlxing/wxbot"
)

func main() {
	configFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	info, err := sjson.NewJson(data)
	if err != nil {
		fmt.Println("config parse err:", err)
		return
	}

	apiKey := info.Get("TuLing").Get("ApiKey").MustString("")
	if apiKey == "" {
		fmt.Println("config tuling apikey is nil")
		return
	}

	configFile.Close()

	wx, err := wxbot.NewWechat(apiKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	wx.Run()
}
