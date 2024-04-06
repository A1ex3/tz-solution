package api

import (
	"context"
	"net/http"
	"tzsolution/http/api/people"
	people_ "tzsolution/people"
	"tzsolution/postgresql"

	"github.com/gin-gonic/gin"
)

type peopleApi struct {
	action          *people_.People
	request         *people.PeopleApiRequest
	response        *people.PeopleApiResponse
	defaultResponse *defaultResponses
}

func NewPeopleApi(db *postgresql.Postgresql) *peopleApi {
	newpeople := people_.NewPeople(db)
	newDefaultResponses := NewDefaultResponses()

	return &peopleApi{
		action: newpeople,
		request: &people.PeopleApiRequest{
			Add: &people.PeopleApiRequestAdd{},
		},
		response: &people.PeopleApiResponse{
			Add: &people.PeopleApiResponseAdd{},
		},
		defaultResponse: newDefaultResponses,
	}
}

// @Summary Add a new person
// @Description Add a new person to the database
// @ID add-people
// @Accept json
// @Produce json
// @Param request body people.PeopleApiRequestAdd true "Person data"
// @Success 200 {object} people.PeopleApiResponseAdd
// @Failure 400 {object} defaultResponses
// @Failure 500 {object} defaultResponses
// @Router /user [post]
func (api *peopleApi) Add(c *gin.Context) {
	if err := c.BindJSON(api.request.Add); err != nil {
		c.JSON(http.StatusBadRequest, api.defaultResponse.GetResponse(http.StatusBadRequest))
		return
	}

	ctx := context.Background()
	id, httpStatusCode := api.action.NewPeople(
		ctx,
		api.request.Add.Name,
		api.request.Add.Surname,
		api.request.Add.Patronymic,
	)

	api.response.Add.Id = id
	c.JSON(
		httpStatusCode,
		api.response.Add,
	)
}
