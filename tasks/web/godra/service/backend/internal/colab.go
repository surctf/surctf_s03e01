package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var (
	COLAB_FLAG = os.Getenv("COLAB_FLAG")
)

func (hc *HandlerContext) GetColabFlag(c *gin.Context) {
	userInterface, _ := c.Get("user")
	user := userInterface.(User)

	var purchases []Purchase
	if err := hc.DB.Where("user_id = ?", user.ID).Find(&purchases).Error; err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": InternalServerError})
		return
	}

	if len(purchases) < 20 {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Чтобы получить флаг, нужно совершить 20 любых покупок"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"flag": COLAB_FLAG})
}
