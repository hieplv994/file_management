package model

import "time"

//File model
type File struct {
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	Hash        string    `json:"hash"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}
