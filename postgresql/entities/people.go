package entities

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type EntityPeople struct {
	id         int
	name       string
	surname    string
	patronymic string
}

type EntityPeopleValidation struct{}

func (people *EntityPeople) GetId() int {
	return people.id
}

func (people *EntityPeople) GetName() string {
	return people.name
}

func (people *EntityPeople) GetSurname() string {
	return people.surname
}

func (people *EntityPeople) GetPatronymic() string {
	return people.patronymic
}

func (*EntityPeople) ValidateId(value int) bool {
	if value > 0 {
		return true
	} else {
		return false
	}
}

func (*EntityPeople) ValidateName(value string) bool {
	const length int = 64
	if len(value) > length {
		return false
	} else {
		return true
	}
}

func (*EntityPeople) ValidateSurname(value string) bool {
	const length int = 64
	if len(value) > length {
		return false
	} else {
		return true
	}
}

func (*EntityPeople) ValidatePatronymic(value string) bool {
	const length int = 64
	if len(value) > length {
		return false
	} else {
		return true
	}
}

func (people *EntityPeople) Add(
	ctx context.Context,
	conn *pgx.Conn,
	name string,
	surname string,
	patronymic string,
) int {
	var id int
	err := conn.QueryRow(
		ctx,
		"INSERT INTO people VALUES (DEFAULT, $1, $2, $3) RETURNING id",
		name,
		surname,
		patronymic,
	).Scan(&id)
	if err != nil {
		logrus.Debugln(err)
		return 0
	}

	return id
}
