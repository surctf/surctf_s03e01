package internal

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Secret      string `json:"secret" gorm:"column:secret"`
	Image       string `json:"image" gorm:"column:image"`
	Cost        int    `json:"cost" gorm:"column:cost"`
	Purchased   bool   `json:"purchased" gorm:"-"`
}

func (hc *HandlerContext) GetProducts(c *gin.Context) {
	var products []Product
	if err := hc.DB.Omit("Secret").Find(&products).Error; err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		return
	}

	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "products.html", gin.H{"User": user, "Products": products})
}

func (hc *HandlerContext) GetProduct(c *gin.Context) {
	productIdParam := c.Param("productId")
	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		return
	}

	var product Product
	if err := hc.DB.Omit("Secret").Where("id = ?", productId).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HTML(http.StatusNotFound, "not_found.html", nil)
			return
		}
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "product_card.html", gin.H{"Product": product})
}

func (hc *HandlerContext) BuyProduct(c *gin.Context) {
	productIdParam := c.Param("productId")
	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		return
	}

	var product Product
	if err := hc.DB.Where("id = ?", productId).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HTML(http.StatusNotFound, "not_found.html", nil)
			return
		}
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		return
	}

	err = hc.DB.Transaction(func(tx *gorm.DB) error {
		// Блокируем таблицу на время выполнения транзакции
		if err := tx.Exec("LOCK TABLE users, purchases IN ACCESS EXCLUSIVE MODE").Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		userInterface, _ := c.Get("user")
		user := userInterface.(User)
		if err := tx.Where("id = ?", user.ID).Find(&user).Error; err != nil {
			c.Error(err)
			c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		}

		if user.Balance < product.Cost {
			c.HTML(http.StatusConflict, "error.html", gin.H{"Error": "у тебя нет столько денег"})
			return nil
		}
		log.Println("TRANS")

		var purchases []Purchase
		if err := tx.Where("username = ?", user.Username).Find(&purchases).Error; err != nil {
			c.Error(err)
			c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		}
		for _, purchase := range purchases {
			if purchase.ProductID == productId {
				c.HTML(http.StatusConflict, "error.html", gin.H{"Error": "ты уже купил этот товар"})
				return errors.New("уже куплен товар")
			}
		}

		purchase := Purchase{
			UserID:    user.ID,
			Username:  user.Username,
			ProductID: productId,
		}
		if err := tx.Create(&purchase).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
			return err
		}

		user.Balance -= product.Cost
		if err := tx.Save(&user).Error; err != nil {
			c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
			return err
		}

		return nil
	})

	if err != nil {
		c.Error(err)
		return
	}
}
