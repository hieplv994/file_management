package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-file-api-v2/driver"
	models "go-file-api-v2/model"
	"go-file-api-v2/repository/repoimpl"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

//HomePage wellcome
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func hashFileMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil

}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

//InsertFileInforInFolder insert all file in directory
func InsertFileInforInFolder(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query["path"][0]
	var files []string
	err := filepath.Walk(path, func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
		ResponseErr(w, http.StatusNotFound)
	}
	for i, p := range files {
		fmt.Println(i, p)
		fileInfo, err := os.Stat(p)
		hash, err := hashFileMd5(p)
		if err != nil {
			ResponseErr(w, http.StatusNotFound)
			return
		}
		statTime := fileInfo.Sys().(*syscall.Stat_t)
		fmt.Println(hash)
		temp := models.File{
			Path:        p,
			Size:        fileInfo.Size(),
			Hash:        hash,
			DateCreated: timespecToTime(statTime.Ctim),
			DateUpdated: timespecToTime(statTime.Mtim),
		}
		fileRepo := repoimpl.NewFileRepo(driver.Postgres.SQL)
		fileRepo.Insert(temp)
	}
	ResponseOk(w, "Insert Success")
}

//GetFilePath insert all file path
func GetFilePath(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query["path"][0]
	fileRepo := repoimpl.NewFileRepo(driver.Postgres.SQL)

	fileTemp, err := fileRepo.SelectByPath(path)
	if err != nil || fileTemp.Path == "" {
		ResponseErr(w, http.StatusNotFound)
		return
	}
	ResponseOk(w, fileTemp)
}

//UpdateFile insert all file path
func UpdateFile(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fileRepo := repoimpl.NewFileRepo(driver.Postgres.SQL)
	var file models.File
	json.Unmarshal(reqBody, &file)
	err := fileRepo.Update(file)
	if err != nil {
		ResponseErr(w, http.StatusNotFound)
		return
	}
	ResponseOk(w, "Update file Success")
}

//ResponseErr handle response error
func ResponseErr(w http.ResponseWriter, statusCode int) {
	jData, err := json.Marshal(models.Error{
		Status:  statusCode,
		Message: http.StatusText(statusCode),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

//ResponseOk handle response success
func ResponseOk(w http.ResponseWriter, data interface{}) {
	if data == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
