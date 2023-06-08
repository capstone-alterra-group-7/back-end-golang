package repositories

import (
	"back-end-golang/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	DashboardGetAll() (int, int, int, int, int, int, []models.TicketTravelerDetail, []models.User, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) DashboardGetAll() (int, int, int, int, int, int, []models.TicketTravelerDetail, []models.User, error) {
	var countUser int64
	var countUserToday int64

	user := models.User{}
	err := r.db.Unscoped().Where("role = 'user'").Order("id DESC").Find(&user).Count(&countUser).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	newUser := []models.User{}
	err = r.db.Unscoped().Where("role = 'user'").Order("id DESC").Find(&newUser).Limit(10).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	// Get the start and end of the current day
	startOfDay := time.Now().Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Count the number of users created today
	err = r.db.Model(&models.User{}).Unscoped().Where("role = 'user' AND created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countUserToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	var countTrain int64
	var countTrainToday int64
	train := models.Train{}
	err = r.db.Unscoped().Find(&train).Count(&countTrain).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	// Count the number of users created today
	err = r.db.Model(&models.Train{}).Unscoped().Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countTrainToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	ticketOrder := []models.TicketTravelerDetail{}

	var countTicketOrder int64
	var countTicketOrderToday int64
	err = r.db.Find(&ticketOrder).Count(&countTicketOrder).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	newTicketOrder := []models.TicketTravelerDetail{}
	err = r.db.Find(&newTicketOrder).Order("id DESC").Limit(10).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	// Count the number of users created today
	err = r.db.Model(&models.TicketOrder{}).Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).Count(&countTicketOrderToday).Error
	if err != nil {
		return 0, 0, 0, 0, 0, 0, []models.TicketTravelerDetail{}, []models.User{}, err
	}

	return int(countUser), int(countUserToday), int(countTrain), int(countTrainToday), int(countTicketOrder), int(countTicketOrderToday), newTicketOrder, newUser, nil
}
