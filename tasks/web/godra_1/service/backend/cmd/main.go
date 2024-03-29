package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"godra/backend/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
	"time"
)

func getSecret(token []byte) []byte {
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write(token)
	return h.Sum(nil)
}

func main() {
	_ = os.Mkdir("./logs", os.ModePerm)
	f, err := os.OpenFile("./logs/log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
		panic(err)
	}
	defer f.Close()

	log := &logrus.Logger{
		// Log into f file handler and on os.Stdout
		Out:   io.MultiWriter(f, os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	if err := configureDB(db); err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	router.Use(internal.JSONLogMiddleware(log),
		internal.TgAuthMiddleware(getSecret([]byte(os.Getenv("BOT_TOKEN")))),
	)
	registeredMW := internal.RegisteredMiddleware(db, router)

	handlerContext := internal.HandlerContext{DB: db}

	router.GET("/api/colab", registeredMW, handlerContext.GetColabFlag)

	router.GET("/api/user", registeredMW, handlerContext.GetUser)
	router.POST("/api/user", internal.ReadRequestBodyMiddleware, handlerContext.CreateUser)
	router.DELETE("/api/user", registeredMW, handlerContext.DeleteUser)

	router.GET("/api/products", registeredMW, handlerContext.GetProducts)
	router.GET("/api/products/:productId", registeredMW, handlerContext.GetProduct)
	router.POST("/api/products/:productId/buy", registeredMW, handlerContext.BuyProduct)

	if err := router.Run(os.Getenv("SERVICE_ADDR")); err != nil {
		log.Fatalln(err)
	}
}

func configureDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	maxOpenConns, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONS"))
	if err != nil {
		return fmt.Errorf("DB_MAX_OPEN_CONS: %w", err)
	}
	maxIdleConns, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONS"))
	if err != nil {
		return fmt.Errorf("DB_MAX_IDLE_CONS: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	if err := db.AutoMigrate(&internal.User{}, &internal.Product{}, &internal.Purchase{}); err != nil {
		return err
	}

	// Создаем записи продуктов, если их нет
	productsJSON := []byte(os.Getenv("PRODUCTS_JSON"))
	p := struct {
		Products []internal.Product `json:"products"`
	}{}
	if err := json.Unmarshal(productsJSON, &p); err != nil {
		return err
	}

	for _, product := range p.Products {
		if err := db.FirstOrCreate(&product).Error; err != nil {
			return err
		}
	}

	return nil
}
