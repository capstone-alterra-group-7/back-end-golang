package routes

import (
	"back-end-golang/controllers"
	"back-end-golang/middlewares"
	"back-end-golang/repositories"
	"back-end-golang/usecases"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func init() {
	middleware.ErrJWTMissing.Code = 401
	middleware.ErrJWTMissing.Message = "Unauthorized"
}

func Init(e *echo.Echo, db *gorm.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userRepository := repositories.NewUserRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepository)
	userController := controllers.NewUserController(userUsecase)

	cloudinaryUsecase := usecases.NewMediaUpload()
	cloudinaryController := controllers.NewCloudinaryController(cloudinaryUsecase)

	stationRepository := repositories.NewStationRepository(db)
	stationUsecase := usecases.NewStationUsecase(stationRepository)
	stationController := controllers.NewStationController(stationUsecase)

	trainStationRepository := repositories.NewTrainStationRepository(db)
	// trainStationUsecase := usecases.NewTrainStationUsecase(trainStationRepository)
	// trainStationController := controllers.NewTrainStationController(trainStationUsecase)

	ticketTravelerDetailRepository := repositories.NewTicketTravelerDetailRepository(db)

	trainRepository := repositories.NewTrainRepository(db)
	trainUsecase := usecases.NewTrainUsecase(trainRepository, trainStationRepository)
	trainController := controllers.NewTrainController(trainUsecase)

	trainCarriageRepository := repositories.NewTrainCarriageRepository(db)
	trainCarriageUsecase := usecases.NewTrainCarriageUsecase(trainCarriageRepository, trainRepository, ticketTravelerDetailRepository)
	trainCarriageController := controllers.NewTrainCarriageController(trainCarriageUsecase)

	travelerDetailRepository := repositories.NewTravelerDetailRepository(db)

	trainSeatRepository := repositories.NewTrainSeatRepository(db)

	paymentRepository := repositories.NewPaymentRepository(db)
	paymentUsecase := usecases.NewPaymentUsecase(paymentRepository)
	paymentController := controllers.NewPaymentController(paymentUsecase)


	historySearchRepository := repositories.NewHistorySearchRepository(db)
	historySearchUsecase := usecases.NewHistorySearchUsecase(historySearchRepository, userRepository)
	historySearchController := controllers.NewHistorySearchController(historySearchUsecase)

	ticketOrderRepository := repositories.NewTicketOrderRepository(db)
	ticketOrderUsecase := usecases.NewTicketOrderUsecase(ticketOrderRepository, ticketTravelerDetailRepository, travelerDetailRepository, trainCarriageRepository, trainRepository, trainSeatRepository, stationRepository, trainStationRepository, paymentRepository, userRepository)
	ticketOrderController := controllers.NewTicketOrderController(ticketOrderUsecase)

	hotelRepository := repositories.NewHotelRepository(db)
	hotelRoomRepository := repositories.NewHotelRoomRepository(db)
	hotelRoomImageRepository := repositories.NewHotelRoomImageRepository(db)
	hotelRoomFacilitiesRepository := repositories.NewHotelRoomFacilitiesRepository(db)

	hotelImageRepository := repositories.NewHotelImageRepository(db)
	hotelFacilitiesRepository := repositories.NewHotelFacilitiesRepository(db)
	hotelPolicyRepository := repositories.NewHotelPoliciesRepository(db)
	hotelUsecase := usecases.NewHotelUsecase(hotelRepository, hotelRoomRepository, hotelRoomImageRepository, hotelRoomFacilitiesRepository, hotelImageRepository, hotelFacilitiesRepository, hotelPolicyRepository, historySearchRepository)
	hotelController := controllers.NewHotelController(hotelUsecase)

	hotelRoomUsecase := usecases.NewHotelRoomUsecase(hotelRepository, hotelRoomRepository, hotelRoomImageRepository, hotelRoomFacilitiesRepository)
	hotelRoomController := controllers.NewHotelRoomController(hotelRoomUsecase)

	hotelOrderRepository := repositories.NewHotelOrderRepository(db)
	hotelOrderUsecase := usecases.NewHotelOrderUsecase(hotelOrderRepository, hotelRepository, hotelImageRepository, hotelFacilitiesRepository, hotelPolicyRepository, hotelRoomRepository, hotelRoomImageRepository, hotelRoomFacilitiesRepository, travelerDetailRepository, paymentRepository, userRepository)
	hotelOrderController := controllers.NewHotelOrderController(hotelOrderUsecase)

	dashboardRepository := repositories.NewDashboardRepository(db)
	dashboardUsecase := usecases.NewDashboardUsecase(dashboardRepository, userRepository, ticketOrderRepository, ticketTravelerDetailRepository, travelerDetailRepository, trainCarriageRepository, trainRepository, trainSeatRepository, stationRepository, trainStationRepository, paymentRepository)
	dashboardController := controllers.NewDashboardController(dashboardUsecase)

	articleRepository := repositories.NewArticleRepository(db)
	articleUsecase := usecases.NewArticleUsecase(articleRepository)
	articleController := controllers.NewArticleController(articleUsecase)

	// Middleware CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Izinkan semua domain
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	api := e.Group("/api/v1")
	public := api.Group("/public")

	// cloudinary
	public.POST("/cloudinary/file-upload", cloudinaryController.FileUpload)
	public.POST("/cloudinary/url-upload", cloudinaryController.UrlUpload)

	// USER
	api.POST("/login", userController.UserLogin)
	api.POST("/register", userController.UserRegister)

	user := api.Group("/user")
	user.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("user"))

	// user account
	user.Any("", userController.UserCredential)
	user.PUT("/update-password", userController.UserUpdatePassword)
	user.PUT("/update-profile", userController.UserUpdateProfile)
	user.PUT("/update-photo-profile", userController.UserUpdatePhotoProfile)
	user.DELETE("/delete-photo-profile", userController.UserDeletePhotoProfile)

	// train ka
	public.GET("/train/search", trainController.SearchTrainAvailable)
	user.POST("/train/order", ticketOrderController.CreateTicketOrder)
	user.PATCH("/train/order", ticketOrderController.UpdateTicketOrder)

	public.GET("/hotel/search", hotelController.SearchHotelAvailable)
	user.GET("/order/ticket", ticketOrderController.GetTicketOrders)
	user.GET("/order/ticket/detail", ticketOrderController.GetTicketOrderByID)

	user.POST("/hotel/order", hotelOrderController.CreateHotelOrder)
	user.PATCH("/hotel/order", hotelOrderController.UpdateHotelOrder)

	user.GET("/order/hotel", hotelOrderController.GetHotelOrders)
	user.GET("/order/hotel/detail", hotelOrderController.GetHotelOrderByID)

	user.GET("/history-search", historySearchController.HistorySearchGetAll)
	user.POST("/history-search", historySearchController.HistorySearchCreate)
	user.DELETE("/history-search/:id", historySearchController.HistorySearchDelete)

	// ADMIN

	admin := api.Group("/admin")
	admin.Use(middlewares.JWTMiddleware, middlewares.RoleMiddleware("admin"))

	// users @ admin
	admin.GET("/user", userController.UserGetAll)
	admin.GET("/user/detail", userController.UserGetDetail)
	admin.POST("/user/register", userController.UserAdminRegister)
	admin.PUT("/user/update/:id", userController.UserAdminUpdate)

	admin.GET("/dashboard", dashboardController.DashboardGetAll)

	admin.GET("/order/ticket", ticketOrderController.GetTicketOrdersByAdmin)
	admin.GET("/order/ticket/detail", ticketOrderController.GetTicketOrderDetailByAdmin)

	admin.GET("/order/hotel", hotelOrderController.GetHotelOrdersByAdmin)
	admin.GET("/order/hotel/detail", hotelOrderController.GetHotelOrderDetailByAdmin)

	// crud station
	public.GET("/station", stationController.GetAllStations)
	public.GET("/station/:id", stationController.GetStationByID)
	admin.GET("/station", stationController.GetAllStationsByAdmin)
	admin.PUT("/station/:id", stationController.UpdateStation)
	admin.POST("/station", stationController.CreateStation)
	admin.DELETE("/station/:id", stationController.DeleteStation)

	// crud train
	public.GET("/train", trainController.GetAllTrains)
	public.GET("/train/:id", trainController.GetTrainByID)
	admin.GET("/train", trainController.GetAllTrainsByAdmin)
	admin.PUT("/train/:id", trainController.UpdateTrain)
	admin.POST("/train", trainController.CreateTrain)
	admin.DELETE("/train/:id", trainController.DeleteTrain)

	public.GET("/train-carriage", trainCarriageController.GetAllTrainCarriages)
	public.GET("/train-carriage/:id", trainCarriageController.GetTrainCarriageByID)
	admin.PUT("/train-carriage/:id", trainCarriageController.UpdateTrainCarriage)
	admin.POST("/train-carriage", trainCarriageController.CreateTrainCarriage)
	admin.DELETE("/train-carriage/:id", trainCarriageController.DeleteTrainCarriage)

	public.GET("/article", articleController.GetAllArticles)
	public.GET("/article/:id", articleController.GetArticleByID)
	admin.PUT("/article/:id", articleController.UpdateArticle)
	admin.POST("/article", articleController.CreateArticle)
	admin.DELETE("/article/:id", articleController.DeleteArticle)

	public.GET("/payment", paymentController.GetAllPayments)
	public.GET("/payment/:id", paymentController.GetPaymentByID)
	admin.PUT("/payment/:id", paymentController.UpdatePayment)
	admin.POST("/payment", paymentController.CreatePayment)
	admin.DELETE("/payment/:id", paymentController.DeletePayment)

	public.GET("/hotel", hotelController.GetAllHotels)
	public.GET("/hotel/:id", hotelController.GetHotelByID)
	admin.PUT("/hotel/:id", hotelController.UpdateHotel)
	admin.POST("/hotel", hotelController.CreateHotel)
	admin.DELETE("/hotel/:id", hotelController.DeleteHotel)

	public.GET("/hotel-room", hotelRoomController.GetAllHotelRooms)
	public.GET("/hotel-room/:id", hotelRoomController.GetHotelRoomByID)
	admin.PUT("/hotel-room/:id", hotelRoomController.UpdateHotelRoom)
	admin.POST("/hotel-room", hotelRoomController.CreateHotelRoom)
	admin.DELETE("/hotel-room/:id", hotelRoomController.DeleteHotelRoom)
}
