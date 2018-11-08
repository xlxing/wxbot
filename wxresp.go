package wxbot

type LoginRedirect struct {
	Skey       string `xml:"skey"`
	Wxsid      string `xml:"wxsid"`
	Wxuin      string `xml:"wxuin"`
	Passticket string `xml:"pass_ticket"`
}

type BaseResponse struct {
	Ret    int32  `json:"Ret"`
	ErrMsg string `json:"ErrMsg"`
}

type Contact struct {
	Uin              int64    `json:"Uin"`
	UserName         string   `json:"UserName"`
	NickName         string   `json:"NickName"`
	HeadImgUrl       string   `json:"HeadImgUrl"`
	ContactFlag      int      `json:"ContactFlag"`
	MemberCount      int      `json:"MemberCount"`
	MemberList       []Member `json:"MemberList"`
	RemarkName       string   `json:"RemarkName"`
	HideInputBarFlag int      `json:"HideInputBarFlag"`
	Sex              int      `json:"Sex"`
	Signatrue        string   `json:"Signatrue"`
	VerifyFlag       int      `json:"VerifyFlag"`
	OwnerUin         int      `json:"OwnerUin"`
	PYInitial        string   `json:"PYInitial"`
	PYQuanPin        string   `json:"PYQuanPin"`
	RemarkPYInitial  string   `json:"RemarkPYInitial"`
	RemarkPYQuanPin  string   `json:"RemarkPYQuanPin"`
	StarFriend       int      `json:"StarFriend"`
	AppAccountFlag   int      `json:"AppAccountFlag"`
	Statues          int      `json:"Statues"`
	AttrStatus       int      `json:"AttrStatus"`
	Province         string   `json:"Province"`
	City             string   `json:"City"`
	Alias            string   `json:"Alias"`
	SnsFlag          int      `json:"SnsFlag"`
	UinFriend        int      `json:"UinFriend"`
	DisplayName      string   `json:"DisplayName"`
	ChatRoomId       int32    `json:"ChatRoomId"`
	KeyWord          string   `json:"KeyWord"`
	EncryChatRoomId  string   `json:"EncryChatRoomId"`
	IsOwner          int      `json:"IsOwner"`
}

type User struct {
	Uin               int64  `json:"Uin"`
	UserName          string `json:"UserName"`
	NickName          string `json:"NickName"`
	HeadImgUrl        string `json:"HeadImgUrl"`
	RemarkName        string `json:"RemarkName"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	HideInputBarFlag  int    `json:"HideInputBarFlag"`
	StarFriend        int    `json:"StarFriend"`
	Sex               int    `json:"Sex"`
	Signatrue         string `json:"Signatrue"`
	AppAccountFlag    int    `json:"AppAccountFlag"`
	VerifyFlag        int    `json:"VerifyFlag"`
	ContactFlag       int    `json:"ContactFlag"`
	WebWxPluginSwitch int    `json:WebWxPluginSwitch`
	HeadImgFlag       int    `json:"HeadImgFlag"`
	SnsFlag           int    `json:"SnsFlag"`
}

type Member struct {
	Uin             int64  `json:"Uin"`
	UserName        string `json:"UserName"`
	NickName        string `json:"NickName"`
	AttrStatus      int    `json:"AttrStatus"`
	PYInitial       string `json:"PYInitial"`
	PYQuanPin       string `json:"PYQuanPin"`
	RemarkPYInitial string `json:"RemarkPYInitial"`
	RemarkPYQuanPin string `json:"RemarkPYQuanPin"`
	MemberStatus    int    `json:"MemberStatus"`
	DisplayName     string `json:"DisplayName"`
	KeyWord         string `json:"KeyWord"`
}

type SyncKey struct {
	Key int `json:"Key"`
	Val int `json:"Val"`
}

type SyncKeys struct {
	Count int       `json:"Count"`
	List  []SyncKey `json:"List"`
}

type MPArticle struct {
	Title  string `json:"Title"`
	Digset string `json:"Digset"`
	Cover  string `json:"Cover"`
	Url    string `json:"Url"`
}

type MPSubscribeMsg struct {
	UserName       string      `json:"UserName"`
	MPArticleCount int32       `json:"MPArticleCount"`
	MPArticleList  []MPArticle `json:"MPArticleList"`
	Time           uint64      `json:"Time"`
	NickName       string      `json:"NickName"`
}

type WXInit struct {
	BaseResponse        BaseResponse     `json:"BaseResponse"`
	Count               int              `json:"Count"`
	ContactList         []Contact        `json:"ContactList"`
	SyncKey             SyncKeys         `json:"SyncKey"`
	User                User             `json:"User"`
	ChatSet             string           `json:"ChatSet"`
	SKey                string           `json:"SKey"`
	ClientVersion       int64            `json:"ClientVersion"`
	SystemTime          int64            `json:"SystemTime"`
	GrayScale           int              `json:"GrayScale"`
	InviteStartCount    int              `json:"InviteStartCount"`
	MPSubscribeMsgCount int              `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList  []MPSubscribeMsg `json:"MPSubscribeMsgList"`
	ClickReportInterval int64            `json:"ClickReportInterval"`
}

