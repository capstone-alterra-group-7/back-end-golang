package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/models"
	"back-end-golang/usecases"
	"net/http"
	"regexp"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PaymentController interface {
	GetAllPayments(c echo.Context) error
	GetPaymentByID(c echo.Context) error
	CreatePayment(c echo.Context) error
	UpdatePayment(c echo.Context) error
	DeletePayment(c echo.Context) error
}

type paymentController struct {
	paymentUsecase usecases.PaymentUsecase
}

func NewPaymentController(paymentUsecase usecases.PaymentUsecase) PaymentController {
	return &paymentController{paymentUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *paymentController) GetAllPayments(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	payments, count, err := c.paymentUsecase.GetAllPayments(page, limit)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching payment",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all payments",
			payments,
			page,
			limit,
			count,
		),
	)
}

func (c *paymentController) GetPaymentByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	payment, err := c.paymentUsecase.GetPaymentByID(uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get payment by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get payment by id",
			payment,
		),
	)

}

func (c *paymentController) CreatePayment(ctx echo.Context) error {
	var paymentInput dtos.PaymentInput
	if err := ctx.Bind(&paymentInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding payment",
				helpers.GetErrorData(err),
			),
		)
	}

	if paymentInput.ImageUrl == "" {
		formHeader, err := ctx.FormFile("file")
		if err != nil {
			if err != nil {
				return ctx.JSON(
					http.StatusInternalServerError,
					helpers.NewErrorResponse(
						http.StatusInternalServerError,
						"Error uploading photo",
						helpers.GetErrorData(err),
					),
				)
			}
		}

		//get file from header
		formFile, err := formHeader.Open()
		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)

		if !re.MatchString(formHeader.Filename) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(err),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().FileUpload(models.File{File: formFile})

		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}
		paymentInput.ImageUrl = uploadUrl
	} else {
		var url models.Url
		url.Url = paymentInput.ImageUrl

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)
		if !re.MatchString(paymentInput.ImageUrl) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(nil),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().RemoteUpload(url)
		if uploadUrl == "" || err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		paymentInput.ImageUrl = uploadUrl
	}

	payment, err := c.paymentUsecase.CreatePayment(&paymentInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a payment",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a payment",
			payment,
		),
	)
}

func (c *paymentController) UpdatePayment(ctx echo.Context) error {
	var paymentInput dtos.PaymentInput
	if err := ctx.Bind(&paymentInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed fetching payment",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	payment, err := c.paymentUsecase.GetPaymentByID(uint(id))
	if payment.ID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get payment by id",
				helpers.GetErrorData(err),
			),
		)
	}

	if paymentInput.ImageUrl == "" {
		formHeader, err := ctx.FormFile("file")
		if err != nil {
			if err != nil {
				return ctx.JSON(
					http.StatusInternalServerError,
					helpers.NewErrorResponse(
						http.StatusInternalServerError,
						"Error uploading photo",
						helpers.GetErrorData(err),
					),
				)
			}
		}

		//get file from header
		formFile, err := formHeader.Open()
		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)

		if !re.MatchString(formHeader.Filename) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(err),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().FileUpload(models.File{File: formFile})

		if err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}
		paymentInput.ImageUrl = uploadUrl
	} else {
		var url models.Url
		url.Url = paymentInput.ImageUrl

		var re = regexp.MustCompile(`.png|.jpeg|.jpg`)
		if !re.MatchString(paymentInput.ImageUrl) {
			return ctx.JSON(
				http.StatusBadRequest,
				helpers.NewErrorResponse(
					http.StatusBadRequest,
					"The provided file format is not allowed. Please upload a JPEG or PNG image",
					helpers.GetErrorData(nil),
				),
			)
		}

		uploadUrl, err := usecases.NewMediaUpload().RemoteUpload(url)
		if uploadUrl == "" || err != nil {
			return ctx.JSON(
				http.StatusInternalServerError,
				helpers.NewErrorResponse(
					http.StatusInternalServerError,
					"Error uploading photo",
					helpers.GetErrorData(err),
				),
			)
		}

		paymentInput.ImageUrl = uploadUrl
	}

	paymentResp, err := c.paymentUsecase.UpdatePayment(uint(id), paymentInput)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed update payment",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated payment",
			paymentResp,
		),
	)
}

func (c *paymentController) DeletePayment(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.paymentUsecase.DeletePayment(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete payment",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted payment",
			nil,
		),
	)
}
