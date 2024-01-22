package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/middlewares"
	"back-end-golang/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HotelRatingsController interface {
	//user
	CreateHotelRating(c echo.Context) error
	GetHotelRatingsByIdOrders(c echo.Context) error
	GetAllHotelRatingsByIdHotels(c echo.Context) error
	//admin
	GetRatingsByHotelsId(c echo.Context) error
}

type hotelRatingsController struct {
	hotelRatingUsecase usecases.HotelRatingsUsecase
}

func NewHotelRatingsController(hotelRatingUsecase usecases.HotelRatingsUsecase) HotelRatingsController {
	return &hotelRatingsController{hotelRatingUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *hotelRatingsController) CreateHotelRating(ctx echo.Context) error {
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
	var hotelRatingsInputDTO dtos.HotelRatingInput

	if err := ctx.Bind(&hotelRatingsInputDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding hotel rating",
				helpers.GetErrorData(err),
			),
		)
	}

	ratings, err := c.hotelRatingUsecase.CreateHotelRating(userId, hotelRatingsInputDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to create hotel rating",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Success to create hotel rating",
			ratings,
		),
	)
}

func (c *hotelRatingsController) GetRatingsByHotelsId(ctx echo.Context) error {
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

	ratingParam := ctx.QueryParam("rating")
	rating, err := strconv.Atoi(ratingParam)
	if err != nil {
		rating = 10
	}

	filter := ctx.QueryParam("filter")
	if filter == "" {
		filter = "all"
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	hotelId := uint(id)
	ratings, count, err := c.hotelRatingUsecase.GetHotelRatingsByHotelID(rating, page, limit, hotelId, filter)
	if err != nil {
		return ctx.JSON(
			http.StatusNotFound,
			helpers.NewErrorResponse(
				http.StatusNotFound,
				"Failed to get hotel rating",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Success fetching hotel rating",
			ratings,
			page,
			limit,
			count,
		),
	)
}

func (c *hotelRatingsController) GetHotelRatingsByIdOrders(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	hotelOrderId := uint(id)
	ratings, err := c.hotelRatingUsecase.GetHotelRatingsByIdOrders(hotelOrderId)
	if err != nil {
		return ctx.JSON(
			http.StatusNotFound,
			helpers.NewErrorResponse(
				http.StatusNotFound,
				"Failed to get hotel rating",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success to get hotel rating",
			ratings,
		),
	)
}

func (c *hotelRatingsController) GetAllHotelRatingsByIdHotels(ctx echo.Context) error {
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

	id, _ := strconv.Atoi(ctx.Param("id"))

	hotelId := uint(id)

	ratings, count, err := c.hotelRatingUsecase.GetAllHotelRatingsByIdHotels(page, limit, hotelId)
	if err != nil {
		return ctx.JSON(
			http.StatusNotFound,
			helpers.NewErrorResponse(
				http.StatusNotFound,
				"Failed fetching hotel rating",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Success fetching hotel rating",
			ratings,
			page,
			limit,
			count,
		),
	)

}
