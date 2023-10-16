package domain

type User struct {
	Id       int64
	Email    string
	Password string
}

//func (u User) ValidateEmail() bool {
//	return u.Email
//}
