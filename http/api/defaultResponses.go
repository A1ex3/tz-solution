package api

type defaultResponses struct {
	Description string `json:"description"`
	StatusCode  int    `json:"statusCode"`
}

func NewDefaultResponses() *defaultResponses {
	return &defaultResponses{}
}

func (dr *defaultResponses) GetResponse(statusCode int) *defaultResponses {
	switch statusCode {
	case 200:
		dr.StatusCode = 200
		dr.Description = "OK"
	case 400:
		dr.StatusCode = 400
		dr.Description = "Bad Request"
	case 401:
		dr.StatusCode = 401
		dr.Description = "Unauthorized"
	case 403:
		dr.StatusCode = 403
		dr.Description = "Forbidden"
	case 404:
		dr.StatusCode = 404
		dr.Description = "Not Found"
	case 500:
		dr.StatusCode = 500
		dr.Description = "Internal Server Error"
	default:
		dr.StatusCode = statusCode
		dr.Description = "Unknown Status Code"
	}

	return dr
}
