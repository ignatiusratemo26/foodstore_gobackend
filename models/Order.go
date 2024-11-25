package models

import (
	"time"

	"gorm.io/gorm"
)

type LatLng struct {
	Lat string `gorm:"type:varchar(50)"`
	Lng string `gorm:"type:varchar(50)"`
}

type OrderItem struct {
	Food     Food    `gorm:"embedded"`
	Price    float64 `gorm:"not null"`
	Quantity int     `gorm:"not null"`
}

type Order struct {
	ID            string         `gorm:"type:varchar(24);primaryKey"`
	Name          string         `gorm:"type:varchar(100);not null"`
	Address       string         `gorm:"type:varchar(255);not null"`
	AddressLatLng LatLng         `gorm:"embedded"`
	TotalPrice    float64        `gorm:"not null"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID"`
	Status        string         `gorm:"type:varchar(100);not null"`
	UserID        string         `gorm:"type:varchar(24);not null"`
	PaymentID     string         `gorm:"type:varchar(100)"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
