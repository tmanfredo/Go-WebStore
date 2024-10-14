package db

import "database/sql"

type Product struct{
	Name string
	Image string
	Price float64
	InStock int
}

func GetAllProducts(connection *sql.DB) ([]Product, error){
	rows, err := connection.Query("SELECT product_name, image_name, price, in_stock FROM product")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var products []Product
	for rows.Next() {
        var product Product
        rows.Scan(&product.Name, &product.Image, &product.Price, &product.InStock)
		products = append(products, product)
	}
	return products, nil
}
