package database

import (
	"time"
)

type User struct {
	Id            int64          `gorm:"primary_key;auto_increment" json:"id"`
	FirstName     string         `gorm:"size:255;not null" json:"first_name"`
	LastName      string         `gorm:"size:255" json:"last_name"`
	Avatar        string         `gorm:"size:255" json:"avatar"`
	Address       string         `gorm:"size:255" json:"address"`
	Balance       float64        `gorm:"size:255;not null" json:"balance"`
	Email         string         `gorm:"size:255;not null" json:"email"`
	Password      string         `gorm:"size:255;not null" json:"-"`
	Cards         []Cart         `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Orders        []Order        `gorm:"foreignKey:UserId" json:"orders,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserId" json:"notifications,omitempty"`
}

type Category struct {
	Id       int64     `gorm:"primary_key;auto_increment" json:"id,omitempty"`
	Name     string    `gorm:"size:255;not null" json:"name,omitempty"`
	Products []Product `gorm:"foreignKey:CategoryId" json:"products,omitempty"`
}

type Brand struct {
	Id       int64     `gorm:"primary_key;auto_increment" json:"id,omitempty"`
	Name     string    `gorm:"size:255;not null" json:"name,omitempty"`
	Products []Product `gorm:"foreignKey:BrandId" json:"products,omitempty"`
}

type Product struct {
	Id          int64          `gorm:"primary_key;auto_increment" json:"id,omitempty"`
	Name        string         `gorm:"size:255;not null" json:"name,omitempty"`
	Description string         `gorm:"size:255;not null" json:"description,omitempty"`
	Size        int            `gorm:"not null" json:"size,omitempty"`
	Season      string         `gorm:"size:255;not null" json:"season,omitempty"`
	Price       float64        `gorm:"not null" json:"price,omitempty"`
	CategoryId  int64          `gorm:"not null" json:"category_id,omitempty"`
	Category    Category       `gorm:"foreignKey:CategoryId" json:"-"`
	BrandId     int64          `gorm:"not null" json:"brand_id,omitempty"`
	Brand       Brand          `gorm:"foreignKey:BrandId" json:"brand,omitempty"`
	Images      []ProductImage `gorm:"foreignKey:ProductId" json:"images,omitempty"`
	Carts       []Cart         `gorm:"foreignKey:ProductId" json:"product,omitempty"`
	Orders      []Order        `gorm:"foreignKey:ProductId" json:"order,omitempty"`
}

type ProductImage struct {
	Id        int64   `gorm:"primary_key;auto_increment" json:"id,omitempty"`
	ProductId int64   `gorm:"not null" json:"product_id,omitempty"`
	Product   Product `gorm:"foreignKey:ProductId" json:"product,omitempty"`
	FilePath  string  `gorm:"size:255;not null" json:"file_path,omitempty"`
}

type Cart struct {
	Id        int64   `gorm:"primary_key;auto_increment" json:"id,omitempty"`
	UserId    int64   `gorm:"not null" json:"user_id,omitempty"`
	User      User    `gorm:"foreignKey:UserId" json:"user,omitempty"`
	ProductId int64   `gorm:"not null" json:"product_id,omitempty"`
	Product   Product `gorm:"foreignKey:ProductId" json:"product,omitempty"`
	Count     int     `gorm:"not null" json:"count,omitempty"`
}

type Order struct {
	Id        int64     `gorm:"primary_key;auto_increment" json:"id,omitempty"`
	UserId    int64     `gorm:"not null" json:"user_id,omitempty"`
	User      User      `gorm:"foreignKey:UserId" json:"user,omitempty"`
	ProductId int64     `gorm:"not null" json:"product_id,omitempty"`
	Product   Product   `gorm:"foreignKey:ProductId" json:"product,omitempty"`
	Count     int       `gorm:"not null" json:"count,omitempty"`
	Status    string    `gorm:"not null" json:"status,omitempty"`
	Time      time.Time `gorm:"not null" json:"time,omitempty"`
}

type Notification struct {
	Id          int64  `gorm:"primary_key;auto_increment" json:"id,omitempty"`
	Title       string `gorm:"size:255;not null" json:"title,omitempty"`
	Description string `gorm:"size:255;not null" json:"description,omitempty"`
	UserId      int64  `gorm:"not null" json:"user_id,omitempty"`
	User        User   `gorm:"foreignKey:UserId" json:"user,omitempty"`
}
