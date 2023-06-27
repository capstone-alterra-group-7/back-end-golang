package controllers_test

import (
	"back-end-golang/controllers"
	"back-end-golang/dtos"
	"back-end-golang/usecases"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestFileUpload(t *testing.T) {
	cloudinaryUsecase := usecases.NewMediaUpload()
	cloudinaryController := controllers.NewCloudinaryController(cloudinaryUsecase)
	LoadEnvTest()
	e := echo.New()
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	imageFile, err := os.Open("./images/Image-pAaMJR.jpg")
	if err != nil {
		t.Fatalf("Failed to open image file: %v", err)
	}
	defer imageFile.Close()

	imagePart, err := writer.CreateFormFile("file", "image.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = io.Copy(imagePart, imageFile)
	if err != nil {
		t.Fatalf("Failed to copy image file to form file: %v", err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/public/cloudinary/file-upload", form)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := cloudinaryController.FileUpload(c); err != nil {
		t.Fatalf("Failed to upload image: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v ", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	t.Logf("Response : %v", response)
}

func TestUrlUpload(t *testing.T) {
	cloudinaryUsecase := usecases.NewMediaUpload()
	cloudinaryController := controllers.NewCloudinaryController(cloudinaryUsecase)
	LoadEnvTest()
	e := echo.New()
	body := map[string]interface{}{
		"url": "https://dragonflyaerospace.com/wp-content/themes/dragonfly/images/page-product/product-chameleon-1.png",
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/public/cloudinary/url-upload", bytes.NewBuffer(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := cloudinaryController.UrlUpload(c); err != nil {
		t.Fatalf("Failed to upload image: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v error %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	t.Logf("Response : %v", response)
}
