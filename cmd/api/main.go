package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kodersky/golang-api-example/internal/app/api/middleware"
	_orderHttpDeliver "github.com/kodersky/golang-api-example/internal/app/api/order/delivery/http"
	_orderRepo "github.com/kodersky/golang-api-example/internal/app/api/order/repository"
	_orderUcase "github.com/kodersky/golang-api-example/internal/app/api/order/usecase"
	"github.com/labstack/echo"

	"github.com/spf13/viper"
)

func init() {
	cwd, _ := os.Getwd()
	viper.SetConfigFile(fmt.Sprintf("%s/.env", cwd))
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`db_host`)
	dbPort := viper.GetInt(`db_port`)
	dbUser := viper.GetString(`db_user`)
	dbPass := viper.GetString(`db_pass`)
	dbName := viper.GetString(`db_name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Bangkok")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		log.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	or := _orderRepo.NewMysqlOrderRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("timeout")) * time.Second

	ou := _orderUcase.NewOrderUsecase(or, timeoutContext)

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	_orderHttpDeliver.NewOrderHandler(e, ou)

	log.Fatal(e.Start(":8081"))
}
