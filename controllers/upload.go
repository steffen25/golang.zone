package controllers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/steffen25/golang.zone/util"
)

const (
	MB = 1 << 20
)

type UploadController struct{}
type UploadImageResponse struct {
	ImageURL string `json:"imageUrl"`
}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (uc *UploadController) UploadImage(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-type")
	if !strings.Contains(contentType, "multipart/form-data") {
		NewAPIError(&APIError{false, "Invalid request body. Request body must be of type multipart/form-data", http.StatusBadRequest}, w)
		return
	}
	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, 2*MB)

	if err := r.ParseMultipartForm(2 * MB); err != nil {
		NewAPIError(&APIError{false, "The file you are uploading is too big. Maximum file size is 2MB", http.StatusBadRequest}, w)
		return
	}

	var Buf bytes.Buffer
	file, header, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			NewAPIError(&APIError{false, "Image is required", http.StatusBadRequest}, w)
			return
		}
		NewAPIError(&APIError{false, "Error processing multipart data", http.StatusBadRequest}, w)
		return
	}
	defer file.Close()
	// Copy the file data to my buffer
	io.Copy(&Buf, file)

	fileExtension := http.DetectContentType(Buf.Bytes())
	validFileExtensions := map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
		"image/gif":  nil,
	}

	if _, ok := validFileExtensions[fileExtension]; !ok {
		NewAPIError(&APIError{false, "Invalid mime type, file must be of image/jpeg, image/png or image/gif", http.StatusBadRequest}, w)
		return
	}

	name := strings.Split(header.Filename, ".")
	fileExt := name[len(name)-1]

	now := time.Now()
	fileName := now.Format("2006-01-02_15-04-05") + "_" + util.GetMD5Hash(now.String())

	err = ioutil.WriteFile("./public/images/"+fileName+"."+fileExt, Buf.Bytes(), 0644)
	if err != nil {
		NewAPIError(&APIError{false, "Could not write file to disk", http.StatusInternalServerError}, w)
		return
	}

	Buf.Reset()

	// TODO: Remove hardcoded url
	imageSrc := util.GetRequestScheme(r) + r.Host + "/api/v1/public/" + fileName + "." + fileExt

	data := UploadImageResponse{imageSrc}

	NewAPIResponse(&APIResponse{Success: true, Message: "Image uploaded successfully", Data: data}, w, http.StatusOK)
}
