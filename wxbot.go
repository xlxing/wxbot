package wxbot

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mdp/qrterminal"
)

const (
	WxUrl       = "https://login.weixin.qq.com"
	WxReferer   = "https://wx.qq.com/"
	WxUA        = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"
	WxAppid     = "wx782c26e4c19acffb"
	RedirectUri = "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxnewloginpage"
)

var Hosts = []string{
	"webpush.wx.qq.com",
	"webpush.weixin.qq.com",
	"webpush.wx2.qq.com",
	"webpush.wx8.qq.com",
	"webpush.web2.wechat.com",
	"webpush.web.wechat.com",
}

type Wechat struct {
	uuid          string
	redirectUrl   string
	loginRedirect LoginRedirect
	baseRequest   map[string]interface{}
	deviceID      string
	wxnet         *http.Client
	user          User
	synckeylist   SyncKeys
	contact       map[string]Contact
	robot         *Robot
}

func NewWechat(apiKey string) (*Wechat, error) {
	// 随机一串数字作为设备id
	rand.Seed(time.Now().UnixNano())
	id := "e" + strconv.Itoa(rand.Int())[0:15]

	jar, _ := cookiejar.New(nil)

	return &Wechat{
		wxnet:    &http.Client{Jar: jar},
		deviceID: id,
		robot:    NewRobot(apiKey),
	}, nil
}

func (wx *Wechat) Run() {
	fmt.Println("获取UUID")
	if err := wx.GetUUid(); err != nil {
		return
	}
	fmt.Println("获取二维码")
	if err := wx.GetLoginQrcode(); err != nil {
		return
	}
	fmt.Println("请扫描二维码")
	if err := wx.ScanQrCodeAndLogin(); err != nil {
		return
	}
	wx.Redirect()
	if err := wx.WxInit(); err != nil {
		return
	}
	wx.WxStatusNotify()
	wx.WxGetContact()
	wx.Loop()
}

func (wx *Wechat) GetUUid() error {
	if wx.uuid != "" {
		return nil
	}
	uri := WxUrl + "/jslogin?appid=" + WxAppid + "&redirect_uri" + RedirectUri + "&fun=new&lang=zh_CN&_=" + strconv.FormatInt(time.Now().Unix(), 10)
	data, err := wx.get(uri)
	if err != nil {
		fmt.Println("read uuid resp err:", err)
		return err
	}

	datas := strings.Split(string(data), ";")
	loginData := make(map[string]string)
	for _, v := range datas {
		kv := strings.Split(v, " = ")
		if len(kv) == 2 {
			loginData[strings.Replace(kv[0], " ", "", -1)] = strings.Replace(kv[1], "\"", "", -1)
		}
	}

	if loginData["window.QRLogin.code"] == "200" {
		wx.uuid = loginData["window.QRLogin.uuid"]
		return nil
	}
	err = fmt.Errorf("get uuid parse resp data err:%s", string(data))
	fmt.Println(err)
	return err
}

func (wx *Wechat) GetLoginQrcode() error {
	uri := WxUrl + "/l/" + wx.uuid
	qrterminal.GenerateHalfBlock(uri, qrterminal.L, os.Stdout)
	return nil
}

func (wx *Wechat) ScanQrCodeAndLogin() error {
	tip := "1"
	for {
		uri := WxUrl + "/cgi-bin/mmwebwx-bin/login?tip=" + tip + "&uuid=" + wx.uuid + "&_=" + strconv.FormatInt(time.Now().Unix(), 10)
		body, err := wx.get(uri)
		if err != nil {
			fmt.Println("scan qrcode err: ", err)
			return err
		}
		data := string(body)
		if strings.Contains(data, "window.code=200") {
			ss := strings.Split(data, "\"")
			if len(ss) < 2 {
				return fmt.Errorf("scan qrcode login parse err")
			}
			wx.redirectUrl = ss[1] + "&fun=new&version=v2&lang=zh_CN"
			fmt.Println("登录成功")
			return nil
		} else if strings.Contains(data, "window.code=201") {
			fmt.Println("扫描成功，请确认登录")
			tip = "0"
		}
	}
}

func (wx *Wechat) Redirect() error {
	body, err := wx.get(wx.redirectUrl)
	if err != nil {
		fmt.Println("redirect err:", err)
		return err
	}
	var data LoginRedirect
	if err := xml.Unmarshal(body, &data); err != nil {
		fmt.Println("redirect parse err:", err)
		return err
	}
	wx.loginRedirect = data
	uin, _ := strconv.ParseFloat(data.Wxuin, 64)
	wx.baseRequest = make(map[string]interface{})
	wx.baseRequest["Uin"] = int64(uin)
	wx.baseRequest["Sid"] = data.Wxsid
	wx.baseRequest["Skey"] = data.Skey
	wx.baseRequest["DeviceID"] = wx.deviceID
	return nil
}

