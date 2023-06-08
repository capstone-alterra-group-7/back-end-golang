package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/models"
	"back-end-golang/repositories"
	"errors"
)

type HotelRoomUsecase interface {
	// admin
	GetAllHotelRooms(page, limit int) ([]dtos.HotelRoomResponse, int, error)
	GetHotelRoomByID(id uint) (dtos.HotelRoomResponse, error)
	CreateHotelRoom(roomInput *dtos.HotelRoomInput) (dtos.HotelRoomResponse, error)
	UpdateHotelRoom(id uint, roomInput dtos.HotelRoomInput) (dtos.HotelRoomResponse, error)
	DeleteHotelRoom(id uint) error
}

type hotelRoomUsecase struct {
	hotelRepo               repositories.HotelRepository
	hotelRoomRepo           repositories.HotelRoomRepository
	hotelRoomImageRepo      repositories.HotelRoomImageRepository
	hotelRoomFacilitiesRepo repositories.HotelRoomFacilitiesRepository
}

func NewHotelRoomUsecase(hotelRepo repositories.HotelRepository, hotelRoomRepo repositories.HotelRoomRepository, hotelRoomImageRepo repositories.HotelRoomImageRepository, hotelRoomFacilitiesRepo repositories.HotelRoomFacilitiesRepository) HotelRoomUsecase {
	return &hotelRoomUsecase{hotelRepo, hotelRoomRepo, hotelRoomImageRepo, hotelRoomFacilitiesRepo}
}

// =============================== ADMIN ================================== \\

// GetAllHotelRooms godoc
// @Summary      Get all hotel room
// @Description  Get all hotel room
// @Tags         Admin - Hotel Room
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllHotelRoomStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel-room [get]
func (u *hotelRoomUsecase) GetAllHotelRooms(page, limit int) ([]dtos.HotelRoomResponse, int, error) {

	rooms, count, err := u.hotelRoomRepo.GetAllHotelRooms(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var hotelRoomResponses []dtos.HotelRoomResponse

	for _, room := range rooms {
		getImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(room.ID)
		if err != nil {
			return hotelRoomResponses, 0, err
		}
		getFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByID(room.ID)
		if err != nil {
			return hotelRoomResponses, 0, err
		}

		var hotelRoomImageResponses []dtos.HotelRoomImageResponse
		for _, image := range getImage {
			hotelImageResponse := dtos.HotelRoomImageResponse{
				HotelID:     image.HotelID,
				HotelRoomID: image.HotelRoomID,
				ImageUrl:    image.ImageUrl,
			}
			hotelRoomImageResponses = append(hotelRoomImageResponses, hotelImageResponse)
		}

		var hotelFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
		for _, facilities := range getFacilities {
			HotelFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
				HotelID:     facilities.HotelID,
				HotelRoomID: facilities.HotelRoomID,
				Name:        facilities.Name,
			}
			hotelFacilitiesResponses = append(hotelFacilitiesResponses, HotelFacilitiesResponse)
		}

		hotelRoomResponse := dtos.HotelRoomResponse{
			HotelRoomID:       room.ID,
			HotelID:           room.HotelID,
			Name:              room.Name,
			SizeOfRoom:        room.SizeOfRoom,
			QuantityOfRoom:    room.QuantityOfRoom,
			Description:       room.Description,
			NormalPrice:       room.NormalPrice,
			Discount:          room.Discount,
			DiscountPrice:     room.DiscountPrice,
			NumberOfGuest:     room.NumberOfGuest,
			MattressSize:      room.MattressSize,
			NumberOfMattress:  room.NumberOfMattress,
			HotelRoomImage:    hotelRoomImageResponses,
			HotelRoomFacility: hotelFacilitiesResponses,
			CreatedAt:         room.CreatedAt,
			UpdatedAt:         room.UpdatedAt,
		}
		hotelRoomResponses = append(hotelRoomResponses, hotelRoomResponse)
	}
	// Apply offset and limit to hotelRoomResponses
	start := (page - 1) * limit
	end := start + limit

	// Ensure that `start` is within the range of hotelRoomResponses
	if start >= count {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of hotelRoomResponses
	if end > count {
		end = count
	}

	subsetHotelRoomResponses := hotelRoomResponses[start:end]

	return subsetHotelRoomResponses, count, nil
}

// GetHotelRoomByID godoc
// @Summary      Get hotel room by ID
// @Description  Get hotel room by ID
// @Tags         Admin - Hotel Room
// @Accept       json
// @Produce      json
// @Param id path integer true "ID Hotel Room"
// @Success      200 {object} dtos.HotelRoomStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel-room/{id} [get]
func (u *hotelRoomUsecase) GetHotelRoomByID(id uint) (dtos.HotelRoomResponse, error) {
	var hotelRoomResponses dtos.HotelRoomResponse
	room, err := u.hotelRoomRepo.GetHotelRoomByID(id)
	if err != nil {
		return hotelRoomResponses, err
	}

	getImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(room.ID)
	if err != nil {
		return hotelRoomResponses, err
	}
	getFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByID(room.ID)
	if err != nil {
		return hotelRoomResponses, err
	}

	var hotelRoomImageResponses []dtos.HotelRoomImageResponse
	for _, image := range getImage {
		hotelRoomImageResponse := dtos.HotelRoomImageResponse{
			HotelID:     image.HotelID,
			HotelRoomID: image.HotelRoomID,
			ImageUrl:    image.ImageUrl,
		}
		hotelRoomImageResponses = append(hotelRoomImageResponses, hotelRoomImageResponse)
	}

	var hotelRoomFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
	for _, facilities := range getFacilities {
		HotelRoomFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
			HotelID:     facilities.HotelID,
			HotelRoomID: facilities.HotelRoomID,
			Name:        facilities.Name,
		}
		hotelRoomFacilitiesResponses = append(hotelRoomFacilitiesResponses, HotelRoomFacilitiesResponse)
	}

	hotelRoomResponse := dtos.HotelRoomResponse{
		HotelRoomID:       room.ID,
		HotelID:           room.HotelID,
		Name:              room.Name,
		SizeOfRoom:        room.SizeOfRoom,
		QuantityOfRoom:    room.QuantityOfRoom,
		Description:       room.Description,
		NormalPrice:       room.NormalPrice,
		Discount:          room.Discount,
		DiscountPrice:     room.DiscountPrice,
		NumberOfGuest:     room.NumberOfGuest,
		MattressSize:      room.MattressSize,
		NumberOfMattress:  room.NumberOfMattress,
		HotelRoomImage:    hotelRoomImageResponses,
		HotelRoomFacility: hotelRoomFacilitiesResponses,
		CreatedAt:         room.CreatedAt,
		UpdatedAt:         room.UpdatedAt,
	}
	return hotelRoomResponse, nil
}

