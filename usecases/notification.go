package usecases

import (
	"back-end-golang/models"
	"back-end-golang/repositories"
	"strconv"
	"strings"
)

type NotificationUsecase interface {
	CreateNotification(userID uint, templateMessageID uint, ticketTravelerDetailRepository *models.TicketTravelerDetail) error
}

type notificationUsecase struct {
	notificationRepo               repositories.NotificationRepository
	templateMessageRepo            repositories.TemplateMessageRepository
	userRepo                       repositories.UserRepository
	ticketTravelerDetailRepository repositories.TicketTravelerDetailRepository
}

func NewNotificationUsecase(notificationRepo repositories.NotificationRepository, templateMessageRepo repositories.TemplateMessageRepository, userRepo repositories.UserRepository, ticketTravelerDetailRepository repositories.TicketTravelerDetailRepository) NotificationUsecase {
	return &notificationUsecase{
		notificationRepo:               notificationRepo,
		templateMessageRepo:            templateMessageRepo,
		userRepo:                       userRepo,
		ticketTravelerDetailRepository: ticketTravelerDetailRepository,
	}
}

func (u *notificationUsecase) CreateNotification(userID uint, templateMessageID uint, ticketTravelerDetailRepository *models.TicketTravelerDetail) error {
	// Dapatkan template message dari database berdasarkan ID
	templateMessage, err := u.templateMessageRepo.GetTemplateMessageByID(templateMessageID)
	if err != nil {
		return err
	}

	// Mendapatkan user dari database berdasarkan ID
	user, err := u.userRepo.UserGetById(userID)
	if err != nil {
		return err
	}

	trainCarriageID := strconv.FormatUint(uint64(ticketTravelerDetailRepository.TrainCarriageID), 10)
	dateOfDeparture := ticketTravelerDetailRepository.DateOfDeparture.Format("2006-01-02")

	content := strings.Replace(templateMessage.Content, "Nama :", user.FullName, 1)
	content = strings.Replace(content, "Nomor Tiket :", ticketTravelerDetailRepository.BoardingTicketCode, 1)
	content = strings.Replace(content, "Nomor Kereta :", trainCarriageID, 1)
	content = strings.Replace(content, "Jam Keberangkatan :", dateOfDeparture, 1)
	content = strings.Replace(content, "Jam Kedatangan :", ticketTravelerDetailRepository.DepartureTime, 1)

	// Buat notifikasi dan simpan ke dalam database
	err = u.notificationRepo.CreateNotificationRepository(userID, templateMessageID, content)
	if err != nil {
		return err
	}

	return nil
}
