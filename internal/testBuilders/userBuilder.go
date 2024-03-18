package testBuilders

import "intern/models"

type UserBuilder struct {
	user models.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (b *UserBuilder) WithID(id int) *UserBuilder {
	b.user.ID = id
	return b
}

func (b *UserBuilder) WithLogin(login string) *UserBuilder {
	b.user.Login = login
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.user.Password = password
	return b
}

func (b *UserBuilder) WithRole(role string) *UserBuilder {
	b.user.Role = role
	return b
}

func (b *UserBuilder) Build() models.User {
	return b.user
}
