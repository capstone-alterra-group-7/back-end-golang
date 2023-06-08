package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/models"
	"back-end-golang/repositories"
	"errors"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type TicketOrderUsecase interface {
	GetTicketOrders(page, limit int, userID uint, status string) ([]dtos.TicketTravelerDetailOrderResponse, int, error)
	GetTicketOrdersByAdmin(page, limit int, search, dateStart, dateEnd, orderBy, filter string) ([]dtos.TicketTravelerDetailOrderResponse, int, error)
	GetTicketOrdersDetailByAdmin(ticketOrderId, trainId uint) (dtos.TicketTravelerDetailOrderResponse, error)
	GetTicketOrderByID(userID, ticketTravelerDetailId, ticketOrderId uint) (dtos.TicketTravelerDetailOrderResponse, error)
	CreateTicketOrder(userID uint, ticketOrderInput dtos.TicketOrderInput) (dtos.TicketOrderResponse, error)
	UpdateTicketOrder(userID, ticketOrderID uint, status string) (dtos.TicketOrderResponse, error)
}

type ticketOrderUsecase struct {
	ticketOrderRepo          repositories.TicketOrderRepository
	ticketTravelerDetailRepo repositories.TicketTravelerDetailRepository
	travelerDetailRepo       repositories.TravelerDetailRepository
	trainCarriageRepo        repositories.TrainCarriageRepository
	trainRepo                repositories.TrainRepository
	trainSeatRepo            repositories.TrainSeatRepository
	stationRepo              repositories.StationRepository
	trainStationRepo         repositories.TrainStationRepository
	paymentRepo              repositories.PaymentRepository
	userRepo                 repositories.UserRepository
}

func NewTicketOrderUsecase(ticketOrderRepo repositories.TicketOrderRepository, ticketTravelerDetailRepo repositories.TicketTravelerDetailRepository, travelerDetailRepo repositories.TravelerDetailRepository, trainCarriageRepo repositories.TrainCarriageRepository, trainRepo repositories.TrainRepository, trainSeatRepo repositories.TrainSeatRepository, stationRepo repositories.StationRepository, trainStationRepo repositories.TrainStationRepository, paymentRepo repositories.PaymentRepository, userRepo repositories.UserRepository) TicketOrderUsecase {
	return &ticketOrderUsecase{ticketOrderRepo, ticketTravelerDetailRepo, travelerDetailRepo, trainCarriageRepo, trainRepo, trainSeatRepo, stationRepo, trainStationRepo, paymentRepo, userRepo}
}

