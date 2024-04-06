package people

type PeopleApiRequest struct{
	Add *PeopleApiRequestAdd
}

type PeopleApiRequestAdd struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}
