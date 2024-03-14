package Controller

import (
	m "Week6/Model"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	maxUserNameLength = 55
	maxAgeLength      = 3
	maxAddressLength  = 150
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM users"
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]

	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name= '" + name[0] + "'"
	}

	if age != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " age='" + age[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		SendErrorResponse(w, 500, "Failed to execute.")
		return
	}

	var user m.Users
	var users []m.Users
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.Email, &user.Password); err != nil {
			SendErrorResponse(w, 404, "Data not found")
			return
		} else {
			users = append(users, user)
		}
	}
	// SendSuccessResponse(w, 200, "Get user successfully!")
	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id := r.Form.Get("id")
	name := r.Form.Get("name")
	age := r.Form.Get("age")
	address := r.Form.Get("address")

	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to parse form")
		return
	}

	// cek klo kelebihan fieldnya
	if len(r.Form) > 3 {
		fmt.Println("Error: Too many fields provided")
		SendErrorResponse(w, 400, "Bad request")
		return
	}

	// Validate individual field lengths
	if len(name) > maxUserNameLength || len(age) > maxAgeLength || len(address) > maxAddressLength {
		fmt.Println("Error: Field length exceeds maximum allowed")
		SendErrorResponse(w, 400, "Bad request")
		return
	}

	if name == "" || age == "" || address == "" {
		fmt.Println("Error: Incomplete data provided")
		SendErrorResponse(w, 400, "Bad request")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		SendErrorResponse(w, 400, "Bad request")
		return
	}
	defer tx.Rollback()

	query := "INSERT INTO users (id, name, age, address) VALUES (?, ?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		SendErrorResponse(w, 505, "Failed to insert data.")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, age, address)
	if err != nil {
		SendErrorResponse(w, 500, "Internal server error")
		return
	}

	SendSuccessResponse(w, 200, "User inserted successfully!")
}

func InsertNewUserGorm(w http.ResponseWriter, r *http.Request) {
	db, err := connectForGorm()
	if err != nil {
		SendErrorResponse(w, 500, "Error connecting to database")
		return
	}

	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")
	address := r.URL.Query().Get("address")
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	agestr, err := strconv.Atoi(age)
	if err != nil {
		SendErrorResponse(w, 505, "Failed to convert")
		return
	}

	users := m.Users{
		Name:     name,
		Age:      agestr,
		Address:  address,
		Email:    email,
		Password: password,
	}

	result := db.Create(&users)
	err = result.Error
	if err != nil {
		SendErrorResponse(w, 404, "Data not found")
	}
	SendSuccessResponse(w, 200, "User inserted successfully!")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	name := r.Form.Get("name")
	age := r.Form.Get("age")
	address := r.Form.Get("address")
	id := mux.Vars(r)["id"] // Mengambil nilai ID dari URL

	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to parse form")
		return
	}

	stmt, err := db.Prepare("UPDATE users SET name=?, age=?, address=? WHERE id=?")
	if err != nil {
		SendErrorResponse(w, 505, "Failed to update data.")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, age, address, id)
	if err != nil {
		SendErrorResponse(w, 500, "Internal server error")
		return
	}

	SendSuccessResponse(w, 200, "User updated successfully!")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id := mux.Vars(r)["id"]

	var err error
	err = r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to parse form")
		return
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		SendErrorResponse(w, 505, "Failed to delete data")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		SendErrorResponse(w, 500, "Internal server error")
		return
	}

	SendSuccessResponse(w, 200, "User deleted successfully!")
}

func DeleteUserGorm(w http.ResponseWriter, r *http.Request) {
	db, err := connectForGorm()
	if err != nil {
		SendErrorResponse(w, 500, "Error connecting to database")
		return
	}

	id := mux.Vars(r)["id"]

	var user m.Users
	result := db.Where("id = ?", &id).Delete(&user)
	err = result.Error
	if err != nil {
		SendErrorResponse(w, 505, "Failed to delete data")
		return
	}

	SendSuccessResponse(w, 200, "User deleted successfully!")
}

