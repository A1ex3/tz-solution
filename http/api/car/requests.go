package car

type CarApiRequest struct {
	Add                         *CarApiRequestAdd
	Update                      *CarApiRequestUpdate
	GetFilteredAndPaginatedInfo *CarApiRequestGetFilteredAndPaginatedInfo
}

type CarApiRequestAdd struct {
	Year   int    `json:"year"`
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
}

type CarApiRequestUpdate struct {
	OwnerId *int    `json:"ownerId,omitempty"`
	RegNum  *string `json:"regNum,omitempty"`
}

type CarApiRequestGetFilteredAndPaginatedInfo struct {
	Year       []int    `json:"year,omitempty"`
	RegNum     []string `json:"regNum,omitempty"`
	Mark       []string `json:"mark,omitempty"`
	Model      []string `json:"model,omitempty"`
	Name       []string `json:"name,omitempty"`
	Surname    []string `json:"surname,omitempty"`
	Patronymic []string `json:"patronymic,omitempty"`
}
