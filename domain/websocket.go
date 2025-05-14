package domain

type WSResponse struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
