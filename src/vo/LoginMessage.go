package vo

// LoginMessage : 登录信息
type LoginMessage struct {
	UserName    string
	LatestTime  string
	AccessToken string
	FileSize    int64
	Status      int8
	Role        string
	HeadPhoto   string
}