func GetDetailUserTrans(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := `SELECT u.ID, u.Name, u.age, u.Address, p.ID, p.Name, p.Price, t.ID, t.Quantity FROM transactions t JOIN users u ON t.UserID = u.ID JOIN products p ON t.ProductID = p.ID`

	userDetailTransRow, err := db.Query(query)
	if err != nil {
		SendErrorResponse(w, 500, "Failed to execute.")
		return
	}

	var userDetailTransactions []m.Transactions
	for userDetailTransRow.Next() {
		//initialiizing each class
		var user m.Users
		var product m.Products
		var trans m.Transactions

		if err := userDetailTransRow.Scan(
			&user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &trans.ID, &trans.Quantity); err != nil {
			SendErrorResponse(w, 500, "Couldn't get user details")
			return
		} else {
			trans.UserID = user
			trans.ProductID = product

			userDetailTransactions = append(userDetailTransactions, trans)
		}
	}
	SendSuccessResponse(w, 200, "Get data successfully!")
}

func GetTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		SendErrorResponse(w, 505, "Failed to convert")
		return
	}

	query := fmt.Sprintf(`SELECT u.ID, u.Name, u.Age, u.Address, p.ID, p.Name, p.Price, t.ID, t.Quantity
        FROM transactions t
        JOIN users u ON t.UserID = u.ID
        JOIN products p ON t.ProductID = p.ID
        WHERE u.ID = %d`, userID)

	userDetailTransRow, err := db.Query(query)
	if err != nil {
		SendErrorResponse(w, 500, "Failed to execute.")
		return
	}
	defer userDetailTransRow.Close()

	var userDetailTransactions []m.Transactions
	for userDetailTransRow.Next() {
		var transaction m.Transactions
		var user m.Users
		var product m.Products

		if err := userDetailTransRow.Scan(
			&user.ID, &user.Name, &user.Age, &user.Address,
			&product.ID, &product.Name, &product.Price,
			&transaction.ID, &transaction.Quantity); err != nil {
			SendErrorResponse(w, 500, "Couldn't get user details")
			return
		}

		transaction.UserID = user
		transaction.ProductID = product

		userDetailTransactions = append(userDetailTransactions, transaction)
	}

	SendSuccessResponse(w, 200, "Get data successfully!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		SendErrorResponse(w, 500, "Email and Password must be filled.")
		return
	}

	var (
		dbPassword string
		userID     int
		name       string
		age        int
		address    string
	)

	query := "SELECT id, password, name, age, address FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&userID, &dbPassword, &name, &age, &address)
	if err != nil {
		if err == sql.ErrNoRows {
			SendErrorResponse(w, 401, "Unauthorized")
			return
		}
		SendErrorResponse(w, 500, "Failed to execute.")
		return
	}

	if password != dbPassword {
		SendErrorResponse(w, 505, "Wrong password.")
		return
	}

	platform := r.Header.Get("platform")
	fmt.Fprintf(w, "Success login from %s", platform)
}

func RawGorm(w http.ResponseWriter, r *http.Request) {
	db, err := connectForGorm()
	if err != nil {
		SendErrorResponse(w, 500, "Error connecting to database")
		return
	}
	id := mux.Vars(r)["id"]

	var user []m.Users
	db.Raw("SELECT name, address FROM users WHERE id = ?", &id).Scan(&user)
	if err != nil {
		SendErrorResponse(w, 500, "Failed to select data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = user
	json.NewEncoder(w).Encode(response)

}

func SelectUserGorm(w http.ResponseWriter, r *http.Request) {
	db, err := connectForGorm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to establish a connection to the database")
		return
	}

	var users []m.Users
	queryResult := db.Last(&users)
	if queryResult.Error != nil {
		SendErrorResponse(w, 500, "Error retrieving product data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}
