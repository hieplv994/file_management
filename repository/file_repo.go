package repository

import (
	models "go-file-api-v2/model"
)

//FileRepo interface
type FileRepo interface {
	Insert(u models.File) error
	SelectByPath(path string) (models.File, error)
	Update(u models.File) error
}
