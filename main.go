package main

import (
	"fmt"
	"log"
	"net/http"

	"Week6/Controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/users", Controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/v2/users/{id}", Controller.RawGorm).Methods("GET")
	router.HandleFunc("/v1/users", Controller.InsertNewUser).Methods("POST")
	router.HandleFunc("/v2/users", Controller.InsertNewUserGorm).Methods("POST")
	router.HandleFunc("/v1/users/{id}", Controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/v1/users/{id}", Controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/v2/users/{id}", Controller.DeleteUserGorm).Methods("DELETE")
	router.HandleFunc("/v1/users", Controller.Login).Methods("POST")
	router.HandleFunc("/v2/selectusers", Controller.SelectUserGorm).Methods("GET")

	router.HandleFunc("/v1/products", Controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/v2/selectproduct", Controller.SelectProductGorm).Methods("GET")
	router.HandleFunc("/v1/products", Controller.InsertNewProducts).Methods("POST")
	router.HandleFunc("/v1/products/{id}", Controller.UpdateProduct).Methods("PUT")
	router.HandleFunc("/v2/products/{id}", Controller.UpdateProductGorm).Methods("PUT")
	router.HandleFunc("/v1/products/{id}", Controller.DeleteProduct).Methods("DELETE")

	router.HandleFunc("/v1/transactions", Controller.GetAllTrans).Methods("GET")
	router.HandleFunc("/v1/transactions/{id}", Controller.GetDetailUserTrans).Methods("GET")
	router.HandleFunc("/v1/transactions/users/{id}", Controller.GetTransactionsByUserID).Methods("GET")
	router.HandleFunc("/v1/transactions/products/{id}", Controller.DeleteSingleProduct).Methods("DELETE")
	router.HandleFunc("/v2/transactions/user/{id}", Controller.DeleteSingleProductGorm).Methods("DELETE")
	router.HandleFunc("/v1/transactions", Controller.InsertNewProductsandTrans).Methods("POST")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
