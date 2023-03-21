package internal

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Username string `json:"username" gorm:"column:username;unique;not null"`
	Balance  int    `json:"balance" gorm:"column:balance;"`
}

type Purchase struct {
	UserID    int64   `json:"-" gorm:"primaryKey;column:user_id;not null"`
	Username  string  `json:"-" gorm:"primaryKey;column:username;not null"`
	ProductID int     `json:"-" gorm:"primaryKey;column:product_id;not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}

func (hc *HandlerContext) CreateUser(c *gin.Context) {
	userID := c.GetInt64("tgID")
	body := c.GetString("requestBody")

	var user User
	if err := json.Unmarshal([]byte(body), &user); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	if len(user.Username) > 24 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Длина юзернейма должна быть не больше 24 символов"})
		return
	}

	user.ID = userID
	if err := hc.DB.Create(user).Error; err != nil {
		var pqErr *pgconn.PgError
		if ok := errors.As(err, &pqErr); ok && pqErr.Code == "23505" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Юзернейм уже занят"})
			return
		}
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (hc *HandlerContext) GetUser(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(User)

	var purchases []Purchase
	if err := hc.DB.Where("user_id = ?", user.ID).Preload("Product").Find(&purchases).Error; err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
	}
	c.JSON(http.StatusOK, gin.H{"user": user, "purchases": purchases})
}

func (hc *HandlerContext) DeleteUser(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(User)

	if err := hc.DB.Where("id = ?", user.ID).Delete(&user).Error; err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
