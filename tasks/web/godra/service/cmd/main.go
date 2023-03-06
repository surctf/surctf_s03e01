package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"godra/internal"
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

	handlerContext := internal.HandlerContext{DB: db}

	router := gin.Default()
	router.Use(internal.JSONLogMiddleware(log),
		internal.TgAuthMiddleware(getSecret([]byte(os.Getenv("BOT_TOKEN")))),
	)

	//router.GET("/signup", handlerContext.GetSignUp)
	router.POST("/signup", handlerContext.PostSignUp)
	//router.GET("/profile", handlerContext.GetPrivacyPolicy)
	//router.DELETE("/profile", handlerContext.GetPrivacyPolicy)
	//
	//router.GET("/products", handlerContext.GetPrivacyPolicy)
	//router.GET("/products/:productId", handlerContext.GetFAQ)
	//router.POST("/products/:productId/buy", handlerContext.GetFAQ)

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
	//if err := db.FirstOrCreate(&models.Event{}).Error; err != nil {
	//	return err
	//}

	return nil
}
