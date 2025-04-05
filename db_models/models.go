package db_models

type User struct {
	Id          int64  `gorm:"primary_key;auto_increment" json:"user_id,omitempty"`
	FirstName   string `gorm:"size:255" json:"first_name,omitempty"`
	LastName    string `gorm:"size:255" json:"last_name,omitempty"`
	Email       string `gorm:"size:255;not null" json:"email,omitempty"`
	Address     string `gorm:"size:255" json:"address,omitempty"`
	PhoneNumber string `gorm:"size:255" json:"phone_number,omitempty"`
	Password    string `gorm:"size:255;not null" json:"password,omitempty"`
}

type Product struct {
	Id            int64          `gorm:"primary_key;auto_increment" json:"product_id,omitempty"`
	Name          string         `gorm:"size:255;not null" json:"name,omitempty"`
	Description   string         `gorm:"not null" json:"description,omitempty"`
	ProductImages ProductImage   `gorm:"foreignKey:ProductId" json:"product_images,omitempty"`
	Price         float64        `gorm:"not null" json:"price,omitempty"`
	CategoryId    int64          `gorm:"not null" json:"category_id,omitempty"`
	Category      *Category      `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	ProductFilter *ProductFilter `gorm:"foreignKey:ProductId" json:"product_filter,omitempty"`
}

type ProductImage struct {
	Id        int64    `gorm:"primary_key;auto_increment" json:"product_image_id,omitempty"`
	ProductId int64    `gorm:"not null" json:"product_id,omitempty"`
	Product   *Product `gorm:"foreignKey:ProductId" json:"product,omitempty"`
}

type ProductFilter struct {
	Id        int64   `gorm:"primary_key;auto_increment" json:"product_filter_id,omitempty"`
	ProductId int64   `gorm:"not null" json:"product_id,omitempty"`
	FilterId  int64   `gorm:"not null" json:"filter_id,omitempty"`
	Product   Product `gorm:"foreignKey:ProductId" json:"product,omitempty"`
	Filter    Filter  `gorm:"foreignKey:FilterId" json:"filter,omitempty"`
}

type Category struct {
	Id       int64     `gorm:"primary_key;auto_increment" json:"category_id,omitempty"`
	Name     string    `gorm:"size:255;not null" json:"category_name,omitempty"`
	Products []Product `gorm:"foreignKey:CategoryId" json:"products,omitempty"`
}

type Filter struct {
	Id             int64           `gorm:"primary_key;auto_increment" json:"filter_id,omitempty"`
	Name           string          `gorm:"size:255;not null" json:"filter_name,omitempty"`
	ProductFilters []ProductFilter `gorm:"foreignKey:FilterId" json:"product_filters,omitempty"`
}