func (wx *Wechat) WxInit() error {
	uri := fmt.Sprintf("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit?pass_ticket=%s&skey=%s&r=%s", wx.loginRedirect.Passticket, wx.loginRedirect.Skey, strconv.FormatInt(time.Now().Unix(), 10))
	// uri := fmt.Sprintf("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxinit?r=%s&pass_ticket=%s", strconv.FormatInt(time.Now().Unix(), 10), wx.loginRedirect.Passticket)
	para := make(map[string]interface{})
	para["BaseRequest"] = wx.baseRequest
	body, err := wx.post(uri, para)
	if err != nil {
		fmt.Println("wxinit err:", err)
		return err
	}

	var wxinitdata WXInit
	if err := json.Unmarshal(body, &wxinitdata); err != nil {
		fmt.Println("wxinit parse err:", err)
		return err
	}

	if wxinitdata.BaseResponse.Ret != 0 {
		err = fmt.Errorf("wxinitdata BaseResponse err", wxinitdata.BaseResponse.ErrMsg)
		return err
	}

	wx.user = wxinitdata.User
	wx.synckeylist = wxinitdata.SyncKey
	return nil
}

func (wx *Wechat) WxStatusNotify() error {
	uri := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxstatusnotify?lang=zh_CN&pass_ticket=" + wx.loginRedirect.Passticket
	para := make(map[string]interface{})
	para["BaseRequest"] = wx.baseRequest
	para["Code"] = 3
	para["FromUserName"] = wx.user.UserName
	para["ToUserName"] = wx.user.UserName
	para["ClientMsgId"] = int(time.Now().Unix())

	body, err := wx.post(uri, para)
	if err != nil {
		fmt.Println("wxstatusnotify err:", err)
		return err
	}

	var data WXStatusNotify
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("wxstatusnotify parse err:", err)
		return err
	}

	if data.BaseResponse.Ret != 0 {
		err = fmt.Errorf("wxstatusnotify BaseResponse err", data.BaseResponse.ErrMsg)
		return err
	}

	return nil
}

func (wx *Wechat) GetStrSyncKey() string {
	strkey := []string{}
	for _, v := range wx.synckeylist.List {
		strkey = append(strkey, strconv.Itoa(v.Key)+"_"+strconv.Itoa(v.Val))
	}
	return strings.Join(strkey, "|")
}

func (wx *Wechat) WxSyncCheck() (int, int) {
	re := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	for _, host := range Hosts {
		uri := fmt.Sprintf("https://%s/cgi-bin/mmwebwx-bin/synccheck", host)
		para := url.Values{}
		para.Add("r", strconv.FormatInt(time.Now().Unix(), 10))
		para.Add("skey", wx.loginRedirect.Skey)
		para.Add("sid", wx.baseRequest["Sid"].(string))
		para.Add("uin", strconv.FormatInt(wx.baseRequest["Uin"].(int64), 10))
		para.Add("deviceid", wx.deviceID)
		para.Add("synckey", wx.GetStrSyncKey())
		para.Add("_", strconv.FormatInt(time.Now().Unix(), 10))
		uri = uri + "?" + para.Encode()

		body, err := wx.get(uri)
		if err != nil {
			continue
		}
		codes := re.FindStringSubmatch(string(body))
		if len(codes) > 2 {
			retcode, _ := strconv.Atoi(codes[1])
			selector, _ := strconv.Atoi(codes[2])
			if retcode == 0 || retcode == 1100 {
				return retcode, selector
			}
		}
	}
	return -1, 0
}

func (wx *Wechat) WxSync() (*WXSync, error) {
	uri := fmt.Sprintf("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync?sid=%s&skey=%s&pass_ticket=%s", wx.baseRequest["Sid"].(string), wx.loginRedirect.Skey, wx.loginRedirect.Passticket)
	para := make(map[string]interface{})
	para["BaseRequest"] = wx.baseRequest
	para["SyncKey"] = wx.synckeylist
	para["rr"] = ^int(time.Now().Unix())
	body, err := wx.post(uri, para)
	if err != nil {
		return nil, err
	}
	var data WXSync
	err = json.Unmarshal(body, &data)
	if data.BaseResponse.Ret != 0 {
		return nil, fmt.Errorf("wxsync err")
	}
	return &data, nil
}

