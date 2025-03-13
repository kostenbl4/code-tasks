package domain

type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Hpass string `json:"-"`
}
