package api

import (
	"context"
	"net/http"
	"strconv"
	"tzsolution/cars"
	"tzsolution/http/api/car"
	"tzsolution/postgresql"

	"github.com/gin-gonic/gin"
)

type carApi struct {
	action          *cars.Cars
	request         *car.CarApiRequest
	response        *car.CarApiResponse
	defaultResponse *defaultResponses
}

func NewCarApi(db *postgresql.Postgresql) *carApi {
	newCars := cars.NewCars(db)
	newDefaultResponses := NewDefaultResponses()

	return &carApi{
		action: newCars,
		request: &car.CarApiRequest{
			Add:                         &car.CarApiRequestAdd{},
			Update:                      &car.CarApiRequestUpdate{},
			GetFilteredAndPaginatedInfo: &car.CarApiRequestGetFilteredAndPaginatedInfo{},
		},
		response: &car.CarApiResponse{
			Add: &car.CarApiResponseAdd{},
			GetInfo: &car.CarApiResponseGetInfo{
				Owner: &car.CarApiResponseGetInfoOwner{},
			},
			GetFilteredAndPaginatedInfo: &car.CarApiResponseGetFilteredAndPaginatedInfo{
				List: []*car.CarApiResponseGetInfo{},
			},
		},
		defaultResponse: newDefaultResponses,
	}
}

// @Summary Add a new car
// @ID add-car
// @Accept json
// @Produce json
// @Param ownerId query integer true "Owner ID"
// @Param body body car.CarApiRequestAdd true "Car details"
// @Success 200 {object} car.CarApiResponseAdd
// @Failure 400 {object} defaultResponses
// @Failure 500 {object} car.CarApiResponseAdd
// @Router / [post]
func (api *carApi) Add(c *gin.Context) {
	api.request.Add = &car.CarApiRequestAdd{}

	ownerIdParam := c.Query("ownerId")
	id, err := strconv.Atoi(ownerIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}

	if err := c.BindJSON(api.request.Add); err != nil {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}

	ctx := context.Background()
	id, httpStatusCode := api.action.NewCar(
		ctx,
		id,
		api.request.Add.RegNum,
		api.request.Add.Mark,
		api.request.Add.Model,
		api.request.Add.Year,
	)

	api.response.Add.Id = id
	c.JSON(
		httpStatusCode,
		api.response.Add,
	)
}

// @Summary Delete Car
// @ID delete-car
// @Produce json
// @Param carId query integer true "Car ID"
// @Success 200 {object} defaultResponses
// @Failure 400 {object} defaultResponses
// @Failure 500 {object} defaultResponses
// @Router / [delete]
func (api *carApi) Delete(c *gin.Context) {
	ownerIdParam := c.Query("carId")
	id, err := strconv.Atoi(ownerIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}

	ctx := context.Background()

	_, httpStatusCode := api.action.DeleteCar(
		ctx,
		id,
	)
	c.JSON(httpStatusCode, api.defaultResponse.GetResponse(httpStatusCode))
}

// @Summary Update Car Information
// @ID update-car
// @Accept json
// @Produce json
// @Param carId query integer true "Car ID"
// @Param body body car.CarApiRequestUpdate true "Details for updating the car"
// @Success 200 {object} defaultResponses "Successful response"
// @Failure 400 {object} defaultResponses "Bad request"
// @Failure 500 {object} defaultResponses "Internal server error"
// @Router / [put]
func (api *carApi) Update(c *gin.Context) {
	api.request.Update = &car.CarApiRequestUpdate{}

	ownerIdParam := c.Query("carId")
	id, err := strconv.Atoi(ownerIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}

	if err := c.BindJSON(api.request.Update); err != nil {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}

	if api.request.Update.OwnerId == nil && api.request.Update.RegNum == nil {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}

	requestBodyMap := make(map[string]interface{})
	if api.request.Update.OwnerId != nil {
		requestBodyMap["ownerId"] = *api.request.Update.OwnerId
	}

	if api.request.Update.RegNum != nil {
		requestBodyMap["regNum"] = *api.request.Update.RegNum
	}

	ctx := context.Background()

	_, httpStatusCode := api.action.UpdateCar(
		ctx,
		id,
		requestBodyMap,
	)

	c.JSON(httpStatusCode, api.defaultResponse.GetResponse(httpStatusCode))
}

