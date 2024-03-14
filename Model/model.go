package Model

type Users struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Products struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}
type Transactions struct {
	ID        int `json:"id"`
	UserID    Users
	ProductID Products
	Quantity  int `json:"quantity"`
}

// response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Users  `json:"data"`
}
type UsersResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    []Users `json:"data"`
}
type ProductResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Products `json:"data"`
}
type ProductsResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []Products `json:"data"`
}
