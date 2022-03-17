package dao

type Users struct {
	User
	Role string
}

type User struct {
	UserName   string
	UserPwd    string
	SignupAt   string
	LastActive string
}

type Policies []Policy

type Policy struct {
	Sub string
	Obj string
	Act string
}
