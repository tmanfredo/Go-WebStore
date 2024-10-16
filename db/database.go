package db

import (
	"database/sql"
    "time"
    "fmt"
)

type Product struct{
	Name string
	Image string
	Price float64
	Instock int
}

type Customer struct{
	First string
	Last string
	Email string
}

type Order struct{
    Product_Id int
    Customer_Id int
    Quantity int
    Price float64 
    Tax float64
    Donation float64
    Timestamp time.Time
}

func AddCustomer (connection *sql.DB, first_name string, last_name string, email string){
    stmt, err := connection.Prepare("INSERT INTO customer (first_name, last_name, email) VALUES (?,?,?)")
    if err != nil {
        fmt.Sprintf("Error inserting into customer: %s", err)
    }
    defer stmt.Close()
    stmt.Query(first_name, last_name, email)
}

func GetAllCustomers (connection *sql.DB) ([]Customer, error){
	stmt, err := connection.Prepare("SELECT first_name, last_name, email FROM customer")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

	rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var customers []Customer

    for rows.Next() {
        var customer Customer
        err := rows.Scan(&customer.First, &customer.Last, &customer.Email)
        if err != nil {
            return nil, err
        }
        customers = append(customers, customer)
    }

	return customers, nil
}

func NumOfCustomers (connection *sql.DB) (int, error){
    stmt, err := connection.Prepare("SELECT COUNT(*) FROM customer")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

	rows, err := stmt.Query()
    if err != nil {
        return 0, err
    }
    defer rows.Close()
    
    var num int
    if rows.Next() {
        rows.Scan(&num)
    }

    return num, nil
}

func GetCustomerById (connection *sql.DB, id int) (*Customer, error){
    var customer Customer
    stmt, err := connection.Prepare("SELECT first_name, last_name, email FROM customer WHERE id=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

	rows, err := stmt.Query(id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    
    if rows.Next() {
        rows.Scan(&customer.First, &customer.Last, &customer.Email)
    }

    return &customer, nil
}

func GetCustomerByEmail (connection *sql.DB, email string) (*Customer, error){
    var customer Customer
    stmt, err := connection.Prepare("SELECT first_name, last_name, email FROM customer WHERE email=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

	rows, err := stmt.Query(email)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    
    if rows.Next() {
        rows.Scan(&customer.First, &customer.Last, &customer.Email)
    }

    return &customer, nil
}

func NumOfOrders (connection *sql.DB) (int, error){
    stmt, err := connection.Prepare("SELECT COUNT(*) FROM orders")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

	rows, err := stmt.Query()
    if err != nil {
        return 0, err
    }
    defer rows.Close()
    
    var num int
    if rows.Next() {
        rows.Scan(&num)
    }

    return num, nil
}

func GetAllOrders(connection *sql.DB) ([]Order, error){
	stmt, err := connection.Prepare("SELECT product_id, customer_id, quantity, price, tax, donation FROM product")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []Order

    for rows.Next() {
        var order Order
        err := rows.Scan(&order.Product_Id, &order.Customer_Id,&order.Quantity, &order.Price, &order.Tax, &order.Donation)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }

	return orders, nil
}

func AddOrder (connection *sql.DB, product_id int, customer_id int, quantity int) {

}

func GetAllProducts(connection *sql.DB) ([]Product, error){
	stmt, err := connection.Prepare("SELECT product_name, image_name, price, in_stock FROM product")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product

    for rows.Next() {
        var product Product
        err := rows.Scan(&product.Name, &product.Image, &product.Price, &product.Instock)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

	return products, nil
}

func GetProductPrice (connection *sql.DB, product_id int)