package vo

// LoginMessage : 登录信息
type LoginMessage struct {
	UserName   string
	LatestTime string
	UserToken  string
	FileSize   int64
	Status     int8
}