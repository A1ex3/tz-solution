package entities

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type EntityCar struct {
	id     int
	year   int
	owner  *EntityPeople
	regNum string
	mark   string
	model  string
}

func (car *EntityCar) GetId() int {
	return car.id
}

func (car *EntityCar) GetOwner() *EntityPeople {
	return car.owner
}

func (car *EntityCar) GetRegNum() string {
	return car.regNum
}

func (car *EntityCar) GetMark() string {
	return car.mark
}

func (car *EntityCar) GetModel() string {
	return car.model
}

func (car *EntityCar) GetYear() int {
	return car.year
}

func (*EntityCar) ValidateId(value int) bool {
	if value > 0 {
		return true
	} else {
		return false
	}
}

func (*EntityCar) ValidateRegNum(value string) bool {
	const length int = 20
	if len(value) > length {
		return false
	}

	pattern := `^[A-Z]{1,3}\d{1,3}[A-Z]{1,3}$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		logrus.Infoln("An error while compiling a regular expression:", err)
		return false
	}
	return matched
}

func (*EntityCar) ValidateMark(value string) bool {
	const length int = 255

	if len(value) > length {
		return false
	} else {
		return true
	}
}

func (*EntityCar) ValidateModel(value string) bool {
	const length int = 255

	if len(value) > length {
		return false
	} else {
		return true
	}
}

func (*EntityCar) ValidateYear(value int) bool {
	currentYear := time.Now().Year()
	minYear := 1950

	if value < minYear || value > currentYear {
		return false
	}

	return true
}

func (car *EntityCar) Add(
	ctx context.Context,
	conn *pgx.Conn,
	ownerId int,
	regNum string,
	mark string,
	model string,
	year int,
) int {
	var id int
	err := conn.QueryRow(
		ctx,
		"INSERT INTO cars VALUES (DEFAULT, $1, $2, $3, $4, $5) RETURNING id",
		ownerId,
		regNum,
		mark,
		model,
		year,
	).Scan(&id)
	if err != nil {
		logrus.Debugln(err)
		return 0
	}

	return id
}

func (car *EntityCar) DeleteByCarId(
	ctx context.Context,
	conn *pgx.Conn,
	carId int,
) int {
	result, err := conn.Exec(
		ctx,
		"DELETE FROM cars WHERE id=$1 RETURNING id",
		carId,
	)
	if err != nil {
		logrus.Debugln(err)
		return 0
	}

	return int(result.RowsAffected())
}

func (car *EntityCar) UpdateByCarId(
	ctx context.Context,
	conn *pgx.Conn,
	id int,
	fields map[string]interface{},
) (int, int) {
	var args []interface{}
	setClause := strings.Builder{}

	for key, value := range fields {
		switch key {
		case "ownerId":
			if ownerId, ok := value.(int); ok {
				if !car.owner.ValidateId(ownerId) {
					return 0, http.StatusBadRequest
				}
				setClause.WriteString(", owner = $")
				setClause.WriteString(strconv.Itoa(len(args) + 1))
				args = append(args, ownerId)
			}
		case "regNum":
			if regNum, ok := value.(string); ok {
				if !car.ValidateRegNum(regNum) {
					return 0, http.StatusBadRequest
				}

				setClause.WriteString(", regnum = $")
				setClause.WriteString(strconv.Itoa(len(args) + 1))
				args = append(args, regNum)
			}
		}
	}

	if setClause.Len() == 0 {
		logrus.Debugln("No fields provided for update")
		return 0, http.StatusBadRequest
	}

	query := "UPDATE cars SET " + setClause.String()[2:] + " WHERE id=$" + strconv.Itoa(len(args)+1) + " RETURNING id"

	args = append(args, id)
	result, err := conn.Exec(
		ctx,
		query,
		args...,
	)

	if err != nil {
		logrus.Debugln(err)
		return 0, http.StatusInternalServerError
	}

	return int(result.RowsAffected()), http.StatusOK
}

func (car *EntityCar) GetInfoByRegNum(
	ctx context.Context,
	conn *pgx.Conn,
	regNum string,
) *EntityCar {
	entityCar := &EntityCar{
		owner: &EntityPeople{},
	}
	err := conn.QueryRow(
		ctx,
		`SELECT people.name, people.surname, people.patronymic, cars.regnum, cars.mark, cars.model, cars.year
		FROM cars 
		JOIN people ON cars.owner = people.id
		WHERE cars.regnum = $1;
		`,
		regNum,
	).Scan(
		&entityCar.owner.name,
		&entityCar.owner.surname,
		&entityCar.owner.patronymic,
		&entityCar.regNum,
		&entityCar.mark,
		&entityCar.model,
		&entityCar.year,
	)
	if err != nil {
		logrus.Debugln(err)
		return nil
	}

	entityCar.owner = &EntityPeople{
		name:       entityCar.owner.name,
		surname:    entityCar.owner.surname,
		patronymic: entityCar.owner.patronymic,
	}

	return entityCar
}

func (car *EntityCar) GetFilteredAndPaginatedInfo(
	ctx context.Context,
	conn *pgx.Conn,
	page uint,
	perPage uint,
	totalPages uint,
	filter string,
	values []interface{},
) ([]*EntityCar, int) {
	logrus.Debugln(page, perPage, totalPages, filter)
	result := make([]*EntityCar, 0)

	if page > totalPages {
		return result, http.StatusOK
	}

	var limit int64 = int64(perPage)
	var offset int64 = int64(perPage * (page - 1))

	var args []interface{}
	whereClause := strings.Builder{}

	switch filter {
	case "regNum":
		for i, value := range values {
			if !car.ValidateRegNum(value.(string)) {
				return result, http.StatusBadRequest
			}
			if i > 0 {
				whereClause.WriteString(" OR ")
			}
			whereClause.WriteString("regnum = $" + strconv.Itoa(i+1))
			args = append(args, value.(string))
		}
	case "model":
		for i, value := range values {
			if !car.ValidateModel(value.(string)) {
				return result, http.StatusBadRequest
			}
			if i > 0 {
				whereClause.WriteString(" OR ")
			}
			whereClause.WriteString("model = $" + strconv.Itoa(i+1))
			args = append(args, value.(string))
		}
	case "mark":
		for i, value := range values {
			if !car.ValidateMark(value.(string)) {
				return result, http.StatusBadRequest
			}
			if i > 0 {
				whereClause.WriteString(" OR ")
			}
			whereClause.WriteString("model = $" + strconv.Itoa(i+1))
			args = append(args, value.(string))
		}
	case "name":
		for i, value := range values {
			if !car.GetOwner().ValidateName(value.(string)) {
				return result, http.StatusBadRequest
			}
			if i > 0 {
				whereClause.WriteString(" OR ")
			}
			whereClause.WriteString("people.name = $" + strconv.Itoa(i+1))
			args = append(args, value.(string))
		}
	case "surname":
		for i, value := range values {
			if !car.GetOwner().ValidateSurname(value.(string)) {
				return result, http.StatusBadRequest
			}
			if i > 0 {
				whereClause.WriteString(" OR ")
			}
			whereClause.WriteString("people.surname = $" + strconv.Itoa(i+1))
			args = append(args, value.(string))
		}
	case "patronymic":
		for i, value := range values {
			if !car.GetOwner().ValidatePatronymic(value.(string)) {
				return result, http.StatusBadRequest
			}
			if i > 0 {
				whereClause.WriteString(" OR ")
			}
			whereClause.WriteString("people.patronymic = $" + strconv.Itoa(i+1))
			args = append(args, value.(string))
		}
	case "year":
		for i, value := range values {
			if !car.ValidateYear(value.(int)) {
				return result, http.StatusBadRequest
			}
			if i > 0 {
				whereClause.WriteString(" OR ")
			}
			whereClause.WriteString("year = $" + strconv.Itoa(i+1))
			args = append(args, value.(int))
		}
	}

	var query string
	logrus.Debugln(whereClause.String())
	if whereClause.Len() > 0 {
		query = `
		SELECT people.name, people.surname, people.patronymic, cars.regnum, cars.mark, cars.model, cars.year
		FROM cars JOIN people ON cars.owner = people.id
		WHERE ` + whereClause.String() + ` LIMIT $` + strconv.Itoa(len(values)+1) + ` OFFSET $` + strconv.Itoa(len(values)+2) + `;
		`
	} else {
		query = `
		SELECT people.name, people.surname, people.patronymic, cars.regnum, cars.mark, cars.model, cars.year
		FROM cars JOIN people ON cars.owner = people.id
		LIMIT $` + strconv.Itoa(len(values)+1) + ` OFFSET $` + strconv.Itoa(len(values)+2) + `;
		`
	}

	logrus.Debug(query)
	args = append(args, limit, offset)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		logrus.Debugln(err)
		return nil, http.StatusInternalServerError
	}

	for rows.Next() {
		car := &EntityCar{
			owner: &EntityPeople{},
		}
		err := rows.Scan(
			&car.owner.name,
			&car.owner.surname,
			&car.owner.patronymic,
			&car.regNum,
			&car.mark,
			&car.model,
			&car.year,
		)
		if err != nil {
			logrus.Debugln(err)
			continue
		}
		result = append(result, car)
	}

	if err := rows.Err(); err != nil {
		logrus.Debugln(err)
	}

	return result, http.StatusOK
}

func (car *EntityCar) Count(
	ctx context.Context,
	conn *pgx.Conn,
) uint {
	var count uint
	err := conn.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM CARS;",
	).Scan(
		&count,
	)
	if err != nil {
		logrus.Debugln(err)
		return 0
	}

	return count
}