// GetTicketOrders godoc
// @Summary      Get Ticket Order User
// @Description  Get Ticket Order User
// @Tags         User - Order
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param status query string false "Filter by status order"
// @Success      200 {object} dtos.GetAllTicketTravelerDetailOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/order/ticket [get]
// @Security BearerAuth
func (u *ticketOrderUsecase) GetTicketOrders(page, limit int, userID uint, status string) ([]dtos.TicketTravelerDetailOrderResponse, int, error) {
	var ticketTravelerDetailResponses []dtos.TicketTravelerDetailOrderResponse

	ticketTravelerDetail, _, err := u.ticketTravelerDetailRepo.GetAllTicketTravelerDetails()
	if err != nil {
		return ticketTravelerDetailResponses, 0, err
	}
	visitedIDs := make(map[uint]map[uint]bool)

	for _, ticketTravelerDetail := range ticketTravelerDetail {
		var getTicketOrder models.TicketOrder
		if status == "" {
			getTicketOrder, err = u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, userID)
			if err != nil {
				return ticketTravelerDetailResponses, 0, err
			}
		} else {
			getTicketOrder, err = u.ticketOrderRepo.GetTicketOrderByStatusAndID(ticketTravelerDetail.TicketOrderID, userID, status)
			if err != nil {
				return ticketTravelerDetailResponses, 0, err
			}
		}

		trainID := ticketTravelerDetail.TrainID
		ticketOrderID := ticketTravelerDetail.TicketOrderID

		if visitedIDs[ticketOrderID] == nil {
			visitedIDs[ticketOrderID] = make(map[uint]bool)
		}
		if visitedIDs[ticketOrderID][trainID] {
			continue
		}
		visitedIDs[ticketOrderID][trainID] = true

		getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetail.TrainID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getTrainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(uint(ticketTravelerDetail.TrainCarriageID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getTrainSeat, err := u.trainSeatRepo.GetTrainSeatByID(uint(ticketTravelerDetail.TrainSeatID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getStationOrigin, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationOriginID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getStationDestination, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationDestinationID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}

		getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByID2(ticketTravelerDetail.TicketOrderID)
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}

		var travelerDetailResponses []dtos.TravelerDetailResponse

		for _, travelerDetail := range getTravelerDetail {
			travelerDetailResponse := dtos.TravelerDetailResponse{
				ID:           int(travelerDetail.ID),
				Title:        travelerDetail.Title,
				FullName:     travelerDetail.FullName,
				IDCardNumber: *travelerDetail.IDCardNumber,
			}
			travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
		}

		getOrderTicket, err := u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, userID)
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getPayment, err := u.paymentRepo.GetPaymentByID(uint(getOrderTicket.PaymentID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}

		ticketTravelerDetailResponse := dtos.TicketTravelerDetailOrderResponse{
			TicketOrderID:    int(getTicketOrder.ID),
			QuantityAdult:    getTicketOrder.QuantityAdult,
			QuantityInfant:   getTicketOrder.QuantityInfant,
			NameOrder:        getTicketOrder.NameOrder,
			EmailOrder:       getTicketOrder.EmailOrder,
			PhoneNumberOrder: getTicketOrder.PhoneNumberOrder,
			TicketOrderCode:  getTicketOrder.TicketOrderCode,
			Payment: dtos.PaymentResponses{
				ID:            int(getPayment.ID),
				Type:          getPayment.Type,
				ImageUrl:      getPayment.ImageUrl,
				Name:          getPayment.Name,
				AccountName:   getPayment.AccountName,
				AccountNumber: getPayment.AccountNumber,
			},
			TravelerDetail: travelerDetailResponses,
			Train: dtos.TrainResponsesSimply{
				TrainID:         getTrain.ID,
				CodeTrain:       getTrain.CodeTrain,
				Name:            getTrain.Name,
				Class:           getTrainCarriage.Class,
				TrainPrice:      ticketTravelerDetail.TrainPrice,
				TrainCarriageID: getTrainCarriage.ID,
				TrainCarriage:   getTrainCarriage.Name,
				TrainSeatID:     getTrainSeat.ID,
				TrainSeat:       getTrainSeat.Name,
			},
			StationOrigin: dtos.StationResponseSimply{
				StationID:  getStationOrigin.ID,
				Origin:     getStationOrigin.Origin,
				Name:       getStationOrigin.Name,
				Initial:    getStationOrigin.Initial,
				ArriveTime: ticketTravelerDetail.DepartureTime,
			},
			StationDestination: dtos.StationResponseSimply{
				StationID:  getStationDestination.ID,
				Origin:     getStationDestination.Origin,
				Name:       getStationDestination.Name,
				Initial:    getStationDestination.Initial,
				ArriveTime: ticketTravelerDetail.ArrivalTime,
			},
			Date:               ticketTravelerDetail.DateOfDeparture,
			BoardingTicketCode: ticketTravelerDetail.BoardingTicketCode,
			Status:             getTicketOrder.Status,
			CreatedAt:          getOrderTicket.CreatedAt,
			UpdatedAt:          getOrderTicket.UpdatedAt,
		}
		ticketTravelerDetailResponses = append(ticketTravelerDetailResponses, ticketTravelerDetailResponse)
	}

	// Apply offset and limit to ticketTravelerDetailResponses
	start := (page - 1) * limit
	end := start + limit

	// Ensure that `start` is within the range of ticketTravelerDetailResponses
	if start >= len(ticketTravelerDetailResponses) {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of ticketTravelerDetailResponses
	if end > len(ticketTravelerDetailResponses) {
		end = len(ticketTravelerDetailResponses)
	}

	subsetTicketTravelerDetailResponses := ticketTravelerDetailResponses[start:end]

	return subsetTicketTravelerDetailResponses, len(ticketTravelerDetailResponses), nil
}

// GetTicketOrdersByAdmin godoc
// @Summary      Get Ticket Order User
// @Description  Get Ticket Order User
// @Tags         Admin - Order
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "search train name"
// @Param date_start query string false "Date start"
// @Param date_end query string false "Date end"
// @Param order_by query string false "Order by name"
// @Param filter query string false "Filter by status order"
// @Success      200 {object} dtos.GetAllTicketTravelerDetailOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/order/ticket [get]
// @Security BearerAuth
func (u *ticketOrderUsecase) GetTicketOrdersByAdmin(page, limit int, search, dateStart, dateEnd, orderBy, filter string) ([]dtos.TicketTravelerDetailOrderResponse, int, error) {
	var ticketTravelerDetailResponses []dtos.TicketTravelerDetailOrderResponse

	ticketTravelerDetail, _, err := u.ticketTravelerDetailRepo.GetAllTicketTravelerDetails()
	if err != nil {
		return ticketTravelerDetailResponses, 0, err
	}
	visitedIDs := make(map[uint]map[uint]bool)

	for _, ticketTravelerDetail := range ticketTravelerDetail {
		var getTicketOrder models.TicketOrder
		if filter == "" {
			getTicketOrder, err = u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, 1)
			if err != nil {
				return ticketTravelerDetailResponses, 0, err
			}
		} else {
			getTicketOrder, err = u.ticketOrderRepo.GetTicketOrderByStatusAndID(ticketTravelerDetail.TicketOrderID, 1, filter)
			if err != nil {
				return ticketTravelerDetailResponses, 0, err
			}
		}

		trainID := ticketTravelerDetail.TrainID
		ticketOrderID := ticketTravelerDetail.TicketOrderID

		if visitedIDs[ticketOrderID] == nil {
			visitedIDs[ticketOrderID] = make(map[uint]bool)
		}
		if visitedIDs[ticketOrderID][trainID] {
			continue
		}
		visitedIDs[ticketOrderID][trainID] = true

		getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetail.TrainID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getTrainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(uint(ticketTravelerDetail.TrainCarriageID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getTrainSeat, err := u.trainSeatRepo.GetTrainSeatByID(uint(ticketTravelerDetail.TrainSeatID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getStationOrigin, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationOriginID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getStationDestination, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationDestinationID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}

		getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByID2(ticketTravelerDetail.TicketOrderID)
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}

		// Filter by search (getTrain.Name)
		if search != "" && !strings.Contains(strings.ToLower(getTrain.Name), strings.ToLower(search)) {
			continue
		}
		var travelerDetailResponses []dtos.TravelerDetailResponse

		for _, travelerDetail := range getTravelerDetail {

			travelerDetailResponse := dtos.TravelerDetailResponse{
				ID:           int(travelerDetail.ID),
				Title:        travelerDetail.Title,
				FullName:     travelerDetail.FullName,
				IDCardNumber: *travelerDetail.IDCardNumber,
			}
			travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
		}

		getOrderTicket, err := u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, 1)
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}
		getPayment, err := u.paymentRepo.GetPaymentByID(uint(getOrderTicket.PaymentID))
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}

		getUser, err := u.userRepo.UserGetById2(getOrderTicket.UserID)
		if err != nil {
			return ticketTravelerDetailResponses, 0, err
		}

		ticketTravelerDetailResponse := dtos.TicketTravelerDetailOrderResponse{
			TicketOrderID:    int(getTicketOrder.ID),
			QuantityAdult:    getTicketOrder.QuantityAdult,
			QuantityInfant:   getTicketOrder.QuantityInfant,
			NameOrder:        getTicketOrder.NameOrder,
			EmailOrder:       getTicketOrder.EmailOrder,
			PhoneNumberOrder: getTicketOrder.PhoneNumberOrder,
			TicketOrderCode:  getTicketOrder.TicketOrderCode,
			User: &dtos.UserInformationResponses{
				ID:          getUser.ID,
				FullName:    getUser.FullName,
				Email:       getUser.Email,
				PhoneNumber: getUser.PhoneNumber,
				BirthDate:   helpers.FormatDateToYMD(getUser.BirthDate),
				Citizen:     getUser.Citizen,
			},
			Payment: dtos.PaymentResponses{
				ID:            int(getPayment.ID),
				Type:          getPayment.Type,
				ImageUrl:      getPayment.ImageUrl,
				Name:          getPayment.Name,
				AccountName:   getPayment.AccountName,
				AccountNumber: getPayment.AccountNumber,
			},
			TravelerDetail: travelerDetailResponses,
			Train: dtos.TrainResponsesSimply{
				TrainID:         getTrain.ID,
				CodeTrain:       getTrain.CodeTrain,
				Name:            getTrain.Name,
				Class:           getTrainCarriage.Class,
				TrainPrice:      ticketTravelerDetail.TrainPrice,
				TrainCarriageID: getTrainCarriage.ID,
				TrainCarriage:   getTrainCarriage.Name,
				TrainSeatID:     getTrainSeat.ID,
				TrainSeat:       getTrainSeat.Name,
			},
			StationOrigin: dtos.StationResponseSimply{
				StationID:  getStationOrigin.ID,
				Origin:     getStationOrigin.Origin,
				Name:       getStationOrigin.Name,
				Initial:    getStationOrigin.Initial,
				ArriveTime: ticketTravelerDetail.DepartureTime,
			},
			StationDestination: dtos.StationResponseSimply{
				StationID:  getStationDestination.ID,
				Origin:     getStationDestination.Origin,
				Name:       getStationDestination.Name,
				Initial:    getStationDestination.Initial,
				ArriveTime: ticketTravelerDetail.ArrivalTime,
			},
			Date:               ticketTravelerDetail.DateOfDeparture,
			BoardingTicketCode: ticketTravelerDetail.BoardingTicketCode,
			Status:             getTicketOrder.Status,
			CreatedAt:          getOrderTicket.CreatedAt,
			UpdatedAt:          getOrderTicket.UpdatedAt,
		}
		ticketTravelerDetailResponses = append(ticketTravelerDetailResponses, ticketTravelerDetailResponse)
	}

	// Apply offset and limit to ticketTravelerDetailResponses
	start := (page - 1) * limit
	end := start + limit

	// Ensure that `start` is within the range of ticketTravelerDetailResponses
	if start >= len(ticketTravelerDetailResponses) {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of ticketTravelerDetailResponses
	if end > len(ticketTravelerDetailResponses) {
		end = len(ticketTravelerDetailResponses)
	}

	// Filter by dateStart and dateEnd
	filteredTicketTravelerDetailResponses := make([]dtos.TicketTravelerDetailOrderResponse, 0)
	if dateStart == "" && dateEnd == "" {
		filteredTicketTravelerDetailResponses = ticketTravelerDetailResponses
	} else {
		for _, ticketTravelerDetailResponse := range ticketTravelerDetailResponses {
			if helpers.FormatDateToYMD(&ticketTravelerDetailResponse.Date) >= dateStart && helpers.FormatDateToYMD(&ticketTravelerDetailResponse.Date) <= dateEnd {
				filteredTicketTravelerDetailResponses = append(filteredTicketTravelerDetailResponses, ticketTravelerDetailResponse)
			}
		}
	}

	// Sort by orderBy (getTrain.Name)
	switch orderBy {
	case "asc":
		sort.Slice(filteredTicketTravelerDetailResponses, func(i, j int) bool {
			return filteredTicketTravelerDetailResponses[i].Train.Name < filteredTicketTravelerDetailResponses[j].Train.Name
		})
	case "desc":
		sort.Slice(filteredTicketTravelerDetailResponses, func(i, j int) bool {
			return filteredTicketTravelerDetailResponses[i].Train.Name > filteredTicketTravelerDetailResponses[j].Train.Name
		})
	}

	// Apply offset and limit to filteredTicketTravelerDetailResponses
	start = (page - 1) * limit
	end = start + limit

	// Ensure that `start` is within the range of filteredTicketTravelerDetailResponses
	if start >= len(filteredTicketTravelerDetailResponses) {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of filteredTicketTravelerDetailResponses
	if end > len(filteredTicketTravelerDetailResponses) {
		end = len(filteredTicketTravelerDetailResponses)
	}

	subsetTicketTravelerDetailResponses := filteredTicketTravelerDetailResponses[start:end]

	return subsetTicketTravelerDetailResponses, len(filteredTicketTravelerDetailResponses), nil
}

// GetTicketOrdersByAdmin godoc
// @Summary      Get Ticket Order User
// @Description  Get Ticket Order User
// @Tags         Admin - Order
// @Accept       json
// @Produce      json
// @Param ticket_order_id query int true "Ticket Order ID"
// @Param train_id query int true "Train ID"
// @Success      200 {object} dtos.TicketTravelerDetailOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/order/ticket/detail [get]
// @Security BearerAuth
func (u *ticketOrderUsecase) GetTicketOrdersDetailByAdmin(ticketOrderId, trainId uint) (dtos.TicketTravelerDetailOrderResponse, error) {
	var ticketTravelerDetailResponses dtos.TicketTravelerDetailOrderResponse

	ticketTravelerDetail, err := u.ticketTravelerDetailRepo.GetTicketTravelerDetailByTicketOrderIDAndTrainID(ticketOrderId, trainId)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	// visitedIDs := make(map[uint]map[uint]bool)

	// for _, ticketTravelerDetail := range ticketTravelerDetail {
	getTicketOrder, err := u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, 1)
	// if filter == "" {
	// 	if err != nil {
	// 		return ticketTravelerDetailResponses, err
	// 	}
	// } else {
	// 	getTicketOrder, err = u.ticketOrderRepo.GetTicketOrderByStatusAndID(ticketTravelerDetail.TicketOrderID, 1, filter)
	// 	if err != nil {
	// 		return ticketTravelerDetailResponses, err
	// 	}
	// }

	// trainID := ticketTravelerDetail.TrainID
	// ticketOrderID := ticketTravelerDetail.TicketOrderID

	// if visitedIDs[ticketOrderID] == nil {
	// 	visitedIDs[ticketOrderID] = make(map[uint]bool)
	// }
	// if visitedIDs[ticketOrderID][trainID] {
	// 	continue
	// }
	// visitedIDs[ticketOrderID][trainID] = true

	getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetail.TrainID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getTrainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(uint(ticketTravelerDetail.TrainCarriageID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getTrainSeat, err := u.trainSeatRepo.GetTrainSeatByID(uint(ticketTravelerDetail.TrainSeatID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getStationOrigin, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationOriginID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getStationDestination, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationDestinationID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByID2(ticketTravelerDetail.TicketOrderID)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	// Filter by search (getTrain.Name)
	// if search != "" && !strings.Contains(strings.ToLower(getTrain.Name), strings.ToLower(search)) {
	// 	continue
	// }
	var travelerDetailResponses []dtos.TravelerDetailResponse

	for _, travelerDetail := range getTravelerDetail {

		travelerDetailResponse := dtos.TravelerDetailResponse{
			ID:           int(travelerDetail.ID),
			Title:        travelerDetail.Title,
			FullName:     travelerDetail.FullName,
			IDCardNumber: *travelerDetail.IDCardNumber,
		}
		travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
	}

	getOrderTicket, err := u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, 1)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getPayment, err := u.paymentRepo.GetPaymentByID(uint(getOrderTicket.PaymentID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	getUser, err := u.userRepo.UserGetById2(getOrderTicket.UserID)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	ticketTravelerDetailResponses = dtos.TicketTravelerDetailOrderResponse{
		TicketOrderID:    int(getTicketOrder.ID),
		QuantityAdult:    getTicketOrder.QuantityAdult,
		QuantityInfant:   getTicketOrder.QuantityInfant,
		NameOrder:        getTicketOrder.NameOrder,
		EmailOrder:       getTicketOrder.EmailOrder,
		PhoneNumberOrder: getTicketOrder.PhoneNumberOrder,
		TicketOrderCode:  getTicketOrder.TicketOrderCode,
		User: &dtos.UserInformationResponses{
			ID:          getUser.ID,
			FullName:    getUser.FullName,
			Email:       getUser.Email,
			PhoneNumber: getUser.PhoneNumber,
			BirthDate:   helpers.FormatDateToYMD(getUser.BirthDate),
			Citizen:     getUser.Citizen,
		},
		Payment: dtos.PaymentResponses{
			ID:            int(getPayment.ID),
			Type:          getPayment.Type,
			ImageUrl:      getPayment.ImageUrl,
			Name:          getPayment.Name,
			AccountName:   getPayment.AccountName,
			AccountNumber: getPayment.AccountNumber,
		},
		TravelerDetail: travelerDetailResponses,
		Train: dtos.TrainResponsesSimply{
			TrainID:         getTrain.ID,
			CodeTrain:       getTrain.CodeTrain,
			Name:            getTrain.Name,
			Class:           getTrainCarriage.Class,
			TrainPrice:      ticketTravelerDetail.TrainPrice,
			TrainCarriageID: getTrainCarriage.ID,
			TrainCarriage:   getTrainCarriage.Name,
			TrainSeatID:     getTrainSeat.ID,
			TrainSeat:       getTrainSeat.Name,
		},
		StationOrigin: dtos.StationResponseSimply{
			StationID:  getStationOrigin.ID,
			Origin:     getStationOrigin.Origin,
			Name:       getStationOrigin.Name,
			Initial:    getStationOrigin.Initial,
			ArriveTime: ticketTravelerDetail.DepartureTime,
		},
		StationDestination: dtos.StationResponseSimply{
			StationID:  getStationDestination.ID,
			Origin:     getStationDestination.Origin,
			Name:       getStationDestination.Name,
			Initial:    getStationDestination.Initial,
			ArriveTime: ticketTravelerDetail.ArrivalTime,
		},
		Date:               ticketTravelerDetail.DateOfDeparture,
		BoardingTicketCode: ticketTravelerDetail.BoardingTicketCode,
		Status:             getTicketOrder.Status,
		CreatedAt:          getOrderTicket.CreatedAt,
		UpdatedAt:          getOrderTicket.UpdatedAt,
	}
	// ticketTravelerDetailResponses = append(ticketTravelerDetailResponses, ticketTravelerDetailResponse)
	// }

	// Apply offset and limit to ticketTravelerDetailResponses
	// start := (page - 1) * limit
	// end := start + limit

	// Ensure that `start` is within the range of ticketTravelerDetailResponses
	// if start >= len(ticketTravelerDetailResponses) {
	// 	return nil, 0, nil
	// }

	// Ensure that `end` does not exceed the length of ticketTravelerDetailResponses
	// if end > len(ticketTravelerDetailResponses) {
	// 	end = len(ticketTravelerDetailResponses)
	// }

	// Filter by dateStart and dateEnd
	// filteredTicketTravelerDetailResponses := make([]dtos.TicketTravelerDetailOrderResponse, 0)
	// if dateStart == "" && dateEnd == "" {
	// 	filteredTicketTravelerDetailResponses = ticketTravelerDetailResponses
	// } else {
	// 	for _, ticketTravelerDetailResponse := range ticketTravelerDetailResponses {
	// 		if helpers.FormatDateToYMD(&ticketTravelerDetailResponse.Date) >= dateStart && helpers.FormatDateToYMD(&ticketTravelerDetailResponse.Date) <= dateEnd {
	// 			filteredTicketTravelerDetailResponses = append(filteredTicketTravelerDetailResponses, ticketTravelerDetailResponse)
	// 		}
	// 	}
	// }

	// Sort by orderBy (getTrain.Name)
	// switch orderBy {
	// case "asc":
	// 	sort.Slice(filteredTicketTravelerDetailResponses, func(i, j int) bool {
	// 		return filteredTicketTravelerDetailResponses[i].Train.Name < filteredTicketTravelerDetailResponses[j].Train.Name
	// 	})
	// case "desc":
	// 	sort.Slice(filteredTicketTravelerDetailResponses, func(i, j int) bool {
	// 		return filteredTicketTravelerDetailResponses[i].Train.Name > filteredTicketTravelerDetailResponses[j].Train.Name
	// 	})
	// }

	// Apply offset and limit to filteredTicketTravelerDetailResponses
	// start = (page - 1) * limit
	// end = start + limit

	// Ensure that `start` is within the range of filteredTicketTravelerDetailResponses
	// if start >= len(filteredTicketTravelerDetailResponses) {
	// 	return nil, 0, nil
	// }

	// Ensure that `end` does not exceed the length of filteredTicketTravelerDetailResponses
	// if end > len(filteredTicketTravelerDetailResponses) {
	// 	end = len(filteredTicketTravelerDetailResponses)
	// }

	// subsetTicketTravelerDetailResponses := filteredTicketTravelerDetailResponses[start:end]

	return ticketTravelerDetailResponses, nil
}

// GetTicketOrderByID godoc
// @Summary      Get Ticket Order User by ID
// @Description  Get Ticket Order User by ID
// @Tags         User - Order
// @Accept       json
// @Produce      json
// @Param ticket_order_id query int true "Ticket Order ID"
// @Param train_id query int true "Train ID"
// @Success      200 {object} dtos.TicketTravelerDetailOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/order/ticket/detail [get]
// @Security BearerAuth
func (u *ticketOrderUsecase) GetTicketOrderByID(userID, ticketOrderId, trainId uint) (dtos.TicketTravelerDetailOrderResponse, error) {
	var ticketTravelerDetailResponses dtos.TicketTravelerDetailOrderResponse

	ticketTravelerDetail, err := u.ticketTravelerDetailRepo.GetTicketTravelerDetailByTicketOrderIDAndTrainID(ticketOrderId, trainId)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getTicketOrder, err := u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, userID)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetail.TrainID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getTrainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(uint(ticketTravelerDetail.TrainCarriageID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getTrainSeat, err := u.trainSeatRepo.GetTrainSeatByID(uint(ticketTravelerDetail.TrainSeatID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getStationOrigin, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationOriginID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getStationDestination, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationDestinationID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByID2(ticketTravelerDetail.TicketOrderID)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	var travelerDetailResponses []dtos.TravelerDetailResponse

	for _, travelerDetail := range getTravelerDetail {
		travelerDetailResponse := dtos.TravelerDetailResponse{
			ID:           int(travelerDetail.ID),
			Title:        travelerDetail.Title,
			FullName:     travelerDetail.FullName,
			IDCardNumber: *travelerDetail.IDCardNumber,
		}
		travelerDetailResponses = append(travelerDetailResponses, travelerDetailResponse)
	}

	getOrderTicket, err := u.ticketOrderRepo.GetTicketOrderByID(ticketTravelerDetail.TicketOrderID, userID)
	if err != nil {
		return ticketTravelerDetailResponses, err
	}
	getPayment, err := u.paymentRepo.GetPaymentByID(uint(getOrderTicket.PaymentID))
	if err != nil {
		return ticketTravelerDetailResponses, err
	}

	ticketTravelerDetailResponses = dtos.TicketTravelerDetailOrderResponse{
		TicketOrderID:    int(getTicketOrder.ID),
		QuantityAdult:    getTicketOrder.QuantityAdult,
		QuantityInfant:   getTicketOrder.QuantityInfant,
		NameOrder:        getTicketOrder.NameOrder,
		EmailOrder:       getTicketOrder.EmailOrder,
		PhoneNumberOrder: getTicketOrder.PhoneNumberOrder,
		TicketOrderCode:  getTicketOrder.TicketOrderCode,
		Payment: dtos.PaymentResponses{
			ID:            int(getPayment.ID),
			Type:          getPayment.Type,
			ImageUrl:      getPayment.ImageUrl,
			Name:          getPayment.Name,
			AccountName:   getPayment.AccountName,
			AccountNumber: getPayment.AccountNumber,
		},
		TravelerDetail: travelerDetailResponses,
		Train: dtos.TrainResponsesSimply{
			TrainID:         getTrain.ID,
			CodeTrain:       getTrain.CodeTrain,
			Name:            getTrain.Name,
			Class:           getTrainCarriage.Class,
			TrainPrice:      ticketTravelerDetail.TrainPrice,
			TrainCarriageID: getTrainCarriage.ID,
			TrainCarriage:   getTrainCarriage.Name,
			TrainSeatID:     getTrainSeat.ID,
			TrainSeat:       getTrainSeat.Name,
		},
		StationOrigin: dtos.StationResponseSimply{
			StationID:  getStationOrigin.ID,
			Origin:     getStationOrigin.Origin,
			Name:       getStationOrigin.Name,
			Initial:    getStationOrigin.Initial,
			ArriveTime: ticketTravelerDetail.DepartureTime,
		},
		StationDestination: dtos.StationResponseSimply{
			StationID:  getStationDestination.ID,
			Origin:     getStationDestination.Origin,
			Name:       getStationDestination.Name,
			Initial:    getStationDestination.Initial,
			ArriveTime: ticketTravelerDetail.ArrivalTime,
		},
		Date:               ticketTravelerDetail.DateOfDeparture,
		BoardingTicketCode: ticketTravelerDetail.BoardingTicketCode,
		Status:             getTicketOrder.Status,
		CreatedAt:          getOrderTicket.CreatedAt,
		UpdatedAt:          getOrderTicket.UpdatedAt,
	}

	return ticketTravelerDetailResponses, nil
}

// CreateTicketOrder godoc
// @Summary      Order ticket KA
// @Description  Order ticket KA
// @Tags         User - Train
// @Accept       json
// @Produce      json
// @Param        request body dtos.TicketOrderInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.TicketOrderCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/train/order [post]
// @Security BearerAuth
func (u *ticketOrderUsecase) CreateTicketOrder(userID uint, ticketOrderInput dtos.TicketOrderInput) (dtos.TicketOrderResponse, error) {
	var ticketOrderResponse dtos.TicketOrderResponse
	sumTrainPrice := 0
	trainPrice := 0
	if ticketOrderInput.QuantityAdult < 1 || ticketOrderInput.PaymentID < 1 || ticketOrderInput.NameOrder == "" || ticketOrderInput.EmailOrder == "" || ticketOrderInput.PhoneNumberOrder == "" || ticketOrderInput.TravelerDetail == nil || ticketOrderInput.TicketTravelerDetailDeparture == nil {
		return ticketOrderResponse, errors.New("Failed to create ticket order")
	}
	createTicketOrder := models.TicketOrder{
		UserID:           userID,
		QuantityAdult:    ticketOrderInput.QuantityAdult,
		QuantityInfant:   ticketOrderInput.QuantityInfant,
		Price:            0,
		WithReturn:       ticketOrderInput.WithReturn,
		PaymentID:        ticketOrderInput.PaymentID,
		TotalAmount:      0,
		NameOrder:        ticketOrderInput.NameOrder,
		EmailOrder:       ticketOrderInput.EmailOrder,
		PhoneNumberOrder: ticketOrderInput.PhoneNumberOrder,
		TicketOrderCode:  "ticket-order-" + uuid.New().String(),
		Status:           "unpaid",
	}

	createTicketOrder, err := u.ticketOrderRepo.CreateTicketOrder(createTicketOrder)
	if err != nil {
		return ticketOrderResponse, err
	}

	var ticketTravelerDetailDepartureResponses []dtos.TicketTravelerDetailResponse

	for _, travelerDetail := range ticketOrderInput.TravelerDetail {
		if travelerDetail.Title == "" || travelerDetail.FullName == "" {
			return ticketOrderResponse, errors.New("Failed to create ticket order")
		}
		createTravelerDetail := models.TravelerDetail{
			UserID:        userID,
			TicketOrderID: createTicketOrder.ID,
			Title:         travelerDetail.Title,
			FullName:      travelerDetail.FullName,
			IDCardNumber:  &travelerDetail.IDCardNumber,
		}
		createTravelerDetail, err := u.travelerDetailRepo.CreateTravelerDetail(createTravelerDetail)
		if err != nil {
			return ticketOrderResponse, err
		}

		for _, ticketTravelerDetailDeparture := range ticketOrderInput.TicketTravelerDetailDeparture {
			if ticketTravelerDetailDeparture.TrainID < 1 || ticketTravelerDetailDeparture.TrainCarriageID < 1 || ticketTravelerDetailDeparture.TrainSeatID < 1 || ticketTravelerDetailDeparture.StationOriginID < 1 || ticketTravelerDetailDeparture.StationDestinationID < 1 || ticketTravelerDetailDeparture.Date == "" {
				return ticketOrderResponse, errors.New("Failed to create ticket order")
			}
			dateDepartureParse, err := helpers.FormatStringToDate(ticketTravelerDetailDeparture.Date)
			if err != nil {
				return ticketOrderResponse, errors.New("Failed to parsing date")
			}

			getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetailDeparture.TrainID))
			if err != nil {
				return ticketOrderResponse, err
			}
			getTrainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(uint(ticketTravelerDetailDeparture.TrainCarriageID))
			if err != nil {
				return ticketOrderResponse, err
			}

			if travelerDetail.IDCardNumber != "" {
				trainPrice = getTrainCarriage.Price
				sumTrainPrice += trainPrice
			} else {
				trainPrice = 0
				sumTrainPrice += trainPrice
			}

			getTrainSeat, err := u.trainSeatRepo.GetTrainSeatByID(uint(ticketTravelerDetailDeparture.TrainSeatID))
			if err != nil {
				return ticketOrderResponse, err
			}
			getStationOrigin, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetailDeparture.StationOriginID))
			if err != nil {
				return ticketOrderResponse, err
			}
			getStationDestination, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetailDeparture.StationDestinationID))
			if err != nil {
				return ticketOrderResponse, err
			}

			trainStationOrigin, err := u.trainStationRepo.GetTrainStationByTrainIDStationID(getTrain.ID, getStationOrigin.ID)
			if err != nil {
				return ticketOrderResponse, err
			}
			trainStationDestination, err := u.trainStationRepo.GetTrainStationByTrainIDStationID(getTrain.ID, getStationDestination.ID)
			if err != nil {
				return ticketOrderResponse, err
			}

			createTicketTravelerDetail := models.TicketTravelerDetail{
				TicketOrderID:        createTicketOrder.ID,
				TravelerDetailID:     createTravelerDetail.ID,
				TrainID:              uint(getTrain.ID),
				TrainPrice:           trainPrice,
				TrainCarriageID:      uint(getTrainCarriage.ID),
				TrainSeatID:          uint(getTrainSeat.ID),
				StationOriginID:      uint(getStationOrigin.ID),
				DepartureTime:        trainStationOrigin.ArriveTime,
				StationDestinationID: uint(getStationDestination.ID),
				ArrivalTime:          trainStationDestination.ArriveTime,
				DateOfDeparture:      dateDepartureParse,
				BoardingTicketCode:   "boarding-ticket-" + uuid.New().String(),
			}
			createTicketTravelerDetail, err = u.ticketTravelerDetailRepo.CreateTicketTravelerDetail(createTicketTravelerDetail)
			if err != nil {
				return ticketOrderResponse, err
			}

			getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByID(createTicketTravelerDetail.TravelerDetailID)
			if err != nil {
				return ticketOrderResponse, err
			}

			ticketTravelerDetailResponse := dtos.TicketTravelerDetailResponse{
				TicketTravelerDetailID: int(createTicketTravelerDetail.ID),
				TravelerDetail: dtos.TravelerDetailResponse{
					ID:           int(getTravelerDetail.ID),
					Title:        getTravelerDetail.Title,
					FullName:     getTravelerDetail.FullName,
					IDCardNumber: *getTravelerDetail.IDCardNumber,
				},
				Train: dtos.TrainResponsesSimply{
					TrainID:         getTrain.ID,
					CodeTrain:       getTrain.CodeTrain,
					Name:            getTrain.Name,
					Class:           getTrainCarriage.Class,
					TrainPrice:      trainPrice,
					TrainCarriageID: getTrainCarriage.ID,
					TrainCarriage:   getTrainCarriage.Name,
					TrainSeatID:     getTrainSeat.ID,
					TrainSeat:       getTrainSeat.Name,
				},
				StationOrigin: dtos.StationResponseSimply{
					StationID:  getStationOrigin.ID,
					Origin:     getStationOrigin.Origin,
					Name:       getStationOrigin.Name,
					Initial:    getStationOrigin.Initial,
					ArriveTime: createTicketTravelerDetail.DepartureTime,
				},
				StationDestination: dtos.StationResponseSimply{
					StationID:  getStationDestination.ID,
					Origin:     getStationDestination.Origin,
					Name:       getStationDestination.Name,
					Initial:    getStationDestination.Initial,
					ArriveTime: createTicketTravelerDetail.ArrivalTime,
				},
				Date:               createTicketTravelerDetail.DateOfDeparture,
				BoardingTicketCode: createTicketTravelerDetail.BoardingTicketCode,
			}
			ticketTravelerDetailDepartureResponses = append(ticketTravelerDetailDepartureResponses, ticketTravelerDetailResponse)
		}

		if createTicketOrder.WithReturn {
			for _, ticketTravelerDetailReturn := range ticketOrderInput.TicketTravelerDetailReturn {
				if ticketTravelerDetailReturn.TrainID < 1 || ticketTravelerDetailReturn.TrainCarriageID < 1 || ticketTravelerDetailReturn.TrainSeatID < 1 || ticketTravelerDetailReturn.StationOriginID < 1 || ticketTravelerDetailReturn.StationDestinationID < 1 || ticketTravelerDetailReturn.Date == "" {
					return ticketOrderResponse, errors.New("Failed to create ticket order")
				}
				dateReturn, err := helpers.FormatStringToDate(ticketTravelerDetailReturn.Date)
				if err != nil {
					return ticketOrderResponse, errors.New("Failed to parsing date")
				}

				getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetailReturn.TrainID))
				if err != nil {
					return ticketOrderResponse, err
				}
				getTrainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(uint(ticketTravelerDetailReturn.TrainCarriageID))
				if err != nil {
					return ticketOrderResponse, err
				}

				if travelerDetail.IDCardNumber != "" {
					trainPrice = getTrainCarriage.Price
					sumTrainPrice += trainPrice
				} else {
					trainPrice = 0
					sumTrainPrice += trainPrice
				}

				getTrainSeat, err := u.trainSeatRepo.GetTrainSeatByID(uint(ticketTravelerDetailReturn.TrainSeatID))
				if err != nil {
					return ticketOrderResponse, err
				}
				getStationOrigin, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetailReturn.StationOriginID))
				if err != nil {
					return ticketOrderResponse, err
				}
				getStationDestination, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetailReturn.StationDestinationID))
				if err != nil {
					return ticketOrderResponse, err
				}

				trainStationOrigin, err := u.trainStationRepo.GetTrainStationByTrainIDStationID(getTrain.ID, getStationOrigin.ID)
				if err != nil {
					return ticketOrderResponse, err
				}
				trainStationDestination, err := u.trainStationRepo.GetTrainStationByTrainIDStationID(getTrain.ID, getStationDestination.ID)
				if err != nil {
					return ticketOrderResponse, err
				}

				createTicketTravelerDetail := models.TicketTravelerDetail{
					TicketOrderID:        createTicketOrder.ID,
					TravelerDetailID:     createTravelerDetail.ID,
					TrainID:              uint(getTrain.ID),
					TrainPrice:           trainPrice,
					TrainCarriageID:      uint(getTrainCarriage.ID),
					TrainSeatID:          uint(getTrainSeat.ID),
					StationOriginID:      uint(getStationOrigin.ID),
					DepartureTime:        trainStationOrigin.ArriveTime,
					StationDestinationID: uint(getStationDestination.ID),
					ArrivalTime:          trainStationDestination.ArriveTime,
					DateOfDeparture:      dateReturn,
					BoardingTicketCode:   "boarding-ticket-" + uuid.New().String(),
				}
				createTicketTravelerDetail, err = u.ticketTravelerDetailRepo.CreateTicketTravelerDetail(createTicketTravelerDetail)
				if err != nil {
					return ticketOrderResponse, err
				}

				getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByID(createTicketTravelerDetail.TravelerDetailID)
				if err != nil {
					return ticketOrderResponse, err
				}

				ticketTravelerDetailResponse := dtos.TicketTravelerDetailResponse{
					TicketTravelerDetailID: int(createTicketTravelerDetail.ID),
					TravelerDetail: dtos.TravelerDetailResponse{
						ID:           int(getTravelerDetail.ID),
						Title:        getTravelerDetail.Title,
						FullName:     getTravelerDetail.FullName,
						IDCardNumber: *getTravelerDetail.IDCardNumber,
					},
					Train: dtos.TrainResponsesSimply{
						TrainID:         getTrain.ID,
						CodeTrain:       getTrain.CodeTrain,
						Name:            getTrain.Name,
						Class:           getTrainCarriage.Class,
						TrainPrice:      trainPrice,
						TrainCarriageID: getTrainCarriage.ID,
						TrainCarriage:   getTrainCarriage.Name,
						TrainSeatID:     getTrainSeat.ID,
						TrainSeat:       getTrainSeat.Name,
					},
					StationOrigin: dtos.StationResponseSimply{
						StationID:  getStationOrigin.ID,
						Origin:     getStationOrigin.Origin,
						Name:       getStationOrigin.Name,
						Initial:    getStationOrigin.Initial,
						ArriveTime: createTicketTravelerDetail.DepartureTime,
					},
					StationDestination: dtos.StationResponseSimply{
						StationID:  getStationDestination.ID,
						Origin:     getStationDestination.Origin,
						Name:       getStationDestination.Name,
						Initial:    getStationDestination.Initial,
						ArriveTime: createTicketTravelerDetail.ArrivalTime,
					},
					Date:               createTicketTravelerDetail.DateOfDeparture,
					BoardingTicketCode: createTicketTravelerDetail.BoardingTicketCode,
				}
				ticketTravelerDetailDepartureResponses = append(ticketTravelerDetailDepartureResponses, ticketTravelerDetailResponse)
			}
		}

	}

	createTicketOrder.Price = sumTrainPrice
	createTicketOrder.TotalAmount = sumTrainPrice * ticketOrderInput.QuantityAdult

	updateTicketOrder, err := u.ticketOrderRepo.UpdateTicketOrder(createTicketOrder)
	if err != nil {
		return ticketOrderResponse, err
	}

	getOrderTicket, err := u.ticketOrderRepo.GetTicketOrderByID(updateTicketOrder.ID, userID)
	if err != nil {
		return ticketOrderResponse, err
	}
	getPayment, err := u.paymentRepo.GetPaymentByID(uint(getOrderTicket.PaymentID))
	if err != nil {
		return ticketOrderResponse, err
	}

	ticketOrderResponse = dtos.TicketOrderResponse{
		TicketOrderID:    int(getOrderTicket.ID),
		QuantityAdult:    getOrderTicket.QuantityAdult,
		QuantityInfant:   getOrderTicket.QuantityInfant,
		Price:            getOrderTicket.Price,
		TotalAmount:      getOrderTicket.TotalAmount,
		NameOrder:        getOrderTicket.NameOrder,
		EmailOrder:       getOrderTicket.EmailOrder,
		PhoneNumberOrder: getOrderTicket.PhoneNumberOrder,
		TicketOrderCode:  getOrderTicket.TicketOrderCode,
		Status:           getOrderTicket.Status,
		Payment: dtos.PaymentResponses{
			ID:            int(getPayment.ID),
			Type:          getPayment.Type,
			ImageUrl:      getPayment.ImageUrl,
			Name:          getPayment.Name,
			AccountName:   getPayment.AccountName,
			AccountNumber: getPayment.AccountNumber,
		},
		TicketTravelerDetail: ticketTravelerDetailDepartureResponses,
		CreatedAt:            getOrderTicket.CreatedAt,
		UpdatedAt:            getOrderTicket.UpdatedAt,
	}

	return ticketOrderResponse, nil
}

