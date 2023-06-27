package controllers_test

import (
	"back-end-golang/controllers"
	"back-end-golang/dtos"
	"back-end-golang/repositories"
	"back-end-golang/usecases"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadEnvTest() {

	// kalau mau testing di local sesuaikan direktory proyeknya...
	// 1. digunakan untuk create article karena untuk mengambil file image di folder images
	// 2. digunakan untuk update juga

	// Catatam Lain :
	// Token harus set Manual / menggunakan func LoginUser() & LoginAdmin()
	err := os.Chdir("D:\\CAPSTONE_ALTERRA\\back-end-golang")
	if err != nil {
		log.Fatalf("Failed to set working directory: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func ConnectDBTest() (*gorm.DB, error) {
	// Load the Asia/Jakarta location
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle the error
	}
	DB_HOST := "localhost"
	DB_USER := "root"
	DB_PASSWORD := ""
	DB_NAME := "capstone_be"
	DB_PORT := "3306"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DB_USER,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_NAME,
	)

	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	dbConn = dbConn.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().In(location)
		},
	})

	return dbConn, nil
}

func TestCreateArticle(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	articleRepo := repositories.NewArticleRepository(db)
	articleUsecase := usecases.NewArticleUsecase(articleRepo)
	articleCtrl := controllers.NewArticleController(articleUsecase)

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

	articleInput := &dtos.ArticleInput{
		Title:       "Test Title",
		Image:       "image.jpg",
		Description: "Test Description",
		Label:       "Test Label",
	}

	_ = writer.WriteField("title", articleInput.Title)
	_ = writer.WriteField("description", articleInput.Description)
	_ = writer.WriteField("label", articleInput.Label)

	writer.Close()
	tokenAfterLogin, err := loginUser()
	fmt.Println("token", tokenAfterLogin)
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/article", form)
	req.Header.Set("Authorization", "Bearer "+tokenAfterLogin)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetRequest(req.WithContext(ctx.Request().Context()))

	err = articleCtrl.CreateArticle(ctx)
	// require.NoError(t, err)
	// require.Equal(t, http.StatusCreated, rec.Code)

	if err != nil {
		t.Fatalf("Failed to create article: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusCreated, rec.Code, rec.Body.String())
	}

	responseBody := rec.Body.String()
	require.NotEmpty(t, responseBody)
	t.Logf("Response body: %v", responseBody)
}
func TestGetAllArticles(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/article?page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	// db, err := configs.ConnectDB()
	// if err != nil {
	// 	t.Fatalf("Failed to connect to database: %v", err)
	// }
	articleRepo := repositories.NewArticleRepository(db)
	articleUsecase := usecases.NewArticleUsecase(articleRepo)
	articleCtrl := controllers.NewArticleController(articleUsecase)
	err = articleCtrl.GetAllArticles(ctx)
	if err != nil {
		t.Errorf("Failed to get all articles: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d error %v", http.StatusOK, rec.Code, rec.Body.String())
	}
	responseBody := rec.Body.String()
	if responseBody == "" {
		t.Errorf("Expected response body not empty but got empty")
	}
	t.Logf("Response body: %v", responseBody)
}

func TestGetArticleByID(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)

	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/article", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	// db, err := configs.ConnectDB()
	// if err != nil {
	// 	t.Fatalf("Failed to connect to database: %v", err)
	// }
	articleRepo := repositories.NewArticleRepository(db)
	articleUsecase := usecases.NewArticleUsecase(articleRepo)
	articleCtrl := controllers.NewArticleController(articleUsecase)
	err = articleCtrl.GetArticleByID(ctx)
	if err != nil {
		t.Errorf("Failed to get article by ID: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d error %v", http.StatusOK, rec.Code, rec.Body.String())
	}
	responseBody := rec.Body.String()
	if responseBody == "" {
		t.Errorf("Expected response body not empty but got empty")
	}
	t.Logf("Response body: %v", responseBody)
}

func TestUpdateArticle(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	articleRepo := repositories.NewArticleRepository(db)
	articleUsecase := usecases.NewArticleUsecase(articleRepo)
	articleCtrl := controllers.NewArticleController(articleUsecase)

	LoadEnvTest()

	e := echo.New()
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	imageFile, err := os.Open("./images/Image-QHCfyF.jpg")
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

	articleInput := &dtos.ArticleInput{
		Title:       "Test Update 2",
		Image:       "image.jpg",
		Description: "Test Update 2",
		Label:       "Test Update 2",
	}

	_ = writer.WriteField("title", articleInput.Title)
	_ = writer.WriteField("description", articleInput.Description)
	_ = writer.WriteField("label", articleInput.Label)

	writer.Close()
	tokenAfterLogin := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODc1MDgwNTEsInJvbGUiOiJhZG1pbiIsInVzZXJJZCI6MX0.06tBCLU-tZzqIjjfgV3-OB562sclKZCUG64jqfvCncM"
	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/article/", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+tokenAfterLogin)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	ctx.SetRequest(req.WithContext(ctx.Request().Context()))

	err = articleCtrl.UpdateArticle(ctx)
	// require.NoError(t, err)
	// require.Equal(t, http.StatusOK, rec.Code)
	if err != nil {
		t.Fatalf("Failed to update article: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	responseBody := rec.Body.String()
	require.NotEmpty(t, responseBody)
	t.Logf("Response body: %v", responseBody)
}

func TestDeleteArticle(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)

	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/public/article", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("1")
	tokenAfterLogin := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODc1MDgwNTEsInJvbGUiOiJhZG1pbiIsInVzZXJJZCI6MX0.06tBCLU-tZzqIjjfgV3-OB562sclKZCUG64jqfvCncM"

	req.Header.Set("Authorization", "Bearer "+tokenAfterLogin)
	// db, err := configs.ConnectDB()
	// if err != nil {
	// 	t.Fatalf("Failed to connect to database: %v", err)
	// }
	articleRepo := repositories.NewArticleRepository(db)
	articleUsecase := usecases.NewArticleUsecase(articleRepo)
	articleCtrl := controllers.NewArticleController(articleUsecase)
	err = articleCtrl.DeleteArticle(ctx)

	if err != nil {
		t.Errorf("Failed delete get article by ID: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d error %v", http.StatusOK, rec.Code, rec.Body.String())
	}
	responseBody := rec.Body.String()
	if responseBody == "" {
		t.Errorf("Expected response body not empty but got empty")
	}
	t.Logf("Response body: %v", responseBody)
}
