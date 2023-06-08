package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/models"
	"back-end-golang/repositories"
	"errors"
)

type HotelUsecase interface {
	// admin
	GetAllHotels(page, limit int) ([]dtos.HotelResponse, int, error)
	GetHotelByID(id uint) (dtos.HotelResponse, error)
	CreateHotel(train *dtos.HotelInput) (dtos.HotelResponse, error)
	UpdateHotel(id uint, hotelInput dtos.HotelInput) (dtos.HotelResponse, error)
	DeleteHotel(id uint) error
}

type hotelUsecase struct {
	hotelRepo           repositories.HotelRepository
	hotelImageRepo      repositories.HotelImageRepository
	hotelFacilitiesRepo repositories.HotelFacilitiesRepository
	hotelPoliciesRepo   repositories.HotelPoliciesRepository
}

func NewHotelUsecase(hotelRepo repositories.HotelRepository, hotelImageRepo repositories.HotelImageRepository, hotelFacilitiesRepo repositories.HotelFacilitiesRepository, hotelPoliciesRepo repositories.HotelPoliciesRepository) HotelUsecase {
	return &hotelUsecase{hotelRepo, hotelImageRepo, hotelFacilitiesRepo, hotelPoliciesRepo}
}

// =============================== ADMIN ================================== \\

