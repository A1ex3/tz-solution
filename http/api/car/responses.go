package car

type CarApiResponse struct {
	Add                         *CarApiResponseAdd
	GetInfo                     *CarApiResponseGetInfo
	GetFilteredAndPaginatedInfo *CarApiResponseGetFilteredAndPaginatedInfo
}

type CarApiResponseAdd struct {
	Id int `json:"id"`
}

type CarApiResponseGetInfo struct {
	RegNum string                      `json:"regNum"`
	Mark   string                      `json:"mark"`
	Model  string                      `json:"model"`
	Year   int                         `json:"year"`
	Owner  *CarApiResponseGetInfoOwner `json:"owner"`
}

type CarApiResponseGetInfoOwner struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type CarApiResponseGetFilteredAndPaginatedInfo struct {
	CurrentPage int                      `json:"currentPage"`
	PerPage     uint                     `json:"perPage"`
	TotalPages  uint                     `json:"totalPages"`
	List        []*CarApiResponseGetInfo `json:"list"`
}
