package controllers_test

import (
	"back-end-golang/controllers"
	"back-end-golang/dtos"
	"back-end-golang/repositories"
	"back-end-golang/usecases"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginUser() (string, error) {
	db, err := ConnectDBTest()
	if err != nil {
		return "", fmt.Errorf("failed to connect database: %v", err)
	}

	LoadEnvTest()
	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)
	userController := controllers.NewUserController(userUsecase)
	e := echo.New()

	userLoginRequest := UserLoginRequest{
		Email:    "user@gmail.com",
		Password: "qweqwe123",
	}

	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserLogin(c)
	if err != nil {
		return "", fmt.Errorf("failed to login user: %v", err)
	}
	if rec.Code != http.StatusOK {
		return "", fmt.Errorf("expected status code %v but got %v", http.StatusOK, rec.Code)
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	token := response.Data.(map[string]interface{})["token"].(string)
	return token, nil
}
func loginAdmin() (string, error) {
	db, err := ConnectDBTest()
	if err != nil {
		return "", fmt.Errorf("failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)
	userController := controllers.NewUserController(userUsecase)
	e := echo.New()

	userLoginRequest := UserLoginRequest{
		Email:    "admin@gmail.com",
		Password: "qweqwe123",
	}

	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserLogin(c)
	if err != nil {
		return "", fmt.Errorf("failed to login user: %v", err)
	}
	if rec.Code != http.StatusOK {
		return "", fmt.Errorf("expected status code %v but got %v", http.StatusOK, rec.Code)
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	token := response.Data.(map[string]interface{})["token"].(string)
	return token, nil
}

func TestUserLogin(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	LoadEnvTest()
	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)
	userController := controllers.NewUserController(userUsecase)
	e := echo.New()
	userLoginRequest := UserLoginRequest{
		Email:    "user@gmail.com",
		Password: "qweqwe123",
	}

	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserLogin(c)
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}
	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}

func TestUserRegister(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)
	userController := controllers.NewUserController(userUsecase)
	e := echo.New()
	birthDate := "2000-01-01"
	isActive := true
	userRegisterRequest := dtos.UserRegisterInput{
		FullName:        "testRegis",
		Email:           "testRegis1568@gmail.com",
		Password:        "qweqwe123",
		ConfirmPassword: "qweqwe123",
		BirthDate:       &birthDate,
		PhoneNumber:     "0851555555151",
		Role:            "user",
		IsActive:        &isActive,
	}

	reqBody, err := json.Marshal(userRegisterRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserRegister(c)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func TestUserUpdateProfile(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)

	userController := controllers.NewUserController(userUsecase)
	e := echo.New()

	token, err := loginUser()
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	birthDate := "2000-01-01"
	userUpdateProfileRequest := dtos.UserUpdateProfileInput{
		FullName:    "testUpdate",
		PhoneNumber: "085555555551",
		BirthDate:   birthDate,
		Citizen:     "Indonesia",
	}

	reqBody, err := json.Marshal(userUpdateProfileRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/update-profile", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserUpdateProfile(c)
	if err != nil {
		t.Fatalf("Failed to update profile user: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func TestUserGetCredential(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)

	userController := controllers.NewUserController(userUsecase)
	e := echo.New()

	token, err := loginUser()
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/get-credential", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserCredential(c)
	if err != nil {
		t.Fatalf("Failed to get credential user: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func TestUserUpdatePhotoProfile(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)

	userController := controllers.NewUserController(userUsecase)
	e := echo.New()

	token, err := loginUser()
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}
	// LoadEnv()
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

	userPhotoProfile := dtos.UserUpdatePhotoProfileInput{
		ProfilePicture: "",
	}

	_ = writer.WriteField("ProfilePicture ", userPhotoProfile.ProfilePicture)

	writer.Close()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/update-photo-profile", form)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err = userController.UserUpdatePhotoProfile(c)
	if err != nil {
		t.Fatalf("Failed to update photo profile user: %v", err)
	}

	if rec.Code != http.StatusOK {

		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	t.Logf("Response body: %v", response)
}

func TestUserDeletePhotoProfile(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)

	userController := controllers.NewUserController(userUsecase)

	e := echo.New()

	token, err := loginUser()
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/user/delete-photo-profile", nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserDeletePhotoProfile(c)
	if err != nil {
		t.Fatalf("Failed to delete photo profile user: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {

		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func TestGetAllUser(t *testing.T) {

	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)

	userController := controllers.NewUserController(userUsecase)

	e := echo.New()

	token, err := loginAdmin()
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/user", nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Add("page", "1")
	c.QueryParams().Add("limit", "10")
	c.QueryParams().Add("search", "test")
	c.QueryParams().Add("filter", "asc")

	err = userController.UserGetAll(c)
	if err != nil {
		t.Fatalf("Failed to get all user: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {

		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func TestUserGetDetail(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)

	userController := controllers.NewUserController(userUsecase)

	e := echo.New()

	token, err := loginUser()

	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/user/detail", nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Add("id", "2")
	c.QueryParams().Add("isDeleted", "false")

	err = userController.UserGetDetail(c)
	if err != nil {
		t.Fatalf("Failed to get detail user: %v", err)
	}

	if rec.Code != http.StatusOK {

		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {

		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func TestUserAdminUpdate(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)
	userController := controllers.NewUserController(userUsecase)

	e := echo.New()

	token, err := loginAdmin()
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	isActive := false
	BirtDate := "2000-01-01"
	userAdminUpdateRequest := dtos.UserRegisterInputUpdateByAdmin{
		FullName:    "testUpdateByAdmin",
		Email:       "testregis090909@gmail.com",
		Role:        "user",
		PhoneNumber: "8445",
		IsActive:    &isActive,
		BirthDate:   &BirtDate,
	}

	reqBody, err := json.Marshal(userAdminUpdateRequest)

	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/user/update/", bytes.NewReader(reqBody))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("2")
	err = userController.UserAdminUpdate(c)

	if err != nil {
		t.Fatalf("Failed to update user by admin: %v", err)
	}

	if rec.Code != http.StatusCreated {

		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse

	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {

		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func TestUserUpdatePassword(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)

	userController := controllers.NewUserController(userUsecase)
	e := echo.New()

	token, err := loginUserUpdatePWD()
	if err != nil {
		t.Fatalf("Failed to login user: %v", err)
	}

	userUpdatePasswordRequest := dtos.UserUpdatePasswordInput{
		OldPassword:     "qweqwe123",
		NewPassword:     "faros123",
		ConfirmPassword: "faros123",
	}

	reqBody, err := json.Marshal(userUpdatePasswordRequest)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/user/update-password", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// c.Set("Authorization", "Bearer "+token)

	err = userController.UserUpdatePassword(c)
	if err != nil {
		t.Fatalf("Failed to update password user: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	t.Logf("Response body: %v", response)
}

func loginUserUpdatePWD() (string, error) {
	db, err := ConnectDBTest()
	if err != nil {
		return "", fmt.Errorf("failed to connect database: %v", err)
	}

	LoadEnvTest()
	userRepo := repositories.NewUserRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, notificationRepo)
	userController := controllers.NewUserController(userUsecase)
	e := echo.New()

	userLoginRequest := UserLoginRequest{
		Email:    "testregis090909@gmail.com",
		Password: "qweqwe123",
	}

	reqBody, err := json.Marshal(userLoginRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = userController.UserLogin(c)
	if err != nil {
		return "", fmt.Errorf("failed to login user: %v", err)
	}
	if rec.Code != http.StatusOK {
		return "", fmt.Errorf("expected status code %v but got %v", http.StatusOK, rec.Code)
	}

	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	token := response.Data.(map[string]interface{})["token"].(string)
	return token, nil
}