// GetAllHotels godoc
// @Summary      Get all hotel
// @Description  Get all hotel
// @Tags         Admin - Hotel
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllHotelStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel [get]
func (u *hotelUsecase) GetAllHotels(page, limit int) ([]dtos.HotelResponse, int, error) {

	hotels, count, err := u.hotelRepo.GetAllHotels(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var hotelResponses []dtos.HotelResponse

	for _, hotel := range hotels {
		getImage, err := u.hotelImageRepo.GetAllHotelImageByID(hotel.ID)
		if err != nil {
			return hotelResponses, 0, err
		}
		getFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(hotel.ID)
		if err != nil {
			return hotelResponses, 0, err
		}

		getPolicy, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(hotel.ID)
		if err != nil {
			return hotelResponses, 0, err
		}

		var hotelImageResponses []dtos.HotelImageResponse
		for _, image := range getImage {
			hotelImageResponse := dtos.HotelImageResponse{
				HotelID:  image.HotelID,
				ImageUrl: image.ImageUrl,
			}
			hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
		}

		var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
		for _, facilities := range getFacilities {
			HotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
				HotelID: facilities.HotelID,
				Name:    facilities.Name,
			}
			hotelFacilitiesResponses = append(hotelFacilitiesResponses, HotelFacilitiesResponse)
		}

		hotelPoliciesResponses := dtos.HotelPoliciesResponse{
			HotelID:            getPolicy.HotelID,
			IsCheckInCheckOut:  getPolicy.IsCheckInCheckOut,
			TimeCheckIn:        getPolicy.TimeCheckIn,
			TimeCheckOut:       getPolicy.TimeCheckOut,
			IsPolicyCanceled:   getPolicy.IsPolicyCanceled,
			PolicyMinimumAge:   getPolicy.PolicyMinimumAge,
			IsCheckInEarly:     getPolicy.IsCheckInEarly,
			IsCheckOutOverdue:  getPolicy.IsCheckOutOverdue,
			IsBreakfast:        getPolicy.IsBreakfast,
			TimeBreakfastStart: getPolicy.TimeBreakfastStart,
			TimeBreakfastEnd:   getPolicy.TimeBreakfastEnd,
			IsSmoking:          getPolicy.IsSmoking,
			IsPet:              getPolicy.IsPet,
		}

		hotelResponse := dtos.HotelResponse{
			HotelID:         hotel.ID,
			Name:            hotel.Name,
			Class:           hotel.Class,
			Description:     hotel.Description,
			PhoneNumber:     hotel.PhoneNumber,
			Email:           hotel.Email,
			Address:         hotel.Address,
			HotelImage:      hotelImageResponses,
			HotelFacilities: hotelFacilitiesResponses,
			HotelPolicy:     hotelPoliciesResponses,
			CreatedAt:       hotel.CreatedAt,
			UpdatedAt:       hotel.UpdatedAt,
		}
		hotelResponses = append(hotelResponses, hotelResponse)
	}
	// Apply offset and limit to hotelResponses
	start := (page - 1) * limit
	end := start + limit

	// Ensure that `start` is within the range of hotelResponses
	if start >= count {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of hotelResponses
	if end > count {
		end = count
	}

	subsetHotelResponses := hotelResponses[start:end]

	return subsetHotelResponses, count, nil
}

// GetHotelByID godoc
// @Summary      Get hotel by ID
// @Description  Get hotel by ID
// @Tags         Admin - Hotel
// @Accept       json
// @Produce      json
// @Param id path integer true "ID Hotel"
// @Success      200 {object} dtos.HotelStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel/{id} [get]
func (u *hotelUsecase) GetHotelByID(id uint) (dtos.HotelResponse, error) {
	var hotelResponses dtos.HotelResponse
	hotel, err := u.hotelRepo.GetHotelByID(id)
	if err != nil {
		return hotelResponses, err
	}

	getImage, err := u.hotelImageRepo.GetAllHotelImageByID(hotel.ID)
	if err != nil {
		return hotelResponses, err
	}
	getFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(hotel.ID)
	if err != nil {
		return hotelResponses, err
	}

	getPolicy, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(hotel.ID)
	if err != nil {
		return hotelResponses, err
	}

	var hotelImageResponses []dtos.HotelImageResponse
	for _, image := range getImage {
		hotelImageResponse := dtos.HotelImageResponse{
			HotelID:  image.HotelID,
			ImageUrl: image.ImageUrl,
		}
		hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
	}

	var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
	for _, facilities := range getFacilities {
		HotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
			HotelID: facilities.HotelID,
			Name:    facilities.Name,
		}
		hotelFacilitiesResponses = append(hotelFacilitiesResponses, HotelFacilitiesResponse)
	}

	hotelPoliciesResponses := dtos.HotelPoliciesResponse{
		HotelID:            getPolicy.HotelID,
		IsCheckInCheckOut:  getPolicy.IsCheckInCheckOut,
		TimeCheckIn:        getPolicy.TimeCheckIn,
		TimeCheckOut:       getPolicy.TimeCheckOut,
		IsPolicyCanceled:   getPolicy.IsPolicyCanceled,
		PolicyMinimumAge:   getPolicy.PolicyMinimumAge,
		IsCheckInEarly:     getPolicy.IsCheckInEarly,
		IsCheckOutOverdue:  getPolicy.IsCheckOutOverdue,
		IsBreakfast:        getPolicy.IsBreakfast,
		TimeBreakfastStart: getPolicy.TimeBreakfastStart,
		TimeBreakfastEnd:   getPolicy.TimeBreakfastEnd,
		IsSmoking:          getPolicy.IsSmoking,
		IsPet:              getPolicy.IsPet,
	}

	hotelResponse := dtos.HotelResponse{
		HotelID:         hotel.ID,
		Name:            hotel.Name,
		Class:           hotel.Class,
		Description:     hotel.Description,
		PhoneNumber:     hotel.PhoneNumber,
		Email:           hotel.Email,
		Address:         hotel.Address,
		HotelImage:      hotelImageResponses,
		HotelFacilities: hotelFacilitiesResponses,
		HotelPolicy:     hotelPoliciesResponses,
		CreatedAt:       hotel.CreatedAt,
		UpdatedAt:       hotel.UpdatedAt,
	}
	return hotelResponse, nil
}