// CreateHotelRoom godoc
// @Summary      Create a new hotel room
// @Description  Create a new hotel room
// @Tags         Admin - Hotel Room
// @Accept       json
// @Produce      json
// @Param        request body dtos.HotelRoomInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.HotelRoomCreeatedResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel-room [post]
// @Security BearerAuth
func (u *hotelRoomUsecase) CreateHotelRoom(roomInput *dtos.HotelRoomInput) (dtos.HotelRoomResponse, error) {
	var hotelRoomResponse dtos.HotelRoomResponse
	if roomInput.Name == "" || roomInput.Description == "" || roomInput.MattressSize == "" || roomInput.HotelRoomImage == nil || roomInput.HotelRoomFacility == nil || roomInput.HotelID < 1 {
		return hotelRoomResponse, errors.New("failed to create hotel room")
	}

	getHotel, err := u.hotelRepo.GetHotelByID(roomInput.HotelID)
	if err != nil {
		return hotelRoomResponse, errors.New("failed to create hotel room, hotel_id not found")
	}

	createHotelRoom := models.HotelRoom{
		HotelID:          getHotel.ID,
		Name:             roomInput.Name,
		SizeOfRoom:       roomInput.SizeOfRoom,
		QuantityOfRoom:   roomInput.QuantityOfRoom,
		Description:      roomInput.Description,
		NormalPrice:      roomInput.NormalPrice,
		Discount:         roomInput.Discount,
		DiscountPrice:    roomInput.DiscountPrice,
		NumberOfGuest:    roomInput.NumberOfGuest,
		MattressSize:     roomInput.MattressSize,
		NumberOfMattress: roomInput.NumberOfMattress,
	}

	createdHotelRoom, err := u.hotelRoomRepo.CreateHotelRoom(createHotelRoom)
	if err != nil {
		return hotelRoomResponse, err
	}

	for _, roomImage := range roomInput.HotelRoomImage {
		if roomImage.ImageUrl == "" {
			return hotelRoomResponse, errors.New("failed to create hotel room ")
		}
		hotelRoomImagee := models.HotelRoomImage{
			HotelID:     createdHotelRoom.HotelID,
			HotelRoomID: createdHotelRoom.ID,
			ImageUrl:    roomImage.ImageUrl,
		}
		_, err = u.hotelRoomImageRepo.CreateHotelRoomImage(hotelRoomImagee)
		if err != nil {
			return hotelRoomResponse, err
		}
	}

	for _, hotelRoomFacilities := range roomInput.HotelRoomFacility {
		if hotelRoomFacilities.Name == "" {
			return hotelRoomResponse, errors.New("failed to create hotel room")
		}

		hotelRoomFacilitiess := models.HotelRoomFacilities{
			HotelID:     createdHotelRoom.HotelID,
			HotelRoomID: createdHotelRoom.ID,
			Name:        hotelRoomFacilities.Name,
		}
		_, err = u.hotelRoomFacilitiesRepo.CreateHotelRoomFacilities(hotelRoomFacilitiess)
		if err != nil {
			return hotelRoomResponse, err
		}
	}

	getImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(createdHotelRoom.ID)
	if err != nil {
		return hotelRoomResponse, err
	}

	getFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByID(createdHotelRoom.ID)
	if err != nil {
		return hotelRoomResponse, err
	}

	var hotelRoomImageResponses []dtos.HotelRoomImageResponse
	for _, image := range getImage {
		hotelImageResponse := dtos.HotelRoomImageResponse{
			HotelID:     image.HotelID,
			HotelRoomID: image.HotelRoomID,
			ImageUrl:    image.ImageUrl,
		}
		hotelRoomImageResponses = append(hotelRoomImageResponses, hotelImageResponse)
	}

	var hotelFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
	for _, facilities := range getFacilities {
		HotelFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
			HotelID:     facilities.HotelID,
			HotelRoomID: facilities.HotelRoomID,
			Name:        facilities.Name,
		}
		hotelFacilitiesResponses = append(hotelFacilitiesResponses, HotelFacilitiesResponse)
	}

	hotelRoomResponse = dtos.HotelRoomResponse{
		HotelRoomID:       createdHotelRoom.ID,
		HotelID:           createdHotelRoom.HotelID,
		Name:              createdHotelRoom.Name,
		SizeOfRoom:        createdHotelRoom.SizeOfRoom,
		QuantityOfRoom:    createdHotelRoom.QuantityOfRoom,
		Description:       createdHotelRoom.Description,
		NormalPrice:       createdHotelRoom.NormalPrice,
		Discount:          createdHotelRoom.Discount,
		DiscountPrice:     createdHotelRoom.DiscountPrice,
		NumberOfGuest:     createdHotelRoom.NumberOfGuest,
		MattressSize:      createdHotelRoom.MattressSize,
		NumberOfMattress:  createdHotelRoom.NumberOfMattress,
		HotelRoomImage:    hotelRoomImageResponses,
		HotelRoomFacility: hotelFacilitiesResponses,
		CreatedAt:         createdHotelRoom.CreatedAt,
		UpdatedAt:         createdHotelRoom.UpdatedAt,
	}
	return hotelRoomResponse, nil
}

