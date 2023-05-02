package internal

import (
	"gorm.io/gorm"
)

type HandlerContext struct {
	DB *gorm.DB
}
