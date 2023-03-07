package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hc *HandlerContext) GetProfile(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(User)
	var purchases []Purchase
	if err := hc.DB.Where("user_id = ?", user.ID).Preload("Product").Find(&purchases).Error; err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
	}
	c.HTML(http.StatusOK, "profile.html", gin.H{"User": user, "Purchases": purchases})
}

func (hc *HandlerContext) DeleteProfile(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(User)

	if err := hc.DB.Where("id = ?", user.ID).Delete(&user).Error; err != nil {
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "internal_error.html", nil)
		return
	}

	hc.redirectWithMethod(c, "GET", "/signup")
}