func (wx *Wechat) SendMsg(content string, userName string) bool {
	uri := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsendmsg?pass_ticket=" + wx.loginRedirect.Passticket
	para := make(map[string]interface{})
	para["BaseRequest"] = wx.baseRequest
	msg := make(map[string]interface{})
	msg["Type"] = 1
	msg["Content"] = content
	msg["FromUserName"] = wx.user.UserName
	msg["ToUserName"] = userName
	msgId := strconv.FormatInt(time.Now().Unix()<<4, 10) + strconv.Itoa(rand.Int())[0:4]
	msg["ClientMsgId"] = msgId
	msg["LocallD"] = msgId
	para["Msg"] = msg
	_, err := wx.post(uri, para)
	if err != nil {
		return false
	}
	return true
}

func (wx *Wechat) WxParseNewMsg(data *WXSync) bool {
	for _, contact := range data.ModContactList {
		if _, ok := wx.contact[contact.UserName]; !ok {
			wx.contact[contact.UserName] = contact
		}
	}
	for _, msg := range data.AddMsgList {
		if msg.MsgType == 1 {
			if !(msg.FromUserName[:2] == "@@") {
				/*
					// test
					if msg.FromUserName == wx.user.UserName {
						name := wx.contact[msg.FromUserName].NickName
						if wx.contact[msg.FromUserName].RemarkName != "" {
							name = wx.contact[msg.FromUserName].RemarkName
						}
						reply := wx.robot.GetReplyMsg(msg.Content, "robot")
						fmt.Printf("[%s]:%s\n", name, msg.Content)
						if reply != "" {
							wx.SendMsg(reply, msg.FromUserName)
							fmt.Printf("[%s]:%s\n", wx.user.NickName, reply)
						}
					}
				*/
				name := wx.contact[msg.FromUserName].NickName
				if wx.contact[msg.FromUserName].RemarkName != "" {
					name = wx.contact[msg.FromUserName].RemarkName
				}
				reply := wx.robot.GetReplyMsg(msg.Content, "robot")
				fmt.Printf("[%s]:%s\n", name, msg.Content)
				if reply != "" {
					wx.SendMsg(reply, msg.FromUserName)
					fmt.Printf("[%s]:%s\n", wx.user.NickName, reply)
				}
			}
		}
	}
	return false
}

func (wx *Wechat) Loop() {
	time30Msec := time.NewTicker(time.Microsecond * 30)
	failCount := 0
	for {
		select {
		case <-time30Msec.C:
			retcode, selector := wx.WxSyncCheck()
			switch retcode {
			case 0:
				switch selector {
				case 0:
				case 2:
					data, err := wx.WxSync()
					if err == nil {
						if data.SyncKey.Count > 0 {
							wx.synckeylist = data.SyncKey
						}
						wx.WxParseNewMsg(data)
					}
				case 7:
				}
			case 1100:
				fmt.Println("账号在别的地方登陆，服务关闭")
				return
			case -1:
				failCount++
				if failCount >= 10 {
					fmt.Println("超过最大失败次数，自动退出")
					return
				}
			}
		}
	}
}

func (wx *Wechat) WxGetContact() error {
	uri := fmt.Sprintf("https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact?lang=zh_CN&pass_ticket=%s&seq=0&skey=%s&r=%s", wx.loginRedirect.Passticket, wx.loginRedirect.Skey, strconv.FormatInt(time.Now().Unix(), 10))
	body, err := wx.get(uri)
	if err != nil {
		fmt.Println("wxgetcontact err:", err)
		return err
	}

	var data WXContact
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("wxgetcontact parse err:", err)
		return err
	}

	if data.BaseResponse.Ret != 0 {
		err = fmt.Errorf("wxgetcontact BaseResponse err", data.BaseResponse.ErrMsg)
		return err
	}
	wx.contact = make(map[string]Contact)
	for _, contact := range data.MemberList {
		wx.contact[contact.UserName] = contact
	}
	fmt.Println("获取到用户：", data.MemberCount)
	return nil
}

func (wx *Wechat) get(uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	// fake req header
	req.Header.Add("Referer", WxReferer)
	req.Header.Add("User-agent", WxUA)

	resp, err := wx.wxnet.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (wx *Wechat) post(uri string, para map[string]interface{}) ([]byte, error) {
	data, err := json.Marshal(para)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "appliction/json;charset=utf-8")
	// fake req header
	req.Header.Add("Referer", WxReferer)
	req.Header.Add("User-agent", WxUA)

	resp, err := wx.wxnet.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
