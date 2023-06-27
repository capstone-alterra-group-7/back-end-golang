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

func TestGetAllDataDashboard(t *testing.T) {
	db, err := ConnectDBTest()
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}
	dashboardRepo := repositories.NewDashboardRepository(db)
	userRepo := repositories.NewUserRepository(db)
	ticketOrderRepo := repositories.NewTicketOrderRepository(db)
	ticketTravelerDetailRepo := repositories.NewTicketTravelerDetailRepository(db)
	travelerDetailRepo := repositories.NewTravelerDetailRepository(db)
	trainCarriageRepo := repositories.NewTrainCarriageRepository(db)
	trainRepo := repositories.NewTrainRepository(db)
	trainSeatRepo := repositories.NewTrainSeatRepository(db)
	stationRepo := repositories.NewStationRepository(db)
	trainStationRepo := repositories.NewTrainStationRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	hotelOrderRepo := repositories.NewHotelOrderRepository(db)
	hotelRepo := repositories.NewHotelRepository(db)
	dashboardUsecase := usecases.NewDashboardUsecase(dashboardRepo, userRepo, ticketOrderRepo, ticketTravelerDetailRepo, travelerDetailRepo, trainCarriageRepo, trainRepo, trainSeatRepo, stationRepo, trainStationRepo, paymentRepo, hotelOrderRepo, hotelRepo)
	dashboardController := controllers.NewDashboardController(dashboardUsecase)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/public/dashboard", nil)
	tokenAfterLogin := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2ODc1MDgwNTEsInJvbGUiOiJhZG1pbiIsInVzZXJJZCI6MX0.06tBCLU-tZzqIjjfgV3-OB562sclKZCUG64jqfvCncM"
	req.Header.Set("Authorization", "Bearer "+tokenAfterLogin)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := dashboardController.DashboardGetAll(c); err != nil {
		t.Fatalf("Failed to get all dashboard: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status code %v but got %v errors %v", http.StatusOK, rec.Code, rec.Body.String())
	}
	var response dtos.StatusOKResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}
