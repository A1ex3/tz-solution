package car

type CarApiRequest struct {
	Add                         *CarApiRequestAdd
	Update                      *CarApiRequestUpdate
	GetFilteredAndPaginatedInfo *CarApiRequestGetFilteredAndPaginatedInfo
}

type CarApiRequestAdd struct {
	RegNum string `json:"regNum"`
	Mark   string `json:"mark"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
}

type CarApiRequestUpdate struct {
	RegNum  *string `json:"regNum,omitempty"`
	OwnerId *int    `json:"ownerId,omitempty"`
}

type CarApiRequestGetFilteredAndPaginatedInfo struct {
	RegNum     []string `json:"regNum,omitempty"`
	Mark       []string `json:"mark,omitempty"`
	Model      []string `json:"model,omitempty"`
	Name       []string `json:"name,omitempty"`
	Surname    []string `json:"surname,omitempty"`
	Patronymic []string `json:"patronymic,omitempty"`
	Year       []int    `json:"year,omitempty"`
}
