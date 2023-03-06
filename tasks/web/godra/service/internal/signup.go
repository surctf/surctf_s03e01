package internal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hc *HandlerContext) PostSignUp(c *gin.Context) {
	userID := c.GetInt64("tgID")
	body := c.GetString("requestBody")

	var user User
	if err := json.Unmarshal([]byte(body), &user); err != nil {
		c.Error(err)
		hc.redirectWithMethod(c, "GET", "/signup")
		return
	}

	user.ID = userID
	user.Purchases = nil
	if err := hc.DB.Create(user).Error; err != nil {
		c.Error(err)
		hc.redirectWithMethod(c, "GET", "/profile")
		return
	}

	hc.redirectWithMethod(c, "GET", "/profile")
}

func (hc *HandlerContext) GetSignUp(c *gin.Context) {
	c.Data(http.StatusOK, ContentTypeHTML, []byte("<p>signup form</p>"))
}
