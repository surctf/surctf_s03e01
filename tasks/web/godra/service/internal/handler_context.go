package internal

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type HandlerContext struct {
	DB *gorm.DB
}

func (hc *HandlerContext) PostSignUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"good": "good"})
}
