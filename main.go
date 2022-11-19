package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	connector()
	initiator()
	router := mux.NewRouter()
	router.HandleFunc("/get_balance", getBalanceByUserId).Methods("GET")
	router.HandleFunc("/update_balance", updateBalance).Methods("POST")
	router.HandleFunc("/from_user_to_reserve", fromUserToReserve).Methods("POST")
	router.HandleFunc("/from_reserve_to_user", fromReserveToUser).Methods("POST")
	http.Handle("/", router)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