// @Summary Get car information by registration number
// @ID getInfoByRegNum-car
// @Produce json
// @Param regNum query string true "Car registration number"
// @Success 200 {object} car.CarApiResponseGetInfo "Successful response"
// @Failure 400 {object} defaultResponses "Bad request"
// @Failure 500 {object} defaultResponses "Internal server error"
// @Router /info [get]
func (api *carApi) GetInfoByRegNum(c *gin.Context) {
	regNumParam := c.Query("regNum")
	if regNumParam != "" {
		ctx := context.Background()
		result, httpStatusCode := api.action.GetInfoByRegNum(
			ctx,
			regNumParam,
		)
		if httpStatusCode != http.StatusOK {
			c.JSON(httpStatusCode, api.defaultResponse.GetResponse(httpStatusCode))
			return
		} else {
			if result != nil {
				api.response.GetInfo.Mark = result.GetMark()
				api.response.GetInfo.Model = result.GetModel()
				api.response.GetInfo.RegNum = result.GetRegNum()
				api.response.GetInfo.Year = result.GetYear()
				api.response.GetInfo.Owner.Name = result.GetOwner().GetName()
				api.response.GetInfo.Owner.Surname = result.GetOwner().GetSurname()
				api.response.GetInfo.Owner.Patronymic = result.GetOwner().GetPatronymic()

				c.JSON(httpStatusCode, api.response.GetInfo)
				return
			} else {
				c.JSON(http.StatusInternalServerError, api.defaultResponse.GetResponse(http.StatusInternalServerError))
				return
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}
}

// @Summary Retrieve filtered and paginated car information
// @Description Retrieve car information based on various filters and pagination settings.
// @ID getFilteredAndPaginatedInfo-car
// @Accept  json
// @Produce  json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param body body car.CarApiRequestGetFilteredAndPaginatedInfo false "Filter parameters for car information"
// @Success 200 {object} car.CarApiResponseGetFilteredAndPaginatedInfo "Successful response"
// @Failure 400 {object} defaultResponses "Bad request"
// @Failure 500 {object} car.CarApiResponseGetFilteredAndPaginatedInfo "Internal server error"
// @Router /info/filter [get]
func (api *carApi) GetFilteredAndPaginatedInfo(c *gin.Context) {
	api.request.GetFilteredAndPaginatedInfo = &car.CarApiRequestGetFilteredAndPaginatedInfo{}

	page := 1
	pageParam := c.Query("page")
	if pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
			return
		}
	}

	requestFilter := ""
	var requestValuesList []interface{}

	if c.Request.ContentLength > 0 {
		if err := c.BindJSON(&api.request.GetFilteredAndPaginatedInfo); err != nil {
			c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
			return
		}

		if len(api.request.GetFilteredAndPaginatedInfo.Year) > 0 {
			requestFilter = "year"
			for _, year := range api.request.GetFilteredAndPaginatedInfo.Year {
				requestValuesList = append(requestValuesList, year)
			}
		} else if len(api.request.GetFilteredAndPaginatedInfo.RegNum) > 0 {
			requestFilter = "regNum"
			for _, regNum := range api.request.GetFilteredAndPaginatedInfo.RegNum {
				requestValuesList = append(requestValuesList, regNum)
			}
		} else if len(api.request.GetFilteredAndPaginatedInfo.Mark) > 0 {
			requestFilter = "mark"
			for _, mark := range api.request.GetFilteredAndPaginatedInfo.Mark {
				requestValuesList = append(requestValuesList, mark)
			}
		} else if len(api.request.GetFilteredAndPaginatedInfo.Model) > 0 {
			requestFilter = "model"
			for _, model := range api.request.GetFilteredAndPaginatedInfo.Model {
				requestValuesList = append(requestValuesList, model)
			}
		} else if len(api.request.GetFilteredAndPaginatedInfo.Name) > 0 {
			requestFilter = "name"
			for _, name := range api.request.GetFilteredAndPaginatedInfo.Name {
				requestValuesList = append(requestValuesList, name)
			}
		} else if len(api.request.GetFilteredAndPaginatedInfo.Surname) > 0 {
			requestFilter = "surname"
			for _, surname := range api.request.GetFilteredAndPaginatedInfo.Surname {
				requestValuesList = append(requestValuesList, surname)
			}
		} else if len(api.request.GetFilteredAndPaginatedInfo.Patronymic) > 0 {
			requestFilter = "patronymic"
			for _, patronymic := range api.request.GetFilteredAndPaginatedInfo.Patronymic {
				requestValuesList = append(requestValuesList, patronymic)
			}
		}
	}

	ctx := context.Background()
	result, totalPages, perPage, httpStatusCode := api.action.GetFilteredAndPaginatedInfo(
		ctx,
		requestFilter,
		requestValuesList,
		uint(page),
	)

	if httpStatusCode != http.StatusOK {
		api.response.GetFilteredAndPaginatedInfo.CurrentPage = page
		api.response.GetFilteredAndPaginatedInfo.PerPage = perPage
		api.response.GetFilteredAndPaginatedInfo.TotalPages = totalPages
		api.response.GetFilteredAndPaginatedInfo.List = []*car.CarApiResponseGetInfo{}

		c.JSON(http.StatusBadRequest, api.response.GetFilteredAndPaginatedInfo)
		return
	} else {
		api.response.GetFilteredAndPaginatedInfo.CurrentPage = page
		api.response.GetFilteredAndPaginatedInfo.PerPage = perPage
		api.response.GetFilteredAndPaginatedInfo.TotalPages = totalPages

		list := make([]*car.CarApiResponseGetInfo, 0)
		for _, val := range result {
			entityReponse := &car.CarApiResponseGetInfo{
				RegNum: val.GetRegNum(),
				Mark:   val.GetMark(),
				Model:  val.GetModel(),
				Year:   val.GetYear(),
				Owner: &car.CarApiResponseGetInfoOwner{
					Name:       val.GetOwner().GetName(),
					Surname:    val.GetOwner().GetSurname(),
					Patronymic: val.GetOwner().GetPatronymic(),
				},
			}
			list = append(list, entityReponse)
		}

		api.response.GetFilteredAndPaginatedInfo.List = list
		c.JSON(http.StatusOK, api.response.GetFilteredAndPaginatedInfo)
		return
	}
}
