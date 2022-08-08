package telegrambot

type User struct {
	Id string
}

func (user *User) Recipient() string {
	return user.Id
}
