package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func dropWithErrorMessage(w http.ResponseWriter, x string) {
	json.NewEncoder(w).Encode(JsonResponse{Type: "error", Message: x})
}

func getBalanceByUserId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getBalanceByUserId")

	query := r.URL.Query()
	user_id := query.Get("user_id")

	if user_id == "" {
		dropWithErrorMessage(w, "Query param 'user_id' wasn't found")
		return
	}

	var user User
	if err := db.Table("users").Take(&user, user_id).Error; err != nil {
		dropWithErrorMessage(w, "User with id="+user_id+" wasn't found")
		return
	}

	var resultData []User
	resultData = append(resultData, user)

	json.NewEncoder(w).Encode(
		JsonResponse{Type: "success", Data: resultData})
}
func updateBalance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("increaseBalanceByValue")

	var user1 User
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &user1)

	if err != nil {
		dropWithErrorMessage(w, "Wrong parameters")
		return
	}

	var resultData []User

	var user2 User
	if err := db.Table("users").Take(&user2, user1.UserId).Error; err != nil {
		db.Table("users").Create(&user1)
		resultData = append(resultData, user1)
		json.NewEncoder(w).Encode(JsonResponse{Type: "success", Data: resultData, Message: "User was created"})
	} else {
		db.Table("users").Take(&user2, user1.UserId).Update("user_balance", user1.UserBalance+user2.UserBalance)
		resultData = append(resultData, user2)
		json.NewEncoder(w).Encode(JsonResponse{Type: "success", Data: resultData, Message: "User's balance was updated"})
	}
}
func fromUserToReserve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fromUserToReserve")

	var reserve Reserve
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &reserve)

	if err != nil {
		dropWithErrorMessage(w, "Wrong parameters")
		return
	}

	if reserve.Price < 0 {
		dropWithErrorMessage(w, "'price' should be positive")
		return
	}

	var user User
	if err := db.Table("users").Take(&user, reserve.UserId).Error; err != nil {
		dropWithErrorMessage(w, "There is no such user_id in 'users' table")
		return
	}
	if user.UserBalance-reserve.Price < 0 {
		dropWithErrorMessage(w, "User has insufficient funds")
		return
	}

	var count int64
	db.Model(&Reserve{}).Where("service_id = ? AND order_id = ?", reserve.ServiceId, reserve.OrderId).Count(&count)
	if count == 0 {
		db.Table("reserves").Create(&reserve)
	} else {
		dropWithErrorMessage(w,
			"Cannot create new reserve - this combination of order_id and service_id already exist")
		return
	}

	db.Table("users").Take(&user, user.UserId).Update("user_balance", user.UserBalance-reserve.Price)

	json.NewEncoder(w).Encode(JsonResponse{Type: "success",
		Message: "Some money from balance of user with id=" + strconv.Itoa(user.UserId) + " was transferred to reserve, new reserve created"})

}
func fromReserveToUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fromReserveToUser")

	var reserve1 Reserve
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &reserve1)

	if err != nil {
		dropWithErrorMessage(w, "Wrong parameters")
		return
	}

	if reserve1.Price < 0 {
		dropWithErrorMessage(w, "'price' should be positive")
		return
	}

	var reserve2 Reserve
	if err := db.Table("reserves").Where("order_id = ? AND service_id = ?",
		reserve1.OrderId, reserve1.ServiceId).First(&reserve2).Error; err != nil {
		dropWithErrorMessage(w, "No reserves with these service_id and order_id")
		return
	}
	if reserve2.Price < reserve1.Price {
		dropWithErrorMessage(w, "User with id="+strconv.Itoa(reserve1.UserId)+" want too much money")
		return
	}

	var resultMessage = ""

	var user1 User
	if err := db.Table("users").Take(&user1, reserve1.UserId).Error; err == nil {
		db.Table("users").Where("user_id = ?", user1.UserId).Update("user_balance", user1.UserBalance+reserve1.Price)
		resultMessage += "User with id=" + strconv.Itoa(user1.UserId) + " increased his balance"
	} else {
		user1.UserId = reserve1.UserId
		user1.UserBalance = reserve1.Price
		db.Table("users").Create(&user1)
		resultMessage += "User with id=" + strconv.Itoa(user1.UserId) + " was created with some money on balance"
	}

	var user2 User
	if reserve2.Price-reserve1.Price > 0 {
		db.Table("users").Take(&user2, reserve2.UserId)
		db.Table("users").Where("user_id = ?", user2.UserId).Update("user_balance", user2.UserBalance+reserve2.Price-reserve1.Price)
		resultMessage += ", user with id=" + strconv.Itoa(user2.UserId) + " got some money back"
	}

	db.Table("reserves").Where("order_id = ? AND service_id = ?",
		reserve2.OrderId, reserve2.ServiceId).Delete(&reserve2)
	resultMessage += ", reserve was deleted"

	json.NewEncoder(w).Encode(JsonResponse{Type: "success", Message: resultMessage})
}
