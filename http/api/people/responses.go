package people

type PeopleApiResponse struct {
	Add *PeopleApiResponseAdd
}

type PeopleApiResponseAdd struct {
	Id int `json:"id"`
}
