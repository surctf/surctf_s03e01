package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hc *HandlerContext) GetProfile(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "profile.html", gin.H{"User": user})
}