// UpdateHotelRoom godoc
// @Summary      Update hotel room
// @Description  Update hotel room
// @Tags         Admin - Hotel Room
// @Accept       json
// @Produce      json
// @Param id path integer true "ID Hotel Room"
// @Param        request body dtos.HotelRoomInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.HotelRoomStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel-room/{id} [put]
// @Security BearerAuth
func (u *hotelRoomUsecase) UpdateHotelRoom(id uint, roomInput dtos.HotelRoomInput) (dtos.HotelRoomResponse, error) {
	var hotelRooms models.HotelRoom
	var hotelRoomResponse dtos.HotelRoomResponse

	if roomInput.Name == "" || roomInput.Description == "" || roomInput.MattressSize == "" || roomInput.HotelRoomImage == nil || roomInput.HotelRoomFacility == nil || roomInput.HotelID < 1 {
		return hotelRoomResponse, errors.New("failed to update hotel room")
	}

	hotelRooms, err := u.hotelRoomRepo.GetHotelRoomByID(id)
	if err != nil {
		return hotelRoomResponse, err
	}

	getHotel, err := u.hotelRepo.GetHotelByID(roomInput.HotelID)
	if err != nil {
		return hotelRoomResponse, errors.New("failed to update hotel room, hotel_id not found")
	}

	hotelRooms.HotelID = getHotel.ID
	hotelRooms.Name = roomInput.Name
	hotelRooms.SizeOfRoom = roomInput.SizeOfRoom
	hotelRooms.QuantityOfRoom = roomInput.QuantityOfRoom
	hotelRooms.Description = roomInput.Description
	hotelRooms.NormalPrice = roomInput.NormalPrice
	hotelRooms.Discount = roomInput.Discount
	hotelRooms.DiscountPrice = roomInput.DiscountPrice
	hotelRooms.NumberOfGuest = roomInput.NumberOfGuest
	hotelRooms.MattressSize = roomInput.MattressSize
	hotelRooms.NumberOfMattress = roomInput.NumberOfMattress

	updatedHotelRoom, err := u.hotelRoomRepo.UpdateHotelRoom(hotelRooms)
	if err != nil {
		return hotelRoomResponse, err
	}

	u.hotelRoomImageRepo.DeleteHotelRoomImage(id)
	u.hotelRoomFacilitiesRepo.DeleteHotelRoomFacilities(id)

	for _, hotelRoomImage := range roomInput.HotelRoomImage {
		if hotelRoomImage.ImageUrl == "" {
			return hotelRoomResponse, errors.New("failed to update hotel room")
		}
		hotelRoomImagee := models.HotelRoomImage{
			HotelID:     updatedHotelRoom.HotelID,
			HotelRoomID: updatedHotelRoom.ID,
			ImageUrl:    hotelRoomImage.ImageUrl,
		}
		_, err = u.hotelRoomImageRepo.UpdateHotelRoomImage(hotelRoomImagee)
		if err != nil {
			return hotelRoomResponse, err
		}
	}

	for _, hotelRoomFacilities := range roomInput.HotelRoomFacility {
		if hotelRoomFacilities.Name == "" {
			return hotelRoomResponse, errors.New("failed to update hotel room")
		}

		hotelRoomFacilitiess := models.HotelRoomFacilities{
			HotelID:     updatedHotelRoom.HotelID,
			HotelRoomID: updatedHotelRoom.ID,
			Name:        hotelRoomFacilities.Name,
		}
		_, err = u.hotelRoomFacilitiesRepo.UpdateHotelRoomFacilities(hotelRoomFacilitiess)
		if err != nil {
			return hotelRoomResponse, err
		}
	}

	getImage, err := u.hotelRoomImageRepo.GetAllHotelRoomImageByID(updatedHotelRoom.ID)
	if err != nil {
		return hotelRoomResponse, err
	}

	getFacilities, err := u.hotelRoomFacilitiesRepo.GetAllHotelRoomFacilitiesByID(updatedHotelRoom.ID)
	if err != nil {
		return hotelRoomResponse, err
	}

	var hotelRoomImageResponses []dtos.HotelRoomImageResponse
	for _, image := range getImage {
		hotelRoomImageResponse := dtos.HotelRoomImageResponse{
			HotelID:     image.HotelID,
			HotelRoomID: image.HotelRoomID,
			ImageUrl:    image.ImageUrl,
		}
		hotelRoomImageResponses = append(hotelRoomImageResponses, hotelRoomImageResponse)
	}

	var hotelRoomFacilitiesResponses []dtos.HotelRoomFacilitiesResponse
	for _, facilities := range getFacilities {
		HotelRoomFacilitiesResponse := dtos.HotelRoomFacilitiesResponse{
			HotelID:     facilities.HotelID,
			HotelRoomID: facilities.HotelRoomID,
			Name:        facilities.Name,
		}
		hotelRoomFacilitiesResponses = append(hotelRoomFacilitiesResponses, HotelRoomFacilitiesResponse)
	}

	hotelRoomResponse = dtos.HotelRoomResponse{
		HotelRoomID:       updatedHotelRoom.ID,
		HotelID:           updatedHotelRoom.HotelID,
		Name:              updatedHotelRoom.Name,
		SizeOfRoom:        updatedHotelRoom.SizeOfRoom,
		QuantityOfRoom:    updatedHotelRoom.QuantityOfRoom,
		Description:       updatedHotelRoom.Description,
		NormalPrice:       updatedHotelRoom.NormalPrice,
		Discount:          updatedHotelRoom.Discount,
		DiscountPrice:     updatedHotelRoom.DiscountPrice,
		NumberOfGuest:     updatedHotelRoom.NumberOfGuest,
		MattressSize:      updatedHotelRoom.MattressSize,
		NumberOfMattress:  updatedHotelRoom.NumberOfMattress,
		HotelRoomImage:    hotelRoomImageResponses,
		HotelRoomFacility: hotelRoomFacilitiesResponses,
		CreatedAt:         updatedHotelRoom.CreatedAt,
		UpdatedAt:         updatedHotelRoom.UpdatedAt,
	}
	return hotelRoomResponse, nil
}

// DeleteHotelRoom godoc
// @Summary      Delete a hotel room
// @Description  Delete a hotel room
// @Tags         Admin - Hotel Room
// @Accept       json
// @Produce      json
// @Param id path integer true "ID Hotel Room"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/hotel-room/{id} [delete]
// @Security BearerAuth
func (u *hotelRoomUsecase) DeleteHotelRoom(id uint) error {
	u.hotelRoomImageRepo.DeleteHotelRoomImage(id)
	u.hotelRoomFacilitiesRepo.DeleteHotelRoomFacilities(id)
	return u.hotelRoomRepo.DeleteHotelRoom(id)
}