type WXStatusNotify struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MsgID        string       `json:"MsgID"`
}

type WXContact struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int          `json:"MemberCount"`
	MemberList   []Contact    `json:"MemberList"`
}

type RecommendInfo struct {
	UserName   string `json:"UserName"`
	NickName   string `json:"NickName"`
	QQNum      int    `json:"QQNum"`
	Province   string `json:"Province"`
	City       string `json:"City"`
	Content    string `json:"Content"`
	Signatrue  string `json:"Signatrue"`
	Alias      string `json:"Alias"`
	Scene      int    `json:"Scene"`
	VerifyFlag int    `json:"VerifyFlag"`
	AttrStatus int    `json:"AttrStatus"`
	Sex        int    `json:"Sex"`
	Ticket     string `json:"Ticket"`
	OpCode     int    `json:"OpCode"`
}

type AppInfo struct {
	AppID string `json:"AppID"`
	Type  int    `json:"Type"`
}

type AddMsg struct {
	MsgId                string        `json:"MsgId"`
	FromUserName         string        `json:"FromUserName"`
	ToUserName           string        `json:"ToUserName"`
	MsgType              int           `json:"MsgType"`
	Content              string        `json:"Content"`
	Status               int           `json:"Status"`
	ImgStatus            int           `json:"ImgStatus"`
	CreateTime           int           `json:"CreateTime"`
	VoiceLength          int           `json:"VoiceLength"`
	PlayLength           int           `json:"PlayLength"`
	FileName             string        `json:"FileName"`
	FileSize             string        `json:"FileSize"`
	MediaId              string        `json:"MediaId"`
	Url                  string        `json:"Url"`
	AppMsgType           int           `json:"AppMsgType"`
	StatusNotifyCode     int           `json:"StatusNotifyCode"`
	StatusNotifyUserName string        `json:"StatusNotifyUserName"`
	RecommendInfo        RecommendInfo `json:"RecommendInfo"`
	ForwardFlag          int           `json:"ForwardFlag"`
	AppInfo              AppInfo       `json:"AppInfo"`
	HasProductId         int           `json:"HasProductId"`
	Ticket               string        `json:"Ticket"`
	ImgHeight            int           `json:"ImgHeight"`
	ImgWidth             int           `json:"ImgWidth"`
	SubMsgType           int           `json:"SubMsgType"`
	NewMsgId             int64         `json:"NewMsgId"`
	OriContent           string        `json:"OriContent"`
	EncryFileName        string        `json:"EncryFileName"`
}

type Profile struct {
	BitFlag  int `json:"BitFlag"`
	UserName struct {
		Buff string `json:"Buff"`
	} `json:"UserName"`
	NickName struct {
		Buff string `json:"Buff"`
	} `json:"NickName"`
	BindUin   int `json:"BindUin"`
	BindEmail struct {
		Buff string `json:"Buff"`
	} `json:"BindEmail"`
	BindMobile struct {
		Buff string `json:"BindMobile"`
	} `json:"BindMobile"`
	Status            int    `json:"Status"`
	Sex               int    `json:"Sex"`
	PersonalCard      int    `json:"PersonalCard"`
	Alias             string `json:"Alias"`
	HeadImgUpdateFlag int    `json:"HeadImgUpdateFlag"`
	HeadImgUrl        string `json:"HeadImgUrl"`
	Signatrue         string `json:"Signatrue"`
}

type WXSync struct {
	BaseResponse           BaseResponse  `json:"BaseResponse"`
	AddMsgCount            int           `json:"AddMsgCount"`
	AddMsgList             []AddMsg      `json:"AddMsgList"`
	ModContactCount        int           `json:"ModContactCount"`
	ModContactList         []Contact     `json:"ModContactList"`
	DelContactCount        int           `json:"DelContactCount"`
	DelContactList         []interface{} `json:"DelContactList"`
	ModChatRoomMemberCount int           `json:"ModChatRoomMemberCount"`
	ModChatRoomMemberList  []interface{} `json:"ModChatRoomMemberList"`
	Profile                Profile       `json:"Profile"`
	ContinueFlag           int           `json:"ContactFlag"`
	SyncKey                SyncKeys      `json:"SyncKey"`
	SKey                   string        `json:"SKey"`
	SyncCheckKey           SyncKeys      `json:"SyncCheckKey"`
}
