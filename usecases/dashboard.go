package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/repositories"

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
}

func NewDashboardUsecase(dashboardRepository repositories.DashboardRepository, userRepo repositories.UserRepository, ticketOrderRepo repositories.TicketOrderRepository, ticketTravelerDetailRepo repositories.TicketTravelerDetailRepository, travelerDetailRepo repositories.TravelerDetailRepository, trainCarriageRepo repositories.TrainCarriageRepository, trainRepo repositories.TrainRepository, trainSeatRepo repositories.TrainSeatRepository, stationRepo repositories.StationRepository, trainStationRepo repositories.TrainStationRepository, paymentRepo repositories.PaymentRepository) DashboardUsecase {
	return &dashboardUsecase{dashboardRepository, userRepo, ticketOrderRepo, ticketTravelerDetailRepo, travelerDetailRepo, trainCarriageRepo, trainRepo, trainSeatRepo, stationRepo, trainStationRepo, paymentRepo}
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
	countUser, countUserToday, countTrain, countTrainToday, countTicketOrder, countTicketOrderToday, newTicketOrder, newUser, err := u.dashboardRepository.DashboardGetAll()
	if err != nil {
		return dashboardResponse, err
	}
	var newTicketOrderResponses []map[string]interface{}
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
		newTicketOrderResponses = append(newTicketOrderResponses, newTicketOrderResponse)
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
			"total_order":       countTicketOrder,
			"total_order_today": countTicketOrderToday,
		},
		NewOrder: newTicketOrderResponses,
		NewUser:  newUserResponses,
	}
	return dashboardResponse, nil
}
