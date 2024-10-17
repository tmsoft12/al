package domain

type Admin struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password []byte `json:"."`
}
