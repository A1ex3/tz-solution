package cars

import (
	"context"
	"net/http"
	"time"
	"tzsolution/postgresql"
	"tzsolution/postgresql/entities"
)

type Cars struct {
	Db *postgresql.Postgresql
}

func NewCars(db *postgresql.Postgresql) *Cars {
	return &Cars{
		Db: db,
	}
}

func (cars *Cars) NewCar(
	ctx context.Context,
	ownerId int,
	regNum string,
	mark string,
	model string,
	year int,
) (int, int) {
	entityCar := entities.EntityCar{}

	if !entityCar.GetOwner().ValidateId(ownerId) ||
		!entityCar.ValidateRegNum(regNum) ||
		!entityCar.ValidateMark(mark) ||
		!entityCar.ValidateModel(model) ||
		!entityCar.ValidateYear(year) {
		return 0, http.StatusBadRequest
	}

	dbContext, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	dbConn, err := cars.Db.Connect(dbContext)
	defer cars.Db.Close(ctx, dbConn)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	result := entityCar.Add(
		dbContext,
		dbConn,
		ownerId,
		regNum,
		mark,
		model,
		year,
	)
	if result > 0 {
		return result, http.StatusOK
	} else {
		return result, http.StatusInternalServerError
	}
}

func (cars *Cars) UpdateCar(
	ctx context.Context,
	carId int,
	fields map[string]interface{},
) (int, int) {
	entityCar := entities.EntityCar{}

	if !entityCar.ValidateId(carId) {
		return 0, http.StatusBadRequest
	}

	dbContext, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	dbConn, err := cars.Db.Connect(dbContext)
	defer cars.Db.Close(ctx, dbConn)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	result, httpStatus := entityCar.UpdateByCarId(
		dbContext,
		dbConn,
		carId,
		fields,
	)
	if result > 0 {
		return result, httpStatus
	} else {
		return result, http.StatusInternalServerError
	}
}

func (cars *Cars) GetInfoByRegNum(
	ctx context.Context,
	regNum string,
) (*entities.EntityCar, int) {
	entityCar := &entities.EntityCar{}

	if !entityCar.ValidateRegNum(regNum) {
		return entityCar, http.StatusBadRequest
	}

	dbContext, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	dbConn, err := cars.Db.Connect(dbContext)
	defer cars.Db.Close(ctx, dbConn)
	if err != nil {
		return entityCar, http.StatusInternalServerError
	}

	return entityCar.GetInfoByRegNum(
		dbContext,
		dbConn,
		regNum,
	), http.StatusOK
}

func (cars *Cars) GetFilteredAndPaginatedInfo(
	ctx context.Context,
	filter string,
	values []interface{},
	page uint,
) ([]*entities.EntityCar, uint, uint, int) {
	const perPage uint = 30
	entityCar := entities.EntityCar{}

	dbContext, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	dbConn, err := cars.Db.Connect(dbContext)
	defer cars.Db.Close(ctx, dbConn)
	if err != nil {
		return []*entities.EntityCar{}, 0, 0, http.StatusInternalServerError
	}

	totalPages := (entityCar.Count(dbContext, dbConn) + perPage - 1) / perPage
	result, httpStatus := entityCar.GetFilteredAndPaginatedInfo(
		dbContext,
		dbConn,
		page,
		perPage,
		totalPages,
		filter,
		values,
	)

	return result, totalPages, perPage, httpStatus
}

func (cars *Cars) DeleteCar(
	ctx context.Context,
	carId int,
) (int, int) {
	entityCar := entities.EntityCar{}

	if !entityCar.ValidateId(carId) {
		return 0, http.StatusBadRequest
	}

	dbContext, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	dbConn, err := cars.Db.Connect(dbContext)
	defer cars.Db.Close(ctx, dbConn)
	if err != nil {
		return 0, http.StatusInternalServerError
	}

	result := entityCar.DeleteByCarId(
		ctx,
		dbConn,
		carId,
	)

	if result > 0 {
		return result, http.StatusOK
	} else {
		return result, http.StatusInternalServerError
	}
}
