package model

//Error model
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
