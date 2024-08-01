package models

type Success struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Authenticated struct {
	Password string `json:"password"`
}

type CreateReview struct {
	Rating  int32  `json:"rating"`
	Comment string `json:"comment"`
}

type UpdateProduct struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Images      string    `json:"images"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	Discount    *Discount `json:"discount"`
}

type Discount struct {
	DiscountPrice float64 `json:"discount_price"`
	Status        bool    `json:"status"`
}

type UpdateUserById struct {
	FullName   string `json:"fullname,omitempty"`
	Username   string `json:"username,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Email      string `json:"email,omitempty"`
	Image      string `json:"image,omitempty"`
	Role       string `json:"role,omitempty"`
	Address    string `json:"address,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

type UpdateUser struct {
	FullName    string `json:"fullname,omitempty"`
	Username    string `json:"username,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Image       string `json:"image,omitempty"`
	NewPasswrod string `json:"new_passwrod,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	Country     string `json:"country,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Password    string `json:"password,omitempty"`
}

type UpdateReview struct {
	ProductId string `json:"product_id,omitempty"`
	Rating    int32  `json:"rating,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

type UpdateReviewById struct {
	UserId    string `json:"user_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`
	Rating    int32  `json:"rating,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

type AddToBasket struct {
	PurchaseDate string  `json:"purchase_date,omitempty"`
	Quantity     int32   `json:"quantity,omitempty"`
	Price        float64 `json:"price,omitempty"`
}