// CreateHotel godoc
// @Summary      Create a new hotel
// @Description  Create a new hotel
// @Tags         Admin - Hotel
// @Accept       json
// @Produce      json
// @Param        request body dtos.HotelInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.HotelCreeatedResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel [post]
// @Security BearerAuth
func (u *hotelUsecase) CreateHotel(hotel *dtos.HotelInput) (dtos.HotelResponse, error) {
	var hotelResponse dtos.HotelResponse
	if hotel.Name == "" || hotel.Email == "" || hotel.Address == "" || hotel.PhoneNumber == "" || hotel.Class == "" || hotel.Description == "" || hotel.HotelFacilities == nil || hotel.HotelImage == nil || hotel.HotelPolicy == nil {
		return hotelResponse, errors.New("failed to create hotel")
	}

	createHotel := models.Hotel{
		Name:        hotel.Name,
		Class:       hotel.Class,
		Description: hotel.Description,
		PhoneNumber: hotel.PhoneNumber,
		Email:       hotel.Email,
		Address:     hotel.Address,
	}

	createdHotel, err := u.hotelRepo.CreateHotel(createHotel)
	if err != nil {
		return hotelResponse, err
	}

	for _, hotelImage := range hotel.HotelImage {
		if hotelImage.ImageUrl == "" {
			return hotelResponse, errors.New("failed to create hotel")
		}
		hotelImagee := models.HotelImage{
			HotelID:  createdHotel.ID,
			ImageUrl: hotelImage.ImageUrl,
		}
		_, err = u.hotelImageRepo.CreateHotelImage(hotelImagee)
		if err != nil {
			return hotelResponse, err
		}
	}

	for _, hotelFacilities := range hotel.HotelFacilities {
		if hotelFacilities.Name == "" {
			return hotelResponse, errors.New("failed to create hotel")
		}

		hotelFacilitiess := models.HotelFacilities{
			HotelID: createdHotel.ID,
			Name:    hotelFacilities.Name,
		}
		_, err = u.hotelFacilitiesRepo.CreateHotelFacilities(hotelFacilitiess)
		if err != nil {
			return hotelResponse, err
		}
	}

	for _, hotelPolicy := range hotel.HotelPolicy {
		if hotelPolicy.TimeCheckIn == "" || hotelPolicy.TimeCheckOut == "" || hotelPolicy.TimeBreakfastStart == "" || hotelPolicy.TimeBreakfastEnd == "" {
			return hotelResponse, errors.New("failed to create hotel")
		}
		hotelPoliciess := models.HotelPolicies{
			HotelID:            createdHotel.ID,
			IsCheckInCheckOut:  hotelPolicy.IsCheckInCheckOut,
			TimeCheckIn:        hotelPolicy.TimeCheckIn,
			TimeCheckOut:       hotelPolicy.TimeCheckOut,
			IsPolicyCanceled:   hotelPolicy.IsPolicyCanceled,
			PolicyMinimumAge:   hotelPolicy.PolicyMinimumAge,
			IsCheckInEarly:     hotelPolicy.IsCheckInEarly,
			IsCheckOutOverdue:  hotelPolicy.IsCheckOutOverdue,
			IsBreakfast:        hotelPolicy.IsBreakfast,
			TimeBreakfastStart: hotelPolicy.TimeBreakfastStart,
			TimeBreakfastEnd:   hotelPolicy.TimeBreakfastEnd,
			IsSmoking:          hotelPolicy.IsSmoking,
			IsPet:              hotelPolicy.IsPet,
		}
		_, err = u.hotelPoliciesRepo.CreateHotelPolicies(hotelPoliciess)
		if err != nil {
			return hotelResponse, err
		}
	}

	getImage, err := u.hotelImageRepo.GetAllHotelImageByID(createdHotel.ID)
	if err != nil {
		return hotelResponse, err
	}

	getFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(createdHotel.ID)
	if err != nil {
		return hotelResponse, err
	}

	getPolicy, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(createdHotel.ID)
	if err != nil {
		return hotelResponse, err
	}

	var hotelImageResponses []dtos.HotelImageResponse
	for _, image := range getImage {
		hotelImageResponse := dtos.HotelImageResponse{
			HotelID:  image.HotelID,
			ImageUrl: image.ImageUrl,
		}
		hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
	}

	var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
	for _, facilities := range getFacilities {
		HotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
			HotelID: facilities.HotelID,
			Name:    facilities.Name,
		}
		hotelFacilitiesResponses = append(hotelFacilitiesResponses, HotelFacilitiesResponse)
	}

	hotelPoliciesResponses := dtos.HotelPoliciesResponse{
		HotelID:            getPolicy.HotelID,
		IsCheckInCheckOut:  getPolicy.IsCheckInCheckOut,
		TimeCheckIn:        getPolicy.TimeCheckIn,
		TimeCheckOut:       getPolicy.TimeCheckOut,
		IsPolicyCanceled:   getPolicy.IsPolicyCanceled,
		PolicyMinimumAge:   getPolicy.PolicyMinimumAge,
		IsCheckInEarly:     getPolicy.IsCheckInEarly,
		IsCheckOutOverdue:  getPolicy.IsCheckOutOverdue,
		IsBreakfast:        getPolicy.IsBreakfast,
		TimeBreakfastStart: getPolicy.TimeBreakfastStart,
		TimeBreakfastEnd:   getPolicy.TimeBreakfastEnd,
		IsSmoking:          getPolicy.IsSmoking,
		IsPet:              getPolicy.IsPet,
	}

	hotelResponse = dtos.HotelResponse{
		HotelID:         createdHotel.ID,
		Name:            createdHotel.Name,
		Class:           createdHotel.Class,
		Description:     createdHotel.Description,
		PhoneNumber:     createdHotel.PhoneNumber,
		Email:           createdHotel.Email,
		Address:         createdHotel.Address,
		HotelImage:      hotelImageResponses,
		HotelFacilities: hotelFacilitiesResponses,
		HotelPolicy:     hotelPoliciesResponses,
		CreatedAt:       createdHotel.CreatedAt,
		UpdatedAt:       createdHotel.UpdatedAt,
	}
	return hotelResponse, nil
}

