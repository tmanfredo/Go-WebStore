package db

import (
	"database/sql"
    "fmt"
    "math"
    "io/ioutil"
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
    Product_Name string
    Customer_Name string
    Quantity int
    Price float64 
    Tax float64
    Donation float64
    Timestamp int64
}

func ResetDB(connection *sql.DB) error {
    // Read the SQL file
    sqlFile, err := ioutil.ReadFile("assets/sql/schema.sql")
    if err != nil {
        return fmt.Errorf("error reading schema file: %w", err)
    }

    // Split and execute each statement in the SQL file (if there are multiple)
    statements := string(sqlFile)
    _, err = connection.Exec(statements)
    if err != nil {
        return fmt.Errorf("error executing SQL statements: %w", err)
    }

    return nil
}



func AddCustomer (connection *sql.DB, first_name string, last_name string, email string){
    stmt, err := connection.Prepare("INSERT INTO customer (first_name, last_name, email) VALUES (?,?,?)")
    if err != nil {
        fmt.Sprintf("Error inserting into customer: %s", err)
    }
    defer stmt.Close()
    stmt.Exec(first_name, last_name, email)
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
	stmt, err := connection.Prepare("SELECT product_id, customer_id, quantity, price, tax, donation FROM orders")
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
        var (
            productId   int
            customerId  int
            order       Order
        )
        err := rows.Scan(&productId, &customerId,&order.Quantity, &order.Price, &order.Tax, &order.Donation)
        if err != nil {
            return nil, err
        }
       
        product, err := GetProductById(connection, productId)
        if err != nil {
            return nil, err
        }

        order.Product_Name = product.Name

        customer, err := GetCustomerById(connection, customerId)
        if err != nil {
            return nil, err
        }

        order.Customer_Name = customer.First + " " + customer.Last

        orders = append(orders, order)
    }

	return orders, nil
}

func AddOrder (connection *sql.DB, product_id int, customer_id int, quantity int, donation bool) error{
    product, err := GetProductById(connection, product_id)
    if err != nil {
        return fmt.Errorf("error fetching product: %w", err)
    }
    tax := 1.08
    total := float64(quantity)*product.Price*tax

    if donation {
        total = math.Ceil(total)
    }

    stmt, err := connection.Prepare("INSERT INTO orders (product_id, customer_id, quantity, price, tax, donation, timestamp) VALUES (?,?,?,?,?,?,NOW())")
    if err != nil {
        return fmt.Errorf("error preparing SQL statement: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(product_id, customer_id, quantity, product.Price, tax, total)
    if err != nil {
        return fmt.Errorf("error executing SQL statement: %w", err)
    }
    return nil
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

func GetProductById (connection *sql.DB, product_id int) (*Product, error){
    var product Product
    stmt, err := connection.Prepare("SELECT product_name, image_name, price, in_stock FROM product WHERE id=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

	rows, err := stmt.Query(product_id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    
    if rows.Next() {
        rows.Scan(&product.Name, &product.Image, &product.Price, &product.Instock)
    }

    return &product, nil
}

