package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Food struct {
	gorm.Model
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `gorm:"type:varchar(100);not null"`
	Price    float64            `gorm:"not null"`
	Tags     []string           `gorm:"type:varchar(100)"`
	Favorite bool               `gorm:"default:false"`
	Stars    int                `gorm:"default:3"`
	ImageUrl string             `gorm:"type:varchar(255);not null"`
	Origins  []string           `gorm:"type:varchar(100);not null"`
	CookTime string             `gorm:"type:varchar(100);not null"`
}
