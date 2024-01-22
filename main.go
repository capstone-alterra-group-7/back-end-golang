package main

import (
	"back-end-golang/configs"
	_ "back-end-golang/docs"
	"back-end-golang/routes"
	"log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           Tripease API Documentation
// @version         1.0
// @termsOfService  http://swagger.io/terms/

// @contact.name   Capstone Alterra Group 7
// @contact.url    https://github.com/capstone-alterra-group-7

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// api.tripease.my.id
// localhost:8088

// @host      api.tripease.my.id
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := configs.ConnectDB()
	if err != nil {
		panic(err)
	}

	err = configs.MigrateDB(db)
	if err != nil {
		panic(err)
	}

	err = configs.TrainSeatSeeder(db)
	if err != nil {
		panic(err)
	}

	err = configs.AccountSeeder(db)
	if err != nil {
		panic(err)
	}

	routes.Init(e, db)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// e.Logger.Fatal(e.Start(":443"))

	// Start server with TLS
	err = e.StartTLS(":443", "/etc/nginx/ssl/fullchain.pem", "/etc/nginx/ssl/privkey.pem")
	if err != nil {
		log.Fatal("StartTLS: ", err)
	}
}
