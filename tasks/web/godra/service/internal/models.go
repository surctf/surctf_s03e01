package internal

type User struct {
	ID        int64    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Username  string   `json:"username" gorm:"column:username;not null"`
	Balance   int      `json:"balance" gorm:"column:balance;default:0"`
	Purchases Purchase `json:"solves" gorm:"constraint:OnDelete:CASCADE;"`
}

type Product struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Image       string `json:"image" gorm:"column:image"`
	Cost        int    `json:"cost" gorm:"column:cost"`
}

type Purchase struct {
	UserID    int64   `json:"-" gorm:"column:user_id;primaryKey;not null"`
	ProductID int     `json:"-" gorm:"column:product_id;primaryKey;not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}
