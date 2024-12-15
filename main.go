package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func connectDatabase() (*sql.DB, error) {
	connect := os.Getenv("CONNECT")
	db, err := sql.Open("mysql", connect)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected Successfully")
	return db, nil
}

type User struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	UserType  string `json:"user_type"`
}

type Product struct {
	ProductID          int     `json:"product_id"`
	FarmerID           int     `json:"farmer_id"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
	Category           string  `json:"category"`
	Quantity           int     `json:"quantity"`
	Price              float64 `json:"price"`
	ImageURL           string  `json:"image_url"`
	Location           string  `json:"location"`
	Status             string  `json:"status"`
	ProductLife        int     `json:"product_life"`
}

type Order struct {
	OrderID     int     `json:"order_id"`
	BuyerID     int     `json:"buyer_id"`
	OrderDate   string  `json:"order_date"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

type OrderItem struct {
	OrderItemID int     `json:"order_item_id"`
	OrderID     int     `json:"order_id"`
	ProductID   int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

func register(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var newUser User
	registerQuery := os.Getenv("REGISTER_QUERY")
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err = db.Exec(registerQuery, newUser.Username, newUser.Password, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Phone, newUser.Address, newUser.UserType)
	if err != nil {
		http.Error(w, "Error during registration", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("hai hhhhs")
	var credentials User
	var storedPassword string
	fmt.Println("hai")
	loginQuery := os.Getenv("LOGIN_QUERY")
	err := json.NewDecoder(r.Body).Decode(&credentials)
	fmt.Println(loginQuery)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = db.QueryRow(loginQuery, credentials.Username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Unknown User", http.StatusNotFound)
		return
	}
	fmt.Println(storedPassword, credentials.Password)
	if storedPassword == credentials.Password {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func productsDisplay(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	productsQuery := os.Getenv("PRODUCTS_QUERY")
	rows, err := db.Query(productsQuery)
	if err != nil {
		http.Error(w, "Error retrieving products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var productsList []Product
	for rows.Next() {
		var pid, fid, plife int
		var pname, imgurl, status, location string
		var price float64
		err = rows.Scan(&pid, &pname, &fid, &price, &imgurl, &location, &status, &plife)
		if err != nil {
			http.Error(w, "Error scanning products", http.StatusInternalServerError)
			return
		}
		productsList = append(productsList, Product{
			ProductID:   pid,
			ProductName: pname,
			FarmerID:    fid,
			Price:       price,
			ImageURL:    imgurl,
			Location:    location,
			Status:      status,
			ProductLife: plife,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(productsList)
}

func productsDisplayFarmerID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")
	farmerID, _ := strconv.Atoi(parts[4])

	switch r.Method {
	case "GET":
		farmerProductsQuery := os.Getenv("FARMERS_PRODUCT_QUERY")
		rows, err := db.Query(farmerProductsQuery, farmerID)
		if err != nil {
			http.Error(w, "Error retrieving products for farmer", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var farmerProductsList []Product
		for rows.Next() {
			var pid, fid, plife, quantity int
			var pname, imgurl, status, location, pdesc, cat string
			var price float64
			err = rows.Scan(&pid, &fid, &pname, &pdesc, &cat, &quantity, &price, &imgurl, &location, &status, &plife)
			if err != nil {
				http.Error(w, "Error scanning farmer's products", http.StatusInternalServerError)
				return
			}
			farmerProductsList = append(farmerProductsList, Product{
				ProductID:          pid,
				FarmerID:           fid,
				ProductName:        pname,
				ProductDescription: pdesc,
				Category:           cat,
				Quantity:           quantity,
				Price:              price,
				ImageURL:           imgurl,
				Location:           location,
				Status:             status,
				ProductLife:        plife,
			})
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(farmerProductsList)
	case "POST":
		var newProduct Product
		newProductQuery := os.Getenv("NEW_PRODUCT_QUERY")
		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		_, err = db.Exec(newProductQuery, newProduct.FarmerID, newProduct.ProductName, newProduct.ProductDescription, newProduct.Category, newProduct.Quantity, newProduct.Price, newProduct.ImageURL, newProduct.Location, newProduct.Status, newProduct.ProductLife)
		if err != nil {
			http.Error(w, "Error inserting new product", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "PUT":
		var updateProduct Product
		updateProductQuery := os.Getenv("UPDATE_PRODUCT_QUERY")
		err := json.NewDecoder(r.Body).Decode(&updateProduct)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		_, err = db.Exec(updateProductQuery, updateProduct.ProductName, updateProduct.ProductDescription, updateProduct.Category, updateProduct.Quantity, updateProduct.Price, updateProduct.ImageURL, updateProduct.Location, updateProduct.Status, updateProduct.ProductLife, farmerID)
		if err != nil {
			http.Error(w, "Error updating product", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case "DELETE":
		deleteProductQuery := os.Getenv("DELETE_PRODUCT_QUERY")
		_, err := db.Exec(deleteProductQuery, farmerID)
		if err != nil {
			http.Error(w, "Error deleting product", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ordersByBuyer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")
	buyerID, _ := strconv.Atoi(parts[4])

	buyerOrderQuery := os.Getenv("BUYER_ORDER_QUERY")
	rows, err := db.Query(buyerOrderQuery, buyerID)
	if err != nil {
		http.Error(w, "Error retrieving orders for buyer", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var buyerOrdersList []Order
	for rows.Next() {
		var oid, bid int
		var orderDate, status string
		var totalAmount float64
		err = rows.Scan(&oid, &bid, &orderDate, &totalAmount, &status)
		if err != nil {
			http.Error(w, "Error scanning orders", http.StatusInternalServerError)
			return
		}
		buyerOrdersList = append(buyerOrdersList, Order{
			OrderID:     oid,
			BuyerID:     bid,
			OrderDate:   orderDate,
			TotalAmount: totalAmount,
			Status:      status,
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buyerOrdersList)
}

func main() {
	db, err := connectDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/v1/users/register", func(w http.ResponseWriter, r *http.Request) {
		register(db, w, r)
	})
	http.HandleFunc("/v1/users/login", func(w http.ResponseWriter, r *http.Request) {
		login(db, w, r)
	})
	http.HandleFunc("/v1/products", func(w http.ResponseWriter, r *http.Request) {
		productsDisplay(db, w, r)
	})
	http.HandleFunc("/v1/products/farmer/", func(w http.ResponseWriter, r *http.Request) {
		productsDisplayFarmerID(db, w, r)
	})
	http.HandleFunc("/v1/orders/", func(w http.ResponseWriter, r *http.Request) {
		ordersByBuyer(db, w, r)
	})

	// Allow CORS for all routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
