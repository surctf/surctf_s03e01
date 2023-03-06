package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hc *HandlerContext) GetProducts(c *gin.Context) {
	var products []Product
	if err := hc.DB.Find(&products).Error; err != nil {
		c.Error(err)
		hc.redirectWithMethod(c, "GET", "/profile")
		return
	}

	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "products.html", gin.H{"User": user, "Products": products})
}
