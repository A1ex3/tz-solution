package people

import (
	"context"
	"net/http"
	"time"
	"tzsolution/postgresql"
	"tzsolution/postgresql/entities"
)

type People struct {
	Db *postgresql.Postgresql
}

type PeopleValidation struct{}

func NewPeople(db *postgresql.Postgresql) *People {
	return &People{
		Db: db,
	}
}

func (people *People) NewPeople(
	ctx context.Context,
	name string,
	surname string,
	patronymic string,
) (int, int) {
	entityPeople := entities.EntityPeople{}
	if !entityPeople.ValidateName(name) || !entityPeople.ValidateSurname(surname) || !entityPeople.ValidatePatronymic(patronymic) {
		return 0, http.StatusBadRequest
	}

	dbContext, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	dbConn, err := people.Db.Connect(dbContext)
	defer people.Db.Close(ctx, dbConn)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	result := entityPeople.Add(
		dbContext,
		dbConn,
		name,
		surname,
		patronymic,
	)
	if result > 0 {
		return result, http.StatusOK
	} else {
		return result, http.StatusInternalServerError
	}
}
