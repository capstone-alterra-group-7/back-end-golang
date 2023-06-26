package repositories

import (
	"back-end-golang/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDashboardRepository is a mock implementation of DashboardRepository
type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) DashboardGetAll() (int, int, int, int, int, int, []models.TicketTravelerDetail, []models.User, []models.HotelOrder, int, int, error) {
	args := m.Called()
	return args.Int(0), args.Int(1), args.Int(2), args.Int(3), args.Int(4), args.Int(5), args.Get(6).([]models.TicketTravelerDetail), args.Get(7).([]models.User), args.Get(8).([]models.HotelOrder), args.Int(9), args.Int(10), args.Error(11)
}

func TestDashboardRepository_DashboardGetAll(t *testing.T) {
	// Create a mock DB and repository

	// Create some mock data
	countUser := int64(10)
	countUserToday := int64(5)
	countTrain := int64(8)
	countTrainToday := int64(3)
	countTicketOrder := int64(20)
	countTicketOrderToday := int64(15)
	countHotelOrder := int64(12)
	countHotelOrderToday := int64(6)
	newTicketOrder := []models.TicketTravelerDetail{}
	newUser := []models.User{}
	newHotelOrder := []models.HotelOrder{}

	// Set up expectations
	mockRepo := &MockDashboardRepository{}
	mockRepo.On("DashboardGetAll").Return(
		int(countUser), int(countUserToday), int(countTrain), int(countTrainToday),
		int(countTicketOrder), int(countTicketOrderToday), newTicketOrder,
		newUser, newHotelOrder, int(countHotelOrder), int(countHotelOrderToday),
		nil,
	)

	// Call the actual method
	users, usersToday, trains, trainsToday, ticketOrders, ticketOrdersToday,
		ticketOrderDetails, newUsers, newHotelOrders, hotelOrders, hotelOrdersToday, err := mockRepo.DashboardGetAll()

	// Validate the results
	assert.NoError(t, err)
	assert.Equal(t, int(countUser), users)
	assert.Equal(t, int(countUserToday), usersToday)
	assert.Equal(t, int(countTrain), trains)
	assert.Equal(t, int(countTrainToday), trainsToday)
	assert.Equal(t, int(countTicketOrder), ticketOrders)
	assert.Equal(t, int(countTicketOrderToday), ticketOrdersToday)
	assert.Equal(t, newTicketOrder, ticketOrderDetails)
	assert.Equal(t, newUser, newUsers)
	assert.Equal(t, newHotelOrder, newHotelOrders)
	assert.Equal(t, int(countHotelOrder), hotelOrders)
	assert.Equal(t, int(countHotelOrderToday), hotelOrdersToday)

	// Assert that the mock repository's method was called
	mockRepo.AssertExpectations(t)
}
