package main

type User struct {
	UserId      int `gorm:"primaryKey" json:"user_id"`
	UserBalance int `json:"user_balance"`
}
type Reserve struct {
	UserId    int `json:"user_id"`
	ServiceId int `json:"service_id"`
	OrderId   int `json:"order_id"`
	Price     int `json:"price"`
}
type JsonResponse struct {
	Type    string `json:"type"`
	Data    []User `json:"data"`
	Message string `json:"message"`
}
