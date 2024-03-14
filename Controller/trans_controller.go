package Controller

import (
	m "Week6/Model"

	"log"
	"net/http"
)

func GetAllTrans(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var transaction m.Transactions
	var transactions []m.Transactions
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	SendSuccessResponse(w, 200, "Get user successfully!")
}
