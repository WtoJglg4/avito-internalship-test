package entities

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" bindig:"required" db:"login"`
	Password string `json:"password" bindig:"required" db:"password_hash"`
	Role     string `json:"-" db:"role"`
}