// UpdateTicketOrder godoc
// @Summary      Update Order ticket KA
// @Description  Update Order ticket KA
// @Tags         User - Train
// @Accept       json
// @Produce      json
// @Param ticket_order_id query int true "Ticket Order ID"
// @Param status query string true "Update Status Order ID"
// @Success      200 {object} dtos.TicketOrderStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /user/train/order [patch]
// @Security BearerAuth
func (u *ticketOrderUsecase) UpdateTicketOrder(userID, ticketOrderID uint, status string) (dtos.TicketOrderResponse, error) {
	var ticketOrderResponse dtos.TicketOrderResponse

	createTicketOrder, err := u.ticketOrderRepo.GetTicketOrderByID(ticketOrderID, userID)
	if err != nil {
		return ticketOrderResponse, err
	}

	createTicketOrder.Status = status

	createTicketOrder, err = u.ticketOrderRepo.UpdateTicketOrder(createTicketOrder)
	if err != nil {
		return ticketOrderResponse, err
	}

	getTicketTravelerDetail, err := u.ticketTravelerDetailRepo.GetTicketTravelerDetailByTicketOrderID(createTicketOrder.ID)
	if err != nil {
		return ticketOrderResponse, err
	}

	var ticketTravelerDetailResponses []dtos.TicketTravelerDetailResponse

	for _, ticketTravelerDetail := range getTicketTravelerDetail {
		getTrain, err := u.trainRepo.GetTrainByID(uint(ticketTravelerDetail.TrainID))
		if err != nil {
			return ticketOrderResponse, err
		}
		getTrainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(uint(ticketTravelerDetail.TrainCarriageID))
		if err != nil {
			return ticketOrderResponse, err
		}
		getTrainSeat, err := u.trainSeatRepo.GetTrainSeatByID(uint(ticketTravelerDetail.TrainSeatID))
		if err != nil {
			return ticketOrderResponse, err
		}
		getStationOrigin, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationOriginID))
		if err != nil {
			return ticketOrderResponse, err
		}
		getStationDestination, err := u.stationRepo.GetStationByID(uint(ticketTravelerDetail.StationDestinationID))
		if err != nil {
			return ticketOrderResponse, err
		}

		getTravelerDetail, err := u.travelerDetailRepo.GetTravelerDetailByID(ticketTravelerDetail.TravelerDetailID)
		if err != nil {
			return ticketOrderResponse, err
		}

		ticketTravelerDetailResponse := dtos.TicketTravelerDetailResponse{
			TicketTravelerDetailID: int(getTravelerDetail.ID),
			TravelerDetail: dtos.TravelerDetailResponse{
				ID:           int(getTravelerDetail.ID),
				Title:        getTravelerDetail.Title,
				FullName:     getTravelerDetail.FullName,
				IDCardNumber: *getTravelerDetail.IDCardNumber,
			},
			Train: dtos.TrainResponsesSimply{
				TrainID:         getTrain.ID,
				CodeTrain:       getTrain.CodeTrain,
				Name:            getTrain.Name,
				Class:           getTrainCarriage.Class,
				TrainPrice:      ticketTravelerDetail.TrainPrice,
				TrainCarriageID: getTrainCarriage.ID,
				TrainCarriage:   getTrainCarriage.Name,
				TrainSeatID:     getTrainSeat.ID,
				TrainSeat:       getTrainSeat.Name,
			},
			StationOrigin: dtos.StationResponseSimply{
				StationID:  getStationOrigin.ID,
				Origin:     getStationOrigin.Origin,
				Name:       getStationOrigin.Name,
				Initial:    getStationOrigin.Initial,
				ArriveTime: ticketTravelerDetail.DepartureTime,
			},
			StationDestination: dtos.StationResponseSimply{
				StationID:  getStationDestination.ID,
				Origin:     getStationDestination.Origin,
				Name:       getStationDestination.Name,
				Initial:    getStationDestination.Initial,
				ArriveTime: ticketTravelerDetail.ArrivalTime,
			},
			Date:               ticketTravelerDetail.DateOfDeparture,
			BoardingTicketCode: ticketTravelerDetail.BoardingTicketCode,
		}
		ticketTravelerDetailResponses = append(ticketTravelerDetailResponses, ticketTravelerDetailResponse)
	}
	getOrderTicket, err := u.ticketOrderRepo.GetTicketOrderByID(createTicketOrder.ID, userID)
	if err != nil {
		return ticketOrderResponse, err
	}
	getPayment, err := u.paymentRepo.GetPaymentByID(uint(getOrderTicket.PaymentID))
	if err != nil {
		return ticketOrderResponse, err
	}

	ticketOrderResponse = dtos.TicketOrderResponse{
		TicketOrderID:    int(getOrderTicket.ID),
		QuantityAdult:    getOrderTicket.QuantityAdult,
		QuantityInfant:   getOrderTicket.QuantityInfant,
		Price:            getOrderTicket.Price,
		TotalAmount:      getOrderTicket.TotalAmount,
		NameOrder:        getOrderTicket.NameOrder,
		EmailOrder:       getOrderTicket.EmailOrder,
		PhoneNumberOrder: getOrderTicket.PhoneNumberOrder,
		TicketOrderCode:  getOrderTicket.TicketOrderCode,
		Status:           getOrderTicket.Status,
		Payment: dtos.PaymentResponses{
			ID:            int(getPayment.ID),
			Type:          getPayment.Type,
			ImageUrl:      getPayment.ImageUrl,
			Name:          getPayment.Name,
			AccountName:   getPayment.AccountName,
			AccountNumber: getPayment.AccountNumber,
		},
		TicketTravelerDetail: ticketTravelerDetailResponses,
		CreatedAt:            getOrderTicket.CreatedAt,
		UpdatedAt:            getOrderTicket.UpdatedAt,
	}

	return ticketOrderResponse, nil
}
