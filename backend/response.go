package backend

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func Success(data interface{}) Response {
	return Response{Success: true, Data: data}
}

func Failure(error string) Response {
	return Response{Success: false, Error: error}
}
