package main

import (
	"fmt"
	"go-file-api-v2/driver"
	"go-file-api-v2/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//host, port, user, password, dbname
const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "123456"
	dbname   = "file_management"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func main() {
	db := driver.Connect(host, port, user, password, dbname)

	err := db.SQL.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB Successfully connected!")
	fmt.Println("File Management Rest API v2.0 - Server start....")
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", handler.HomePage)
	myRouter.HandleFunc("/get_file_info", handler.GetFilePath)
	myRouter.HandleFunc("/update_file_info", handler.UpdateFile).Methods("PUT")
	myRouter.HandleFunc("/insert_file_info", handler.InsertFileInforInFolder)

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
