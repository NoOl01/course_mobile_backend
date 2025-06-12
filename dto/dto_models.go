package dto

type RegisterDto struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SendOtpDto struct {
	Email string `json:"email"`
}

type ResetPasswordDto struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ResetPasswordWithCodeDto struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}

type ProfileDto struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Address   string `json:"address,omitempty"`
	Email     string `json:"email,omitempty"`
}

type CartDto struct {
	ProductId int64 `json:"product_id"`
	Count     int   `json:"count"`
}

type BuyProductDto struct {
	ProductId int64 `json:"product_id"`
	Count     int   `json:"count"`
}

type ProductWithImageResult struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	CategoryId int64   `json:"category_id"`
	BrandId    int64   `json:"brand_id"`
	Image      string  `json:"image"`
	IsLiked    bool    `json:"is_liked"`
	InCart     bool    `json:"in_cart"`
}

type ProductWithCount struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Image   string  `json:"image"`
	IsLiked bool    `json:"is_liked"`
	InCart  bool    `json:"in_cart"`
	Count   int     `json:"count"`
}

type ProductInfoResult struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	Brand       string   `json:"brand"`
	Images      []string `json:"image"`
	IsLiked     bool     `json:"is_liked"`
	InCart      bool     `json:"in_cart"`
}

type SettingMoneyAction struct {
	Action string  `json:"action"`
	Money  float64 `json:"money"`
}
