package Controller

import (
	m "Week6/Model"

	"encoding/json"
	"fmt"

	// "log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const ( // Define maximum length for each form field
	maxProductNameLength = 55
	maxPriceLength       = 3
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"
	name := r.URL.Query()["name"]
	price := r.URL.Query()["price"]

	if name != nil {
		fmt.Println(name[0])
		query = query + " WHERE name= '" + name[0] + "'"
	}

	if price != nil {
		if name[0] != "" {
			query = query + " AND"
		} else {
			query = query + " WHERE"
		}
		query += " price='" + price[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		SendErrorResponse(w, 500, "Failed to execute.")
		return
	}

	var product m.Products
	var products []m.Products
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			SendErrorResponse(w, 404, "Data not found")
			return
		} else {
			products = append(products, product)
		}
	}
	SendSuccessResponse(w, 200, "Get user successfully!")
}

func InsertNewProductsandTrans(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id := r.Form.Get("id")
	userid := r.Form.Get("userid")
	productId := r.Form.Get("productid")
	quantity := r.Form.Get("quantity")

	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to parse form")
		return
	}

	// Validasi jumlah field
	if len(r.Form) >= 4 {
		fmt.Println("Error: Too many fields provided")
		SendErrorResponse(w, 400, "Bad request")
		return
	}

	// cek klo ada field yang kosong
	if userid == "" || productId == "" || quantity == "" {
		fmt.Println("Error: Too many fields provided")
		SendErrorResponse(w, 400, "Bad request")
		return
	}

	// cek di DB dengan ID produk
	var count int
	tx, err := db.Begin()
	err = tx.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", productId).Scan(&count)
	if err != nil {
		SendErrorResponse(w, 500, "Failed to execute.")
		return
	}

	// Jika produk belum ada, tambahkan produk baru
	if count == 0 {
		_, err = tx.Exec("INSERT INTO products (id, name, price) VALUES (?, '', 0)", productId)
		if err != nil {
			SendErrorResponse(w, 505, "Failed to insert data.")
			return
		}
	}

	// Insert transaksi baru
	_, err = tx.Exec("INSERT INTO transactions (id, userID, productID, quantity) VALUES (?, ?, ?, ?)", id, userid, productId, quantity)
	if err != nil {
		SendErrorResponse(w, 505, "Failed to insert data.")
		return
	}

	SendSuccessResponse(w, 200, "Data inserted successfully")
}

func InsertNewProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id := r.Form.Get("id")
	name := r.Form.Get("name")
	price := r.Form.Get("price")

	var err error
	err = r.ParseForm()
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
	if len(name) > maxProductNameLength || len(price) > maxPriceLength {
		fmt.Println("Error: Field length exceeds maximum allowed")
		SendErrorResponse(w, 400, "Bad request")
		return
	}

	if name == "" || price == "" {
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

	query := "INSERT INTO products (id, name, price) VALUES (?, ?, ?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		SendErrorResponse(w, 505, "Failed to insert data.")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, price)
	if err != nil {
		SendErrorResponse(w, 500, "Internal server error")
		return
	}

	SendSuccessResponse(w, 200, "User inserted successfully!")
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// Parsing data formulir
	err := r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to parse form")
		return
	}

	name := r.Form.Get("name")
	price := r.Form.Get("price")
	id := mux.Vars(r)["id"]

	// Memperbarui data di database
	_, err = db.Exec("UPDATE products SET name=?, price=? WHERE id=?", name, price, id)
	if err != nil {
		SendErrorResponse(w, 500, "Internal server error")
		return
	}

	SendSuccessResponse(w, 200, "Product updated successfully!")
}

func UpdateProductGorm(w http.ResponseWriter, r *http.Request) {
	db, err := connectForGorm()
	if err != nil {
		SendErrorResponse(w, 500, "Error connecting to database")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	price := r.URL.Query().Get("price")
	if id == "" {
		SendErrorResponse(w, 505, "Bad request: Missing ID")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		SendErrorResponse(w, 505, "Bad request: Invalid ID")
		return
	}
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		SendErrorResponse(w, 505, "Bad request: Invalid price")
		return
	}

	var product m.Products
	db.Find(&product, &idInt)
	product.ID = idInt
	product.Price = priceInt

	if err := db.Save(&product).Error; err != nil {
		SendErrorResponse(w, 500, "Failed to update data")
		return
	}
	SendSuccessResponse(w, 200, "Products updated successfully!")
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id := mux.Vars(r)["id"]

	var err error
	err = r.ParseForm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to parse form")
		return
	}

	stmt, err := db.Prepare("DELETE FROM products WHERE id=?")
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

func DeleteSingleProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		SendErrorResponse(w, 505, "Failed to convert")
		return
	}

	_, err = db.Exec("DELETE FROM transactions WHERE productID = ?", idInt)
	if err != nil {
		SendErrorResponse(w, 505, "Failed to delete data.")
		return
	}

	stmt, err := db.Prepare("DELETE FROM products WHERE id=?")
	if err != nil {
		SendErrorResponse(w, 505, "Failed to delete data.")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(idInt)
	if err != nil {
		SendErrorResponse(w, 500, "Internal server error")
		return
	}

	SendSuccessResponse(w, 200, "Product updated successfully!")
}

func DeleteSingleProductGorm(w http.ResponseWriter, r *http.Request) {
	db, err := connectForGorm()
	if err != nil {
		SendErrorResponse(w, 500, "Error connecting to database")
		return
	}

	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		SendErrorResponse(w, 505, "Failed to convert")
		return
	}
	var trans m.Transactions
	result := db.Where("productID = ?", &idInt).Delete(&trans)
	err = result.Error
	if err != nil {
		SendErrorResponse(w, 505, "Failed to delete data")
		return
	}

	var product m.Products
	result2 := db.Where("id = ?", &idInt).Delete(&product)
	err = result2.Error
	if err != nil {
		SendErrorResponse(w, 505, "Failed to delete data")
		return
	}
	SendSuccessResponse(w, 200, "Product deleted successfully!")

}

func SelectProductGorm(w http.ResponseWriter, r *http.Request) {
	db, err := connectForGorm()
	if err != nil {
		SendErrorResponse(w, 500, "Failed to establish a connection to the database")
		return
	}

	var products []m.Products
	queryResult := db.Last(&products)
	if queryResult.Error != nil {
		SendErrorResponse(w, 500, "Error retrieving product data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Success"
	response.Data = products
	json.NewEncoder(w).Encode(response)
}
