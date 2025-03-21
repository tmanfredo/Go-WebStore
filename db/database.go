package db

import (
	"database/sql"
    "fmt"
    "math"
    "go-store/types"
)

func GetUserSecurity (connection *sql.DB, username string, password string) (string, int, error){
    
    stmt, err := connection.Prepare("SELECT first_name, role FROM users WHERE email = ? AND password = ?")
    if err != nil {
        fmt.Sprintf("Error 1: %s", err)
        return "",0, err
    }
    defer stmt.Close()

    rows, err := stmt.Query(username, password)
        if err != nil {
            fmt.Printf("error2: %s", err)
            return "", 0, err
        }
    defer rows.Close()
    if rows.Next() {
        var security int
        var first_name string
        err := rows.Scan(&first_name, &security)
            if err != nil {
                fmt.Printf("error3: %s",err)
                return "", 0, err
            }
        
        return first_name, security, nil
    }
   
    return "", 0, nil

}

func SearchCustomers (connection *sql.DB, input string, searchTerm string) ([]types.Customer, error){
    
        query := fmt.Sprintf("SELECT * FROM customer WHERE %s LIKE ?", input)
        stmt, err := connection.Prepare(query)
        if err != nil {
            fmt.Printf("error1: %s", err)
            return nil, err
        }
        defer stmt.Close()

        rows, err := stmt.Query(fmt.Sprintf("%%%s%%", searchTerm))
        if err != nil {
            fmt.Printf("error2: %s", err)
            return nil, err
        }
        defer rows.Close()


        var customers []types.Customer

        for rows.Next() {
            var customer types.Customer
            err := rows.Scan(&customer.Id, &customer.First, &customer.Last, &customer.Email)
            if err != nil {
                fmt.Printf("error3: %s",err)
                return nil, err
            }
            customers = append(customers, customer)
        }

        return customers, nil
}

func AddCustomer (connection *sql.DB, first_name string, last_name string, email string){
    stmt, err := connection.Prepare("INSERT INTO customer (first_name, last_name, email) VALUES (?,?,?)")
    if err != nil {
        fmt.Sprintf("Error inserting into customer: %s", err)
    }
    defer stmt.Close()
    stmt.Exec(first_name, last_name, email)
}

