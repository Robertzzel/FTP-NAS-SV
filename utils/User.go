package utils

type User struct {
	Name     string
	Password string
	BasePath string
}

func (user *User) IsLogenIn() bool {
	return user.Name != "" && user.Password != ""
}
