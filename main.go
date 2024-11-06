package main

import (
	"net/http"

	"strconv"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	etag "github.com/pablor21/echo-etag/v4"
	"go-store/templates"
	"go-store/types"
	"database/sql"
	"go-store/db"
	"github.com/go-sql-driver/mysql"
)

func connect() *sql.DB {
	config := mysql.Config{
		User: "tmanfredo",
		Passwd: "031820",
		DBName: "tmanfredo",
	}

	sql, _ := sql.Open("mysql", config.FormatDSN())
	return sql
}

func main() {
	e := echo.New()

	e.Use(etag.Etag())

	e.Static("assets", "./assets")

	
	e.GET("/store", func(ctx echo.Context) error {
		connection := connect()
		storeProducts, _ := db.GetAllProducts(connection);
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(storeProducts)))
	})

	e.GET("/order_entry", func(ctx echo.Context) error {
		connection := connect()
		storeProducts, _ := db.GetAllProducts(connection);
		return Render(ctx, http.StatusOK, templates.OrderEntry(storeProducts))
	})

	e.GET("/search_results", func(ctx echo.Context) error {
		connection := connect()
		input := ctx.QueryParam("field")
		searchTerm := ctx.QueryParam("input")
		customerSearch, _ := db.SearchCustomers(connection, input, searchTerm)
		return Render(ctx, http.StatusOK, templates.UserSearch(customerSearch))
	})


	e.GET("/admin", func(ctx echo.Context) error {
		connection := connect()
		customers, _ := db.GetAllCustomers(connection)
		orders, _ := db.GetAllOrders(connection)
		numOrders, _ := db.NumOfOrders(connection)
		products, _ := db.GetAllProducts(connection)
		return Render(ctx, http.StatusOK, templates.Admin(customers, orders, numOrders, products))
	})

	// Handle the form submission and return the purchase confirmation view
	e.POST("/purchase", func(ctx echo.Context) error {
		connection := connect()
	customer, _ := db.GetCustomerByEmail(connection, ctx.FormValue("email"))
	welcome := ""
	if customer == nil {
		db.AddCustomer(connection, ctx.FormValue("first"), ctx.FormValue("last"), ctx.FormValue("email"))
		customer, _ = db.GetCustomerByEmail(connection, ctx.FormValue("email"))
		welcome = "Looks like you're new! An account has been made for you."
	} else {
		welcome = "Welcome back!"
	}
	timestamp,_ := strconv.Atoi(ctx.FormValue("timestamp"))
	quantity, _ := strconv.Atoi(ctx.FormValue("quantity"))
	product := ctx.FormValue("product")
	dbProduct,_ := db.GetProductByName(connection, product)
	price := dbProduct.Price 
	tax := 1.08
	subtotal :=  (price * float64(quantity))
	total :=  subtotal* tax
		purchase := types.PurchaseInfo{
			Welcome:  welcome,
			First:    customer.First,
			Last:     customer.Last,
			Email:    customer.Email,
			Product:  product,
			Price:    price,
			Quantity: quantity,
			Donate:   ctx.FormValue("donate"),
			Tax:      tax,
			Subtotal: subtotal,
			Total:    total,
		}
		
		//add order but only if it isn't already in there (checked inside of AddOrder)
		db.AddOrder(connection, dbProduct.Id, customer.Id, quantity, ctx.FormValue("donate"), (int64)(timestamp))
		return Render(ctx, http.StatusOK, templates.Base(templates.PurchaseConfirmation(purchase)))
	})

	e.GET("/order_placed", func(ctx echo.Context) error {
		connection := connect()
		customer, _ := db.GetCustomerByEmail(connection, ctx.QueryParam("email"))
		if customer == nil {
			db.AddCustomer(connection, ctx.QueryParam("first"), ctx.QueryParam("last"), ctx.QueryParam("email"))
			customer, _ = db.GetCustomerByEmail(connection, ctx.QueryParam("email"))
		} 
		quantity, _ := strconv.Atoi(ctx.QueryParam("quantity"))
		timestamp,_ := strconv.Atoi(ctx.QueryParam("timestamp"))
		dbProduct,_ := db.GetProductByName(connection, ctx.QueryParam("product"))
		price := dbProduct.Price 
		tax := 1.08
		total :=  (price * float64(quantity)) * tax
		order := types.OrderInfo {
			First: ctx.QueryParam("first"),
			Last: ctx.QueryParam("last"),
			Quantity: quantity,
			Product: ctx.QueryParam("product"),
			Total: total,
		}
		
		//add order but only if it isn't already in there (checked inside of AddOrder)
		db.AddOrder(connection, dbProduct.Id, customer.Id, quantity, "No", (int64)(timestamp))
		return Render(ctx, http.StatusOK, templates.OrderPlaced(order))
	})

	e.Logger.Fatal(e.Start(":8000"))
}

// INFO: This is a simplified render method that replaces `echo`'s with a custom
// one. This should simplify rendering out of an echo route.
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