func GetAllCustomers (connection *sql.DB) ([]types.Customer, error){
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

    var customers []types.Customer

    for rows.Next() {
        var customer types.Customer
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

func GetCustomerById(connection *sql.DB, id int) (*types.Customer, error) {
    var customer types.Customer
    stmt, err := connection.Prepare("SELECT * FROM customer WHERE id=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    rows, err := stmt.Query(id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // Check if there are any rows returned
    if rows.Next() {
        // Scan the values into the customer struct
        err := rows.Scan(&customer.Id,&customer.First, &customer.Last, &customer.Email)
        if err != nil {
            return nil, err
        }
        // Return the found customer
        return &customer, nil
    }
    
    // If no rows were found, return nil for the customer and nil for the error
    return nil, nil
}


func GetCustomerByEmail (connection *sql.DB, email string) (*types.Customer, error){
    var customer types.Customer
    stmt, err := connection.Prepare("SELECT * FROM customer WHERE email=?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    rows, err := stmt.Query(email)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    // Check if there are any rows returned
    if rows.Next() {
        // Scan the values into the customer struct
        err := rows.Scan(&customer.Id,&customer.First, &customer.Last, &customer.Email)
        if err != nil {
            return nil, err
        }
        // Return the found customer
        return &customer, nil
    }
    
    // If no rows were found, return nil for the customer and nil for the error
    return nil, nil
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

func GetAllOrders(connection *sql.DB) ([]types.Order, error){
	stmt, err := connection.Prepare("SELECT product_id, customer_id, quantity, price, tax, donation, timestamp FROM orders")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []types.Order

    for rows.Next() {
        var (
            productId   int
            customerId  int
            order       types.Order
        )
        err := rows.Scan(&productId, &customerId,&order.Quantity, &order.Price, &order.Tax, &order.Donation, &order.Timestamp)
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

func CheckOrderExistence(connection *sql.DB, product_id int, customer_id int, timestamp int64) (bool, error) {
    stmt, err := connection.Prepare("SELECT product_id, customer_id, timestamp FROM orders WHERE product_id = ? AND customer_id = ? AND timestamp = ?")
    if err != nil {
        return false, err
    }
    defer stmt.Close()

    rows, err := stmt.Query(product_id, customer_id, timestamp)
    if err != nil {
        return false, err
    }
    defer rows.Close()

    if rows.Next() {
        return true, nil;
    }
    return false, nil;
}

func AddOrder (connection *sql.DB, product_id int, customer_id int, quantity int, donation string, timestamp int64) error{
    orderExist, _ := CheckOrderExistence(connection, product_id, customer_id, timestamp)
	if !orderExist {
        product, _ := GetProductById(connection, product_id)
        if product.Instock >= quantity {

            tax := 1.08
            total := float64(quantity)*product.Price*tax

            if donation == "Yes"{
                total = math.Ceil(total)
            }
        
            productUpdate, err := connection.Prepare("UPDATE product SET in_stock = in_stock - ? WHERE id = ?")
            if err != nil {
                return fmt.Errorf("error preparing SQL statement: %w", err)
            }
            defer productUpdate.Close()

            _, err = productUpdate.Exec(quantity, product_id)
            if err != nil {
                return fmt.Errorf("error executing SQL statement: %w", err)
            }
        

            stmt, err := connection.Prepare("INSERT INTO orders (product_id, customer_id, quantity, price, tax, donation, timestamp) VALUES (?,?,?,?,?,?,?)")
            if err != nil {
                return fmt.Errorf("error preparing SQL statement: %w", err)
            }
            defer stmt.Close()
            _, err = stmt.Exec(product_id, customer_id, quantity, product.Price, tax, total, timestamp)
            if err != nil {
                return fmt.Errorf("error executing SQL statement: %w", err)
            }
        }
    }
    return nil
    
}

func GetOrdersByProduct(connection *sql.DB, product_id int) (bool,error){
    stmt, _ := connection.Prepare("SELECT * FROM orders WHERE product_id = ? LIMIT 1")
    defer stmt.Close()

    rows, _ := stmt.Query(product_id)
    defer rows.Close()

    if rows.Next() {
       return true, nil;
    }

	    return false, nil
    
}

func GetAllProducts(connection *sql.DB) ([]types.Product, error){
	stmt, err := connection.Prepare("SELECT * FROM product")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []types.Product

    for rows.Next() {
        var product types.Product
        err := rows.Scan(&product.Id, &product.Name, &product.Image, &product.Price, &product.Instock, &product.Inactive)
        if err != nil {
            return nil, err
        }
        products = append(products, product)
    }

	return products, nil
}

func GetProductById (connection *sql.DB, product_id int) (*types.Product, error){
    var product types.Product
    stmt, err := connection.Prepare("SELECT * FROM product WHERE id=?")
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
        rows.Scan(&product.Id,&product.Name, &product.Image, &product.Price, &product.Instock, &product.Inactive)
    }

    return &product, nil
}

func GetProductByName (connection *sql.DB, product_name string) (*types.Product, error){
    var product types.Product
    stmt, err := connection.Prepare("SELECT * FROM product WHERE product_name LIKE ?")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

	rows, err := stmt.Query(product_name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    
    if rows.Next() {
        rows.Scan(&product.Id,&product.Name, &product.Image, &product.Price, &product.Instock, &product.Inactive)
    }

    return &product, nil
}

func CreateProduct (connection *sql.DB, name string, image string, quantity int, price float64, inactive int) (error){
    stmt, err := connection.Prepare("INSERT INTO product (product_name, image_name, price, in_stock, inactive) VALUES (?,?,?,?,?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

	_, err = stmt.Query(name, image, price, quantity, inactive)
    if err != nil {
        return err
    }
    defer stmt.Close()

    return nil
}

func UpdateProduct (connection *sql.DB, id int, name string, image string, quantity int, price float64, inactive int) (error){
    stmt, err := connection.Prepare("UPDATE product SET product_name = ?, image_name = ?, price = ?, in_stock = ?, inactive = ? WHERE id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

	_, err = stmt.Query(name, image, price, quantity, inactive, id)
    if err != nil {
        return err
    }
    defer stmt.Close()

    return nil
}

func DeleteProduct (connection *sql.DB, id int) (error){
    stmt, err := connection.Prepare("DELETE FROM product WHERE id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

	_, err = stmt.Query(id)
    if err != nil {
        return err
    }
    defer stmt.Close()

    return nil
}