package testBuilders

import (
	"intern/models"
	"time"
)

type ActorBuilder struct {
	actor models.Actor
}

func NewActorBuilder() *ActorBuilder {
	return &ActorBuilder{}
}

func (b *ActorBuilder) WithID(id int) *ActorBuilder {
	b.actor.ID = id
	return b
}

func (b *ActorBuilder) WithFirstName(firstName string) *ActorBuilder {
	b.actor.FirstName = firstName
	return b
}

func (b *ActorBuilder) WithLastName(lastName string) *ActorBuilder {
	b.actor.LastName = lastName
	return b
}

func (b *ActorBuilder) WithGender(gender byte) *ActorBuilder {
	b.actor.Gender = gender
	return b
}

func (b *ActorBuilder) WithBirthday(birthday time.Time) *ActorBuilder {
	b.actor.Birthday = birthday
	return b
}

func (b *ActorBuilder) Build() models.Actor {
	return b.actor
}
