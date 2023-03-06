package internal

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HandlerContext struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (hc *HandlerContext) redirectWithMethod(c *gin.Context, method string, path string) {
	c.Request.URL.Path = path
	c.Request.Method = method
	hc.Router.HandleContext(c)
}
