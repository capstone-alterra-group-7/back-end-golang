package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/repositories"
	"time"

	"github.com/labstack/echo/v4"
)

type DashboardUsecase interface {
	DashboardGetAll() (dtos.DashboardResponse, error)
}

type dashboardUsecase struct {
	dashboardRepository repositories.DashboardRepository
	userRepo            repositories.UserRepository

	ticketOrderRepo          repositories.TicketOrderRepository
	ticketTravelerDetailRepo repositories.TicketTravelerDetailRepository
	travelerDetailRepo       repositories.TravelerDetailRepository
	trainCarriageRepo        repositories.TrainCarriageRepository
	trainRepo                repositories.TrainRepository
	trainSeatRepo            repositories.TrainSeatRepository
	stationRepo              repositories.StationRepository
	trainStationRepo         repositories.TrainStationRepository
	paymentRepo              repositories.PaymentRepository
	hotelOrderRepo           repositories.HotelOrderRepository
	hotelRepo                repositories.HotelRepository
}

func NewDashboardUsecase(dashboardRepository repositories.DashboardRepository, userRepo repositories.UserRepository, ticketOrderRepo repositories.TicketOrderRepository, ticketTravelerDetailRepo repositories.TicketTravelerDetailRepository, travelerDetailRepo repositories.TravelerDetailRepository, trainCarriageRepo repositories.TrainCarriageRepository, trainRepo repositories.TrainRepository, trainSeatRepo repositories.TrainSeatRepository, stationRepo repositories.StationRepository, trainStationRepo repositories.TrainStationRepository, paymentRepo repositories.PaymentRepository, hotelOrderRepo repositories.HotelOrderRepository, hotelRepo repositories.HotelRepository) DashboardUsecase {
	return &dashboardUsecase{dashboardRepository, userRepo, ticketOrderRepo, ticketTravelerDetailRepo, travelerDetailRepo, trainCarriageRepo, trainRepo, trainSeatRepo, stationRepo, trainStationRepo, paymentRepo, hotelOrderRepo, hotelRepo}
}

// DashboardGetAll godoc
// @Summary      Get dashboard
// @Description  Get dashboard
// @Tags         Admin - Dashboard
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.DashboardStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/dashboard [get]
// @Security BearerAuth
func (u *dashboardUsecase) DashboardGetAll() (dtos.DashboardResponse, error) {
	var dashboardResponse dtos.DashboardResponse
	countUser, countUserToday, countTrain, countTrainToday, countTicketOrder, countTicketOrderToday, newTicketOrder, newUser, newHotelOrder, countHotelOrder, countHotelOrderToday, err := u.dashboardRepository.DashboardGetAll()
	if err != nil {
		return dashboardResponse, err
	}
	var newOrderResponses []map[string]interface{}
	for _, ticketTravelerDetail := range newTicketOrder {

		getTicketOrder, err := u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, 1)
		if err != nil {
			return dashboardResponse, err
		}

		getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetail.TrainID))
		if err != nil {
			return dashboardResponse, err
		}
		newTicketOrderResponse := map[string]interface{}{
			"id":             getTicketOrder.ID,
			"order_name":     getTrain.Name,
			"type":           "Kereta",
			"booking_number": ticketTravelerDetail.BoardingTicketCode,
			"price":          ticketTravelerDetail.TrainPrice,
			"created_at":     getTicketOrder.CreatedAt,
			"updated_at":     getTicketOrder.UpdatedAt,
			"deleted_at":     getTicketOrder.DeletedAt,
		}
		newOrderResponses = append(newOrderResponses, newTicketOrderResponse)
	}

	for _, hotelOrder := range newHotelOrder {

		getHotelOrder, err := u.hotelOrderRepo.GetHotelOrderByID(hotelOrder.ID, 1)
		if err != nil {
			return dashboardResponse, err
		}

		getHotel, err := u.hotelRepo.GetHotelByID2(uint(getHotelOrder.HotelID))
		if err != nil {
			return dashboardResponse, err
		}
		newHotelOrderResponse := map[string]interface{}{
			"id":             getHotelOrder.ID,
			"order_name":     getHotel.Name,
			"type":           "Hotel",
			"booking_number": getHotelOrder.HotelOrderCode,
			"price":          getHotelOrder.Price,
			"created_at":     getHotelOrder.CreatedAt,
			"updated_at":     getHotelOrder.UpdatedAt,
			"deleted_at":     getHotelOrder.DeletedAt,
		}
		newOrderResponses = append(newOrderResponses, newHotelOrderResponse)
	}

	var newUserResponses []map[string]interface{}
	for _, user := range newUser {
		newUserResponse := map[string]interface{}{
			"id":              user.ID,
			"full_name":       user.FullName,
			"profile_picture": user.ProfilePicture,
			"created_at":      user.CreatedAt,
			"updated_at":      user.UpdatedAt,
			"deleted_at":      user.DeletedAt,
		}
		newUserResponses = append(newUserResponses, newUserResponse)
	}

	// Sorting newOrderResponses berdasarkan created_at descending menggunakan bubble sort
	for i := 0; i < len(newOrderResponses)-1; i++ {
		for j := 0; j < len(newOrderResponses)-i-1; j++ {
			timeI := newOrderResponses[j]["created_at"].(time.Time)
			timeJ := newOrderResponses[j+1]["created_at"].(time.Time)
			if timeI.Before(timeJ) {
				newOrderResponses[j], newOrderResponses[j+1] = newOrderResponses[j+1], newOrderResponses[j]
			}
		}
	}

	// Batasi menjadi 10 data terbaru
	var limitedOrderResponses []map[string]interface{}
	if len(newOrderResponses) > 10 {
		limitedOrderResponses = newOrderResponses[:10]
	} else {
		limitedOrderResponses = newOrderResponses
	}

	dashboardResponse = dtos.DashboardResponse{
		CountUser: echo.Map{
			"total_user":       countUser,
			"total_user_today": countUserToday,
		},
		CountTrain: echo.Map{
			"total_train":       countTrain,
			"total_train_today": countTrainToday,
		},
		CountOrder: echo.Map{
			"total_order":       countTicketOrder + countHotelOrder,
			"total_order_today": countTicketOrderToday + countHotelOrderToday,
		},
		NewOrder: limitedOrderResponses,
		NewUser:  newUserResponses,
	}
	return dashboardResponse, nil
}
