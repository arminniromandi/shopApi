package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OTPCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"index;not null" json:"email"`
	Code      string    `gorm:"not null" json:"code"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	ImageURL    string         `json:"image_url"`
	Description string         `json:"description"`
	Products    []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CategoryID      uint           `gorm:"not null;index" json:"category_id"`
	Category        *Category      `json:"category,omitempty"`
	Name            string         `gorm:"not null" json:"name"`
	ImageURL        string         `json:"image_url"`
	Description     string         `json:"description"`
	Price           float64        `gorm:"not null" json:"price"`
	DiscountPercent float64        `gorm:"default:0" json:"discount_percent"`
	IsSpecialOffer  bool           `gorm:"default:false;index" json:"is_special_offer"`
	Weight          string         `json:"weight"`
	Color           string         `json:"color"`
	Comments        []Comment      `gorm:"foreignKey:ProductID" json:"comments,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      *User          `json:"user,omitempty"`
	Content   string         `gorm:"not null" json:"content"`
	Rating    int            `gorm:"default:5" json:"rating"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	ProductID uint           `gorm:"not null" json:"product_id"`
	Product   *Product       `json:"product,omitempty"`
	Quantity  int            `gorm:"default:1" json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
