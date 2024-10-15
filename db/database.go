package db

import (
	"database/sql"
)

type Product struct{
	Name string
	Image string
	Price float64
	InStock int
}

type Customer struct{
	First string
	Last string
	Email string
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
        err := rows.Scan(&product.Name, &product.Image, &product.Price, &product.InStock)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

	return products, nil
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
