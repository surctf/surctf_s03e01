package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type TgUserUnmarshal struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func isValidInitData(initData url.Values) bool {
	_, ok1 := initData["auth_date"]
	_, ok2 := initData["query_id"]
	_, ok3 := initData["user"]
	_, ok4 := initData["hash"]
	return ok1 && ok2 && ok3 && ok4
}

func TgAuthMiddleware(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		initData := c.Query("initData")

		if len(initData) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": "не найдены телеграм данные (tg_data_not_found)"})
			return
		}

		initDataDecoded, err := url.ParseQuery(initData)
		if err != nil || !isValidInitData(initDataDecoded) {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": "неправильный формат 'initData' (wrong_initdata_format)"})
			return
		}

		h := hmac.New(sha256.New, secret)
		h.Write([]byte(fmt.Sprintf("auth_date=%s\nquery_id=%s\nuser=%s",
			initDataDecoded["auth_date"][0],
			initDataDecoded["query_id"][0],
			initDataDecoded["user"][0],
		)))

		if hex.EncodeToString(h.Sum(nil)) != initDataDecoded["hash"][0] {
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"error": "ошибка при проверке подлинности initData (fake_init_data)"})
			return
		}

		var user TgUserUnmarshal
		json.Unmarshal([]byte(initDataDecoded["user"][0]), &user)

		c.Set("tgID", user.ID)
		c.Set("tgUsername", user.Username)
		c.Next()
	}
}

func JSONLogMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := time.Now().Sub(start).Milliseconds()

		entry := log.WithFields(logrus.Fields{
			"client_ip": c.ClientIP(),
			"duration":  duration,
			"method":    c.Request.Method,
			"path":      c.Request.RequestURI,
			"status":    c.Writer.Status(),
			"referrer":  c.Request.Referer(),
		})

		if c.Writer.Status() >= 500 || len(c.Errors.String()) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Infof("Successfully performed: [%s] %s", c.Request.Method, c.Request.URL.Path)
		}
	}
}
