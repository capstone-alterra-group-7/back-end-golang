package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/models"
	"back-end-golang/repositories"
)

type HistorySearchHotelUseCase interface {
	GetHistorySearchHotelByID(userId, id uint) (models.HistorySearchHotel, error)
	GetAllHistorySearchHotels(userID uint, page, limit int) ([]dtos.HistorySearchHotelResponse, int, error)
	CreateHistorySearchHotel(userID, hotelID uint) (dtos.HistorySearchHotelResponse, error)
	DeleteHistorySearchHotel(userID, id uint) error
}

type historySearchHotelUsecase struct {
	HistorySearchHotelRepository repositories.HistorySearchHotelRepository
	hotelRepository              repositories.HotelRepository
	userRepository               repositories.UserRepository
}

func NewHistorySearchHotelUsecase(HistorySearchHotelRepository repositories.HistorySearchHotelRepository, hotelRepository repositories.HotelRepository, userRepository repositories.UserRepository) HistorySearchHotelUseCase {
	return &historySearchHotelUsecase{HistorySearchHotelRepository, hotelRepository, userRepository}
}

func (u *historySearchHotelUsecase) GetHistorySearchHotelByID(userId, id uint) (models.HistorySearchHotel, error) {
	history, err := u.HistorySearchHotelRepository.GetHistorySearchHotelByID(userId, id)
	if err != nil {
		return models.HistorySearchHotel{}, err
	}

	return history, nil
}

func (u *historySearchHotelUsecase) GetAllHistorySearchHotels(userID uint, page, limit int) ([]dtos.HistorySearchHotelResponse, int, error) {
	var historySearchResponses []dtos.HistorySearchHotelResponse

	histories, count, err := u.HistorySearchHotelRepository.GetHistorySearchHotelsByUserID(userID, page, limit)
	if err != nil {
		return historySearchResponses, count, err
	}

	for _, history := range histories {
		hotel, err := u.hotelRepository.GetHotelByID(history.HotelID)
		if err != nil {
			return historySearchResponses, count, err
		}

		historySearchResponse := dtos.HistorySearchHotelResponse{
			ID:        history.ID,
			UserID:    history.UserID,
			HotelID:   history.HotelID,
			HotelName: hotel.Name,
		}
		historySearchResponses = append(historySearchResponses, historySearchResponse)
	}

	return historySearchResponses, count, nil
}

func (u *historySearchHotelUsecase) CreateHistorySearchHotel(userID, hotelID uint) (dtos.HistorySearchHotelResponse, error) {
	var (
		history         models.HistorySearchHotel
		historyResponse dtos.HistorySearchHotelResponse
	)

	history.UserID = userID
	history.HotelID = hotelID

	history, err := u.HistorySearchHotelRepository.CreateHistorySearchHotel(history)
	if err != nil {
		return historyResponse, err
	}

	hotel, err := u.hotelRepository.GetHotelByID(history.HotelID)
	if err != nil {
		return historyResponse, err
	}

	historyResponse.ID = history.ID
	historyResponse.UserID = history.UserID
	historyResponse.HotelID = history.HotelID
	historyResponse.HotelName = hotel.Name

	return historyResponse, nil
}

func (u *historySearchHotelUsecase) DeleteHistorySearchHotel(userID, id uint) error {
	err := u.HistorySearchHotelRepository.DeleteHistorySearchHotel(userID, id)
	if err != nil {
		return err
	}

	return nil
}
