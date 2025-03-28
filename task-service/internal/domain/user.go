package domain

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"login"`
	Hpass    string `json:"-"`
}