// UpdateHotel godoc
// @Summary      Update hotel
// @Description  Update hotel
// @Tags         Admin - Hotel
// @Accept       json
// @Produce      json
// @Param id path integer true "ID hotel"
// @Param        request body dtos.HotelInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.HotelStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel/{id} [put]
// @Security BearerAuth
func (u *hotelUsecase) UpdateHotel(id uint, hotel dtos.HotelInput) (dtos.HotelResponse, error) {
	var hotels models.Hotel
	var hotelResponse dtos.HotelResponse

	if hotel.Name == "" || hotel.Email == "" || hotel.Address == "" || hotel.PhoneNumber == "" || hotel.Class == "" || hotel.Description == "" || hotel.HotelFacilities == nil || hotel.HotelImage == nil || hotel.HotelPolicy == nil {
		return hotelResponse, errors.New("failed to update hotel")
	}

	hotels, err := u.hotelRepo.GetHotelByID(id)
	if err != nil {
		return hotelResponse, err
	}

	hotels.Name = hotel.Name
	hotels.Class = hotel.Class
	hotels.Description = hotel.Description
	hotels.PhoneNumber = hotel.PhoneNumber
	hotels.Email = hotel.Email
	hotels.Address = hotel.Address

	updatedHotel, err := u.hotelRepo.UpdateHotel(hotels)
	if err != nil {
		return hotelResponse, err
	}

	u.hotelImageRepo.DeleteHotelImage(id)
	u.hotelFacilitiesRepo.DeleteHotelFacilities(id)
	u.hotelPoliciesRepo.DeleteHotelPolicies(id)

	for _, hotelImage := range hotel.HotelImage {
		if hotelImage.ImageUrl == "" {
			return hotelResponse, errors.New("failed to update hotel")
		}
		hotelImagee := models.HotelImage{
			HotelID:  updatedHotel.ID,
			ImageUrl: hotelImage.ImageUrl,
		}
		_, err = u.hotelImageRepo.UpdateHotelImage(hotelImagee)
		if err != nil {
			return hotelResponse, err
		}
	}

	for _, hotelFacilities := range hotel.HotelFacilities {
		if hotelFacilities.Name == "" {
			return hotelResponse, errors.New("failed to update hotel")
		}

		hotelFacilitiess := models.HotelFacilities{
			HotelID: updatedHotel.ID,
			Name:    hotelFacilities.Name,
		}
		_, err = u.hotelFacilitiesRepo.UpdateHotelFacilities(hotelFacilitiess)
		if err != nil {
			return hotelResponse, err
		}
	}

	for _, hotelPolicy := range hotel.HotelPolicy {
		if hotelPolicy.TimeCheckIn == "" || hotelPolicy.TimeCheckOut == "" || hotelPolicy.TimeBreakfastStart == "" || hotelPolicy.TimeBreakfastEnd == "" {
			return hotelResponse, errors.New("failed to update hotel")
		}
		hotelPoliciess := models.HotelPolicies{
			HotelID:            updatedHotel.ID,
			IsCheckInCheckOut:  hotelPolicy.IsCheckInCheckOut,
			TimeCheckIn:        hotelPolicy.TimeCheckIn,
			TimeCheckOut:       hotelPolicy.TimeCheckOut,
			IsPolicyCanceled:   hotelPolicy.IsPolicyCanceled,
			PolicyMinimumAge:   hotelPolicy.PolicyMinimumAge,
			IsCheckInEarly:     hotelPolicy.IsCheckInEarly,
			IsCheckOutOverdue:  hotelPolicy.IsCheckOutOverdue,
			IsBreakfast:        hotelPolicy.IsBreakfast,
			TimeBreakfastStart: hotelPolicy.TimeBreakfastStart,
			TimeBreakfastEnd:   hotelPolicy.TimeBreakfastEnd,
			IsSmoking:          hotelPolicy.IsSmoking,
			IsPet:              hotelPolicy.IsPet,
		}
		_, err = u.hotelPoliciesRepo.UpdateHotelPolicies(hotelPoliciess)
		if err != nil {
			return hotelResponse, err
		}
	}

	getImage, err := u.hotelImageRepo.GetAllHotelImageByID(updatedHotel.ID)
	if err != nil {
		return hotelResponse, err
	}

	getFacilities, err := u.hotelFacilitiesRepo.GetAllHotelFacilitiesByID(updatedHotel.ID)
	if err != nil {
		return hotelResponse, err
	}

	getPolicy, err := u.hotelPoliciesRepo.GetHotelPoliciesByIDHotel(updatedHotel.ID)
	if err != nil {
		return hotelResponse, err
	}

	var hotelImageResponses []dtos.HotelImageResponse
	for _, image := range getImage {
		hotelImageResponse := dtos.HotelImageResponse{
			HotelID:  image.HotelID,
			ImageUrl: image.ImageUrl,
		}
		hotelImageResponses = append(hotelImageResponses, hotelImageResponse)
	}

	var hotelFacilitiesResponses []dtos.HotelFacilitiesResponse
	for _, facilities := range getFacilities {
		HotelFacilitiesResponse := dtos.HotelFacilitiesResponse{
			HotelID: facilities.HotelID,
			Name:    facilities.Name,
		}
		hotelFacilitiesResponses = append(hotelFacilitiesResponses, HotelFacilitiesResponse)
	}

	hotelPoliciesResponses := dtos.HotelPoliciesResponse{
		HotelID:            getPolicy.HotelID,
		IsCheckInCheckOut:  getPolicy.IsCheckInCheckOut,
		TimeCheckIn:        getPolicy.TimeCheckIn,
		TimeCheckOut:       getPolicy.TimeCheckOut,
		IsPolicyCanceled:   getPolicy.IsPolicyCanceled,
		PolicyMinimumAge:   getPolicy.PolicyMinimumAge,
		IsCheckInEarly:     getPolicy.IsCheckInEarly,
		IsCheckOutOverdue:  getPolicy.IsCheckOutOverdue,
		IsBreakfast:        getPolicy.IsBreakfast,
		TimeBreakfastStart: getPolicy.TimeBreakfastStart,
		TimeBreakfastEnd:   getPolicy.TimeBreakfastEnd,
		IsSmoking:          getPolicy.IsSmoking,
		IsPet:              getPolicy.IsPet,
	}

	hotelResponse = dtos.HotelResponse{
		HotelID:         updatedHotel.ID,
		Name:            updatedHotel.Name,
		Class:           updatedHotel.Class,
		Description:     updatedHotel.Description,
		PhoneNumber:     updatedHotel.PhoneNumber,
		Email:           updatedHotel.Email,
		Address:         updatedHotel.Address,
		HotelImage:      hotelImageResponses,
		HotelFacilities: hotelFacilitiesResponses,
		HotelPolicy:     hotelPoliciesResponses,
		CreatedAt:       updatedHotel.CreatedAt,
		UpdatedAt:       updatedHotel.UpdatedAt,
	}
	return hotelResponse, nil
}

// DeleteHotel godoc
// @Summary      Delete a hotel
// @Description  Delete a hotel
// @Tags         Admin - Hotel
// @Accept       json
// @Produce      json
// @Param id path integer true "ID Hotel"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel/{id} [delete]
// @Security BearerAuth
func (u *hotelUsecase) DeleteHotel(id uint) error {
	u.hotelImageRepo.DeleteHotelImage(id)
	u.hotelFacilitiesRepo.DeleteHotelFacilities(id)
	u.hotelPoliciesRepo.DeleteHotelPolicies(id)
	return u.hotelRepo.DeleteHotel(id)
}
