package internal

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Product struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"column:name"`
	Description string  `json:"description" gorm:"column:description"`
	Secret      *string `json:"secret,omitempty" gorm:"column:secret"`
	Image       string  `json:"image" gorm:"column:image"`
	Cost        int     `json:"cost" gorm:"column:cost"`
	Purchased   bool    `json:"purchased" gorm:"-"`
}

func (hc *HandlerContext) GetProducts(c *gin.Context) {
	var products []Product
	if err := hc.DB.Omit("Secret").Find(&products).Error; err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	userID := c.GetInt64("tgID")
	var purchases []Purchase
	if err := hc.DB.Where("user_id = ?", userID).Find(&purchases).Error; err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
	}

	for i, product := range products {
		for _, purchase := range purchases {
			if purchase.ProductID == product.ID {
				products[i].Purchased = true
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (hc *HandlerContext) GetProduct(c *gin.Context) {
	productIdParam := c.Param("productId")
	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	var product Product
	if err := hc.DB.Omit("Secret").Where("id = ?", productId).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Нет такого товара"})
			return
		}
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	userID := c.GetInt64("tgID")
	var purchases []Purchase
	if err := hc.DB.Where("user_id = ?", userID).Find(&purchases).Error; err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
	}

	for _, purchase := range purchases {
		if purchase.ProductID == product.ID {
			product.Purchased = true
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (hc *HandlerContext) BuyProduct(c *gin.Context) {
	productIdParam := c.Param("productId")
	productId, err := strconv.Atoi(productIdParam)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	var product Product
	if err := hc.DB.Where("id = ?", productId).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Нет такого товара"})
			return
		}
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	err = hc.DB.Transaction(func(tx *gorm.DB) error {
		// Блокируем таблицу на время выполнения транзакции
		if err := tx.Exec("LOCK TABLE users, purchases IN ACCESS EXCLUSIVE MODE").Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
			return err
		}

		userInterface, _ := c.Get("user")
		user := userInterface.(User)
		if err := tx.Where("id = ?", user.ID).Find(&user).Error; err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
			return err
		}

		if user.Balance < product.Cost {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Не хватает денег на балансе"})
			return nil
		}

		var purchases []Purchase
		if err := tx.Where("username = ?", user.Username).Find(&purchases).Error; err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
			return err
		}
		for _, purchase := range purchases {
			if purchase.ProductID == productId {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Ты уже купил этот товар"})
				return nil
			}
		}

		purchase := Purchase{
			UserID:    user.ID,
			Username:  user.Username,
			ProductID: productId,
		}
		if err := tx.Create(&purchase).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
			return err
		}

		user.Balance -= product.Cost
		if err := tx.Save(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
			return err
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
		return nil
	})

	if err != nil {
		c.Error(err)
		return
	}
}
