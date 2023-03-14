package responses

type TestResponse[T interface{}] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
