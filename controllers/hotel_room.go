package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type HotelRoomController interface {
	//admin
	GetAllHotelRooms(c echo.Context) error
	GetHotelRoomByID(c echo.Context) error
	CreateHotelRoom(c echo.Context) error
	UpdateHotelRoom(c echo.Context) error
	DeleteHotelRoom(c echo.Context) error
}

type hotelRoomController struct {
	hotelRoomUsecase usecases.HotelRoomUsecase
}

func NewHotelRoomController(hotelRoomUsecase usecases.HotelRoomUsecase) HotelRoomController {
	return &hotelRoomController{hotelRoomUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *hotelRoomController) GetAllHotelRooms(ctx echo.Context) error {
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

	hotelRooms, count, err := c.hotelRoomUsecase.GetAllHotelRooms(page, limit)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all hotel room",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all hotel room",
			hotelRooms,
			page,
			limit,
			count,
		),
	)
}

func (c *hotelRoomController) GetHotelRoomByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	hotelRoom, err := c.hotelRoomUsecase.GetHotelRoomByID(uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get hotel room by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get hotel room by id",
			hotelRoom,
		),
	)

}

func (c *hotelRoomController) CreateHotelRoom(ctx echo.Context) error {
	var hotelRoomDTO dtos.HotelRoomInput
	if err := ctx.Bind(&hotelRoomDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding hotel",
				helpers.GetErrorData(err),
			),
		)
	}

	hotelRoom, err := c.hotelRoomUsecase.CreateHotelRoom(&hotelRoomDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a hotel room",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a hotel room",
			hotelRoom,
		),
	)
}

func (c *hotelRoomController) UpdateHotelRoom(ctx echo.Context) error {

	var hotelRoomInput dtos.HotelRoomInput
	if err := ctx.Bind(&hotelRoomInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding hotel",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	hotelRoom, err := c.hotelRoomUsecase.GetHotelRoomByID(uint(id))
	if hotelRoom.HotelRoomID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get hotel room by id",
				helpers.GetErrorData(err),
			),
		)
	}

	hotelRoomResp, err := c.hotelRoomUsecase.UpdateHotelRoom(uint(id), hotelRoomInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to updated a hotel room",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated hotel room",
			hotelRoomResp,
		),
	)
}

func (c *hotelRoomController) DeleteHotelRoom(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.hotelRoomUsecase.DeleteHotelRoom(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete hotel room",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted hotel room",
			nil,
		),
	)
}
