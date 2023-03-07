package internal

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
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

	if len(user.Username) < 6 {
		c.HTML(http.StatusNotAcceptable, "signup.html", gin.H{"Error": "юзернейм должен быть длиннее 6 символов"})
		return
	}

	user.ID = userID
	if err := hc.DB.Create(user).Error; err != nil {
		var pqErr *pgconn.PgError
		if ok := errors.As(err, &pqErr); ok && pqErr.Code == "23505" {
			c.HTML(http.StatusNotAcceptable, "signup.html", gin.H{"Error": "скорее всего юзернейм уже кем-то используется или ты уже зареган"})
			return
		}
		c.Error(err)
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{"Error": "что-то пошло не так, попробуй еще раз"})
		return
	}

	hc.redirectWithMethod(c, "GET", "/profile")
}

func (hc *HandlerContext) GetSignUp(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}
