package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/models"
	"back-end-golang/repositories"
	"errors"
)

type PaymentUsecase interface {
	GetAllPayments(page, limit int) ([]dtos.PaymentResponse, int, error)
	GetPaymentByID(id uint) (dtos.PaymentResponse, error)
	CreatePayment(payment *dtos.PaymentInput) (dtos.PaymentResponse, error)
	UpdatePayment(id uint, payment dtos.PaymentInput) (dtos.PaymentResponse, error)
	DeletePayment(id uint) error
}

type paymentUsecase struct {
	paymentRepo repositories.PaymentRepository
}

func NewPaymentUsecase(PaymentRepo repositories.PaymentRepository) PaymentUsecase {
	return &paymentUsecase{PaymentRepo}
}

// GetAllPayments godoc
// @Summary      Get all payments
// @Description  Get all payments
// @Tags         Admin - Payment
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllPaymentStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/payment [get]
func (u *paymentUsecase) GetAllPayments(page, limit int) ([]dtos.PaymentResponse, int, error) {
	payments, count, err := u.paymentRepo.GetAllPayments(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var paymentResponses []dtos.PaymentResponse
	for _, payment := range payments {
		paymentResponse := dtos.PaymentResponse{
			ID:            int(payment.ID),
			Type:          payment.Type,
			ImageUrl:      payment.ImageUrl,
			Name:          payment.Name,
			AccountName:   payment.AccountName,
			AccountNumber: payment.AccountNumber,
			CreatedAt:     &payment.CreatedAt,
			UpdatedAt:     &payment.UpdatedAt,
		}
		paymentResponses = append(paymentResponses, paymentResponse)
	}

	return paymentResponses, count, nil
}

// GetPaymentByID godoc
// @Summary      Get payment by ID
// @Description  Get payment by ID
// @Tags         Admin - Payment
// @Accept       json
// @Produce      json
// @Param id path integer true "ID payment"
// @Success      200 {object} dtos.PaymentStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/payment/{id} [get]
func (u *paymentUsecase) GetPaymentByID(id uint) (dtos.PaymentResponse, error) {
	var paymentResponses dtos.PaymentResponse
	payment, err := u.paymentRepo.GetPaymentByID2(id)
	if err != nil {
		return paymentResponses, err
	}
	paymentResponse := dtos.PaymentResponse{
		ID:            int(payment.ID),
		Type:          payment.Type,
		ImageUrl:      payment.ImageUrl,
		Name:          payment.Name,
		AccountName:   payment.AccountName,
		AccountNumber: payment.AccountNumber,
		CreatedAt:     &payment.CreatedAt,
		UpdatedAt:     &payment.UpdatedAt,
	}
	return paymentResponse, nil
}

// CreatePayment godoc
// @Summary      Create a new payment
// @Description  Create a new payment
// @Tags         Admin - Payment
// @Accept       json
// @Produce      json
// @Param        file formData file true "Image file"
// @Param		 type formData string true "Type payment"
// @Param		 name formData string true "Name payment"
// @Param		 account_name formData string true "Account name payment"
// @Param		 account_number formData string true "Account number payment"
// @Success      201 {object} dtos.PaymentCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/payment [post]
// @Security BearerAuth
func (u *paymentUsecase) CreatePayment(payment *dtos.PaymentInput) (dtos.PaymentResponse, error) {
	var paymentResponses dtos.PaymentResponse
	if payment.Type == "" || payment.Name == "" || payment.AccountName == "" || payment.AccountNumber == "" || payment.ImageUrl == "" {
		return paymentResponses, errors.New("Failed to create payment")
	}
	createPayment := models.Payment{
		Type:          payment.Type,
		ImageUrl:      payment.ImageUrl,
		Name:          payment.Name,
		AccountName:   payment.AccountName,
		AccountNumber: payment.AccountNumber,
	}

	createdPayment, err := u.paymentRepo.CreatePayment(createPayment)
	if err != nil {
		return paymentResponses, err
	}

	paymentResponse := dtos.PaymentResponse{
		ID:            int(createdPayment.ID),
		Type:          createdPayment.Type,
		ImageUrl:      createdPayment.ImageUrl,
		Name:          createdPayment.Name,
		AccountName:   createdPayment.AccountName,
		AccountNumber: createdPayment.AccountNumber,
		CreatedAt:     &createdPayment.CreatedAt,
		UpdatedAt:     &createdPayment.UpdatedAt,
	}
	return paymentResponse, nil
}

// UpdatePayment godoc
// @Summary      Update payment
// @Description  Update payment
// @Tags         Admin - Payment
// @Accept       json
// @Produce      json
// @Param		 id path integer true "ID payment"
// @Param        file formData file true "Image file"
// @Param		 type formData string true "Type payment"
// @Param		 name formData string true "Name payment"
// @Param		 account_name formData string true "Account name payment"
// @Param		 account_number formData string true "Account number payment"
// @Success      200 {object} dtos.PaymentStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/payment/{id} [put]
// @Security BearerAuth
func (u *paymentUsecase) UpdatePayment(id uint, paymentInput dtos.PaymentInput) (dtos.PaymentResponse, error) {
	var payment models.Payment
	var paymentResponse dtos.PaymentResponse
	if payment.Type == "" || payment.Name == "" || payment.AccountName == "" || payment.AccountNumber == "" || payment.ImageUrl == "" {
		return paymentResponse, errors.New("Failed to update payment")
	}

	payment, err := u.paymentRepo.GetPaymentByID2(id)
	if err != nil {
		return paymentResponse, err
	}

	payment.Type = paymentInput.Type
	payment.Name = paymentInput.Name
	payment.AccountName = paymentInput.AccountName
	payment.AccountNumber = paymentInput.AccountNumber

	payment, err = u.paymentRepo.UpdatePayment(payment)

	if err != nil {
		return paymentResponse, err
	}

	paymentResponse.ID = int(payment.ID)
	paymentResponse.Type = payment.Type
	paymentResponse.Name = payment.Name
	paymentResponse.AccountName = payment.AccountName
	paymentResponse.AccountNumber = payment.AccountNumber
	paymentResponse.CreatedAt = &payment.CreatedAt
	paymentResponse.UpdatedAt = &payment.UpdatedAt

	return paymentResponse, nil

}

// DeletePayment godoc
// @Summary      Delete a payment
// @Description  Delete a payment
// @Tags         Admin - Payment
// @Accept       json
// @Produce      json
// @Param id path integer true "ID payment"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/payment/{id} [delete]
// @Security BearerAuth
func (u *paymentUsecase) DeletePayment(id uint) error {
	return u.paymentRepo.DeletePayment(id)
}
