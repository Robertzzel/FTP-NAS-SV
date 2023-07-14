package utils

type User struct {
	Name     string
	Password string
	loggedIn bool
}

func (user *User) LogIn() (bool, error) {
	return false, nil
}

func (user *User) IsLogenIn() bool {
	return user.loggedIn
}
