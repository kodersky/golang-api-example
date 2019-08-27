package main

import (
	"database/sql"
	"fmt"
	"github.com/kodersky/golang-api-example/internal/app/api/order"
	"googlemaps.github.io/maps"
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

	"context"
	"github.com/spf13/viper"
)

func init() {
	cwd, _ := os.Getwd()
	viper.SetConfigFile(fmt.Sprintf("%s/config.yaml", cwd))
	err := viper.ReadInConfig()
	if err != nil {
		viper.SetConfigFile(fmt.Sprint("/opt/config.yaml"))
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`db.host`)
	dbPort := viper.GetInt(`db.port`)
	dbUser := viper.GetString(`db.user`)
	dbPass := viper.GetString(`db.pass`)
	dbName := viper.GetString(`db.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	//val.Add("loc", "Asia/Bangkok")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	dbConn.SetMaxOpenConns(1)
	dbConn.SetMaxIdleConns(1)

	for {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
		err := dbConn.PingContext(ctx)
		if err != nil {
			log.Println(err)
		}
		cancel()
		time.Sleep(time.Second)
	}

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

	gm, err := maps.NewClient(maps.WithAPIKey(viper.GetString(`gmaps_apikey`)))

	if err != nil {
		log.Fatal(err)
	}

	c := newClient(gm)

	ou := _orderUcase.NewOrderUsecase(or, timeoutContext, c)

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	_orderHttpDeliver.NewOrderHandler(e, ou)

	log.Fatal(e.Start(":8081"))
}

func newClient(gm order.GoogleMapClient) *order.Client {
	return &order.Client{
		Client: gm,
	}
}
