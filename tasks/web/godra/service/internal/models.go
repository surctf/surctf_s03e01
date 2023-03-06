package internal

type User struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Username string `json:"username" gorm:"column:username;not null"`
	Balance  int    `json:"balance" gorm:"column:balance;"`
}

type Product struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Image       string `json:"image" gorm:"column:image"`
	Cost        int    `json:"cost" gorm:"column:cost"`
	Purchased   bool   `json:"purchased" gorm:"-"`
}

type Purchase struct {
	UserID    int64   `json:"-" gorm:"primaryKey;column:user_id;not null"`
	Username  string  `json:"-" gorm:"primaryKey;column:username;not null"`
	ProductID int     `json:"-" gorm:"primaryKey;column:product_id;not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}
