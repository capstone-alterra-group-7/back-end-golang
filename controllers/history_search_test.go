package controllers_test

import (
	"back-end-golang/controllers"
	"back-end-golang/dtos"
	"back-end-golang/repositories"
	"back-end-golang/usecases"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetAllDataHistorySearch(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	LoadEnvTest()
	historySearchRepo := repositories.NewHistorySearchRepository(db)
	userRepo := repositories.NewUserRepository(db)
	historySearchUsecase := usecases.NewHistorySearchUsecase(historySearchRepo, userRepo)
	historySearchController := controllers.NewHistorySearchController(historySearchUsecase)
	e := echo.New()

	token, err := loginUser()
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/user/history-search", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Add("page", "1")
	c.QueryParams().Add("limit", "10")
	if err := historySearchController.HistorySearchGetAll(c); err != nil {
		t.Fatalf("Failed to get all history search: %v", err)
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
