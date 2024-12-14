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
	_ "github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func connect_database() (*sql.DB, error) {
	connect := os.Getenv("CONNECT")
	db, err := sql.Open("mysql", connect)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected Successfully")
	return db, nil
}

type users struct {
	User_id    int    `json:"user_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	User_type  string `json:"user_type"`
}

type products struct {
	Product_id          int     `json:"product_id"`
	Farmer_id           int     `json:"farmer_id"`
	Product_name        string  `json:"product_name"`
	Product_description string  `json:"product_description"`
	Category            string  `json:"category"`
	Quantity            int     `json:"quantity"`
	Price               float64 `json:"price"`
	Image_url           string  `json:"image_url"`
	Location            string  `json:"location"`
	Status              string  `json:"status"`
	Product_life        int     `json:"product_life"`
}

type orders struct {
	Order_id     int     `json:"order_id"`
	Buyer_id     int     `json:"buyer_id"`
	Order_date   string  `json:"order_date"`
	Total_amount float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

type order_items struct {
	Order_item_id int     `json:"order_item_id"`
	Order_id      int     `json:"order_id"`
	Product_id    int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Unit_price    float64 `json:"unit_price"`
}

func register(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var newuser users
	register_query := os.Getenv("REGISTER_QUERY")
	err := json.NewDecoder(r.Body).Decode(&newuser)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(register_query, newuser.Username, newuser.Password, newuser.First_name, newuser.Last_name, newuser.Email, newuser.Phone, newuser.Address, newuser.User_type)
	if err != nil {
		panic(err)
	}
}

func login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var credentials users
	var stored_password string
	login_query := os.Getenv("LOGIN_QUERY")
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		panic(err)
	}
	err = db.QueryRow(login_query, credentials.Username).Scan(&stored_password)
	if err != nil {
		http.Error(w, "Unknown User", http.StatusNotFound)
	}
	if stored_password == credentials.Password {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader((http.StatusUnauthorized))
	}
}

func products_display(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	products_query := os.Getenv("PRODUCTS_QUERY")
	rows, err := db.Query(products_query)
	if err != nil {
		panic(err)
	}
	var productslist []products
	for rows.Next() {
		var pid, fid, plife int
		var pname, imgurl, status, location string
		var price float64
		err = rows.Scan(&pid, &pname, &fid, &price, &imgurl, &location, &status, &plife)
		if err != nil {
			panic(err)
		}
		productslist = append(productslist, products{
			Product_id:   pid,
			Product_name: pname,
			Farmer_id:    fid,
			Price:        price,
			Image_url:    imgurl,
			Location:     location,
			Status:       status,
			Product_life: plife,
		})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(productslist)
	defer rows.Close()

}

func products_display_farmerid(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	parts := strings.Split(urlPath, "/")
	farmerID, _ := strconv.Atoi(parts[4])
	switch r.Method {
	case "GET":
		farmer_products_query := os.Getenv("FARMERS_PRODUCT_QUERY")
		rows, err := db.Query(farmer_products_query, farmerID)
		if err != nil {
			panic(err)
		}
		var farmer_products_list []products
		for rows.Next() {
			var pid, fid, plife, quantity int
			var pname, imgurl, status, location, pdesc, cat string
			var price float64
			err = rows.Scan(&pid, &fid, &pname, &pdesc, &cat, &quantity, &price, &imgurl, &location, &status, &plife)
			if err != nil {
				panic(err)
			}
			farmer_products_list = append(farmer_products_list, products{
				Product_id:          pid,
				Farmer_id:           fid,
				Product_name:        pname,
				Product_description: pdesc,
				Category:            cat,
				Quantity:            quantity,
				Price:               price,
				Image_url:           imgurl,
				Location:            location,
				Status:              status,
				Product_life:        plife,
			})

		}

		w.WriteHeader((http.StatusOK))
		json.NewEncoder(w).Encode(farmer_products_list)
		defer rows.Close()
	case "POST":
		var newproduct products
		new_product_query := os.Getenv("NEW_PRODUCT_QUERY")
		err := json.NewDecoder(r.Body).Decode(&newproduct)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec(new_product_query, newproduct.Farmer_id, newproduct.Product_name, newproduct.Product_description, newproduct.Category, newproduct.Quantity, newproduct.Price, newproduct.Image_url, newproduct.Location, newproduct.Status, newproduct.Product_life)
		if err != nil {
			panic(err)
		}
	case "PUT":
		var updateproduct products
		update_product_query := os.Getenv("UPDATE_PRODUCT_QUERY")
		err := json.NewDecoder(r.Body).Decode(&updateproduct)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec(update_product_query, updateproduct.Product_name, updateproduct.Product_description, updateproduct.Category, updateproduct.Quantity, updateproduct.Price, updateproduct.Image_url, updateproduct.Location, updateproduct.Status, updateproduct.Product_life, farmerID)
		if err != nil {
			panic(err)
		}
	case "DELETE":
		delete_product_query := os.Getenv("DELETE_PRODUCT_QUERY")
		_, err := db.Exec(delete_product_query, farmerID)
		if err != nil {
			panic(err)
		}
	}
}

func ordersbybuyer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path
	parts := strings.Split(urlpath, "/")
	buyerID, _ := strconv.Atoi(parts[4])
	buyer_order_query := os.Getenv("BUYER_ORDER_QUERY")
	rows, err := db.Query(buyer_order_query, buyerID)
	if err != nil {
		panic(err)
	}
	var buyer_orders_list []orders
	for rows.Next() {
		var oid, bid int
		var order_date, status string
		var t_amount float64
		err = rows.Scan(&oid, &bid, &order_date, &t_amount, &status)
		if err != nil {
			panic(err)
		}
		buyer_orders_list = append(buyer_orders_list, orders{
			Order_id:     oid,
			Buyer_id:     bid,
			Order_date:   order_date,
			Total_amount: t_amount,
			Status:       status,
		})
	}

	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(buyer_orders_list)
	defer rows.Close()
}

func ordersbyfarmer(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path
	parts := strings.Split(urlpath, "/")
	farmerID, _ := strconv.Atoi(parts[4])
	farmer_order_query := os.Getenv("FARMER_ORDER_QUERY")
	rows, err := db.Query(farmer_order_query, farmerID)
	if err != nil {
		panic(err)
	}
	var farmer_orders_list []orders
	for rows.Next() {
		var oid, bid int
		var order_date, status string
		var t_amount float64
		err = rows.Scan(&oid, &bid, &order_date, &t_amount, &status)
		if err != nil {
			panic(err)
		}
		farmer_orders_list = append(farmer_orders_list, orders{
			Order_id:     oid,
			Buyer_id:     bid,
			Order_date:   order_date,
			Total_amount: t_amount,
			Status:       status,
		})
	}

	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(farmer_orders_list)
	defer rows.Close()
}

func getorders(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path
	parts := strings.Split(urlpath, "/")
	orderID, _ := strconv.Atoi(parts[3])
	order_query := os.Getenv("ORDER_QUERY")
	rows, err := db.Query(order_query, orderID)
	if err != nil {
		panic(err)
	}
	var orders_list []orders
	for rows.Next() {
		var oid, bid int
		var order_date, status string
		var t_amount float64
		err = rows.Scan(&oid, &bid, &order_date, &t_amount, &status)
		if err != nil {
			panic(err)
		}
		orders_list = append(orders_list, orders{
			Order_id:     oid,
			Buyer_id:     bid,
			Order_date:   order_date,
			Total_amount: t_amount,
			Status:       status,
		})
	}
	w.WriteHeader((http.StatusOK))
	json.NewEncoder(w).Encode(orders_list)
	defer rows.Close()
}

func postorders(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var neworder order_items
	new_order_query := os.Getenv("NEW_ORDER_QUERY")
	err := json.NewDecoder(r.Body).Decode(&neworder)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(new_order_query, neworder.Order_id, neworder.Product_id, neworder.Quantity, neworder.Unit_price)
	if err != nil {
		panic(err)
	}

}

func main() {
	db, err := connect_database()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/v1/users/register", func(w http.ResponseWriter, r *http.Request) {
		register(db, w, r)
	})
	http.HandleFunc("/v1/users/login", func(w http.ResponseWriter, r *http.Request) {
		login(db, w, r)
	})
	http.HandleFunc("/v1/products", func(w http.ResponseWriter, r *http.Request) {
		products_display(db, w, r)
	})
	http.HandleFunc("/v1/products/farmers/{farmer_id}", func(w http.ResponseWriter, r *http.Request) {
		products_display_farmerid(db, w, r)
	})
	http.HandleFunc("/v1/orders/buyer/{buyer_id}", func(w http.ResponseWriter, r *http.Request) {
		ordersbybuyer(db, w, r)
	})
	http.HandleFunc("/v1/orders/farmer/{farmer_id}", func(w http.ResponseWriter, r *http.Request) {
		ordersbyfarmer(db, w, r)
	})
	http.HandleFunc("/v1/orders/{order_id}", func(w http.ResponseWriter, r *http.Request) {
		getorders(db, w, r)
	})
	http.HandleFunc("/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		postorders(db, w, r)
	})
	http.ListenAndServe(":8080", nil)
}
