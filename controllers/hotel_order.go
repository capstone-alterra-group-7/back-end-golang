package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/middlewares"
	"back-end-golang/models"
	"back-end-golang/usecases"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HotelOrderController interface {
	GetHotelOrders(c echo.Context) error
	GetHotelOrdersByAdmin(c echo.Context) error
	GetHotelOrderDetailByAdmin(c echo.Context) error
	GetHotelOrderByID(c echo.Context) error
	CreateHotelOrder(c echo.Context) error
	CreateHotelOrder2(c echo.Context) error
	UpdateHotelOrder(c echo.Context) error
	CsvHotelOrder(ctx echo.Context) error
}

type hotelOrderController struct {
	hotelOrderUsecase usecases.HotelOrderUsecase
}

func NewHotelOrderController(hotelOrderUsecase usecases.HotelOrderUsecase) HotelOrderController {
	return &hotelOrderController{hotelOrderUsecase}
}

func (c *hotelOrderController) GetHotelOrders(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}

	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 1000
	}

	searchParam := ctx.QueryParam("search")
	nameParam := ctx.QueryParam("name")
	addressParam := ctx.QueryParam("address")
	orderDateParam := ctx.QueryParam("order_date")
	sortParam := ctx.QueryParam("order_by")
	statusParam := ctx.QueryParam("status")

	hotelOrder, count, err := c.hotelOrderUsecase.GetHotelOrders(page, limit, userId, searchParam, nameParam, addressParam, orderDateParam, sortParam, statusParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get a hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully to get order hotels",
			hotelOrder,
			page,
			limit,
			count,
		),
	)
}

func (c *hotelOrderController) GetHotelOrdersByAdmin(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 1000
	}

	searchParam := ctx.QueryParam("search")
	dateStartParam := ctx.QueryParam("date_start")
	dateEndParam := ctx.QueryParam("date_end")
	orderByParam := ctx.QueryParam("order_by")
	filterParam := ctx.QueryParam("filter")
	ratingClass, _ := strconv.Atoi(ctx.QueryParam("rating_class"))

	hotelOrder, count, err := c.hotelOrderUsecase.GetHotelOrdersByAdmin(page, limit, ratingClass, searchParam, dateStartParam, dateEndParam, orderByParam, filterParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get a hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully to get order hotels",
			hotelOrder,
			page,
			limit,
			count,
		),
	)
}

func (c *hotelOrderController) GetHotelOrderDetailByAdmin(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}
	hotelOrderIdParam := ctx.QueryParam("hotel_order_id")
	hotelOrderId, _ := strconv.Atoi(hotelOrderIdParam)

	hotelOrder, err := c.hotelOrderUsecase.GetHotelOrdersDetailByAdmin(uint(hotelOrderId))
	if err != nil {
		return ctx.JSON(
			http.StatusNotFound,
			helpers.NewErrorResponse(
				http.StatusNotFound,
				"Failed to get a hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get order hotels",
			hotelOrder,
		),
	)
}

func (c *hotelOrderController) GetHotelOrderByID(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}

	hotelOrderIdParam := ctx.QueryParam("hotel_order_id")
	hotelOrderId, _ := strconv.Atoi(hotelOrderIdParam)

	isCheckInParam := ctx.QueryParam("update_check_in")
	isCheckIn, _ := strconv.ParseBool(isCheckInParam)

	isCheckOutParam := ctx.QueryParam("update_check_out")
	isCheckOut, _ := strconv.ParseBool(isCheckOutParam)

	hotelOrder, err := c.hotelOrderUsecase.GetHotelOrderByID(userId, uint(hotelOrderId), isCheckIn, isCheckOut)
	if err != nil {
		return ctx.JSON(
			http.StatusNotFound,
			helpers.NewErrorResponse(
				http.StatusNotFound,
				"Failed to get a hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get order hotels",
			hotelOrder,
		),
	)
}

func (c *hotelOrderController) CreateHotelOrder(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}

	var hotelOrderInput dtos.HotelOrderInput
	if err := ctx.Bind(&hotelOrderInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	hotelOrder, err := c.hotelOrderUsecase.CreateHotelOrder(userId, hotelOrderInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a hotel order",
			hotelOrder,
		),
	)
}

func (c *hotelOrderController) CreateHotelOrder2(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}

	var hotelOrderInput dtos.HotelOrderInput
	if err := ctx.Bind(&hotelOrderInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	hotelOrder, err := c.hotelOrderUsecase.CreateHotelOrderMidtrans(userId, hotelOrderInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a hotel order",
			hotelOrder,
		),
	)
}

func (c *hotelOrderController) UpdateHotelOrder(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}

	var hotelOrderInput dtos.HotelOrderInput
	if err := ctx.Bind(&hotelOrderInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	hotelOrderIDParam := ctx.QueryParam("hotel_order_id")
	hotelOrderID, _ := strconv.Atoi(hotelOrderIDParam)

	statusParam := ctx.QueryParam("status")

	hotelOrder, err := c.hotelOrderUsecase.UpdateHotelOrder(userId, uint(hotelOrderID), statusParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to update a hotel order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to update a hotel order",
			hotelOrder,
		),
	)

}
func (c *hotelOrderController) CsvHotelOrder(ctx echo.Context) error {
	var cloudinary models.Url

	// Mendapatkan data hotel order dari sumber data yang sesuai
	hotelOrders, err := c.hotelOrderUsecase.CsvHotelOrder()
	if err != nil {
		fmt.Println("Gagal mendapatkan data hotel order:", err)
		return err
	}

	uploadUrl, err := usecases.NewMediaUpload().RemoteUpload(models.Url{Url: "hotel_orders.csv"})
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
	cloudinary.Url = uploadUrl

	// Buat file CSV
	file, err := os.Create("hotel_orders.csv")
	if err != nil {
		fmt.Println("Gagal membuat file CSV:", err)
		return err
	}
	defer file.Close()

	// Buat writer CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Tulis header kolom
	header := []string{"Hotel Order Code", "Hotel", "Hotel Room", "Check In Date", "Check Out Date", "Number of Night", "Price", "Total Amount", "Name Order", "Email Order", "Phone Number Order", "Status", "Order Date"}
	if err := writer.Write(header); err != nil {
		fmt.Println("Gagal menulis header kolom:", err)
		return err
	}

	// Tulis data
	for _, order := range hotelOrders {
		row := []string{
			order.HotelOrderCode,
			order.Hotel,
			order.HotelRoom,
			order.CheckIn,
			order.CheckOut,
			strconv.Itoa(order.NumberOfNight),
			strconv.Itoa(order.Price),
			strconv.Itoa(order.TotalAmount),
			order.NameOrder,
			order.EmailOrder,
			order.PhoneNumberOrder,
			order.Status,
			order.CreatedAt.String(),
		}

		if err := writer.Write(row); err != nil {
			fmt.Println("Gagal menulis data:", err)
			return err
		}
	}

	fmt.Println("File CSV berhasil dibuat.")

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully uploaded CSV file",
			cloudinary,
		),
	)
}
