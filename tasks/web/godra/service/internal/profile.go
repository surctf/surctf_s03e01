package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hc *HandlerContext) GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"route": "profile"})
}
