package repoimpl

import (
	"database/sql"
	"fmt"
	models "go-file-api-v2/model"
	repo "go-file-api-v2/repository"
)

//FileRepoImpl model
type FileRepoImpl struct {
	Db *sql.DB
}

//NewFileRepo instance
func NewFileRepo(db *sql.DB) repo.FileRepo {
	return &FileRepoImpl{
		Db: db,
	}
}

//Insert files information
func (f *FileRepoImpl) Insert(file models.File) error {
	insertStatement := `
	INSERT INTO files (path, size, hash, date_created, date_updated)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := f.Db.Exec(insertStatement, file.Path, file.Size, file.Hash, file.DateCreated, file.DateUpdated)
	if err != nil {
		return err
	}

	fmt.Println("Record added: ", file)

	return nil
}

//SelectByPath files information
func (f *FileRepoImpl) SelectByPath(path string) (models.File, error) {
	row := f.Db.QueryRow("SELECT * FROM files WHERE path = $1", path)
	file := models.File{}
	err := row.Scan(&file.Path, &file.Size, &file.Hash, &file.DateCreated, &file.DateUpdated)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return file, err
	case nil:
		return file, nil
	default:
		panic(err)
	}
}

//Update file information
func (f *FileRepoImpl) Update(file models.File) error {
	updateStatement := `
	UPDATE files SET size=$2, hash=$3, date_created=$4, date_updated=$5
	WHERE path=$1`

	_, err := f.Db.Exec(updateStatement, file.Path, file.Size, file.Hash, file.DateCreated, file.DateUpdated)
	if err != nil {
		return err
	}

	fmt.Println("Record updated: ", file)
	return nil
}
