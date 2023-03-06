package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hc *HandlerContext) PostPurchase(c *gin.Context) {
	user, _ := c.Get("user")
	// написать создание покупки
	c.HTML(http.StatusOK, "products.html", gin.H{"User": user})
}
