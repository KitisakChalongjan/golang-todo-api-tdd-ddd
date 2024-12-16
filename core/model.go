package core

type ApiRespose struct {
	Error string      `json:"message"`
	Data  interface{} `json:"data"`
}
