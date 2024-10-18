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

var connection *sql.DB

func main() {
	e := echo.New()


	// TODO: Fill in your products here with name -> price as the key -> value pair.
	dbcfg := mysql.Config{
        User:   "tmanfredo",
        Passwd: "031820",
		DBName: "tmanfredo",
	}

	connection, err := sql.Open("mysql", dbcfg.FormatDSN())
    if err != nil {
        e.Logger.Fatal(err)
    }
	
	defer connection.Close()



	productMap := map[string]struct{
		Price float64
		Image string
	}{
		"Super Mario Odyssey" : {Price: 39.99, Image: "assets/images/odyssey.png"},
		"Undertale" : {Price: 19.99, Image: "assets/images/undertale.png"},
		"Mario Kart 8 Deluxe" : {Price: 39.99, Image: "assets/images/kart.png"},
	}

	
	e.Use(etag.Etag())

	// INFO: If you wanted to load a CSS file, you'd do something like this:
	// `<link rel="stylesheet" href="assets/styles/styles.css">`
	e.Static("assets", "./assets")

	e.GET("/store", func(ctx echo.Context) error {
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(productMap)))
	})
	e.GET("/dbQueries", func(ctx echo.Context) error {
		
		customers, _ := db.GetAllCustomers(connection)
		numCustomers, _ := db.NumOfCustomers(connection)
		customer1, _ := db.GetCustomerById(connection, 1)
		customer2, _ := db.GetCustomerById(connection, 3)
		db.AddCustomer(connection, "Thomas", "Manfredo", "tmanfredo@mines.edu")
		customer2Added, _ := db.GetCustomerById(connection, 3)
		customer3, _ := db.GetCustomerByEmail(connection, "mmouse@mines.edu")
		customer4, _ := db.GetCustomerByEmail(connection, "tmanfredo@mines.edu")
		numOrdersNone, _ := db.NumOfOrders(connection)
		db.AddOrder(connection, 1, 2, 1, false)
		products, _ := db.GetAllProducts(connection)
		orders, _ := db.GetAllOrders(connection)
		numOrders, _ := db.NumOfOrders(connection)

		data := types.TemplateData{
			Products:      products,
			Customers:     customers,
			NumCustomers:  numCustomers,
			Customer1:     customer1,
			Customer2:     customer2,
			Customer2Added: customer2Added,
			Customer3:     customer3,
			Customer4:     customer4,
			Orders:        orders,
			NumOrders:     numOrders,
			NumOrdersNone: numOrdersNone,
		}
		return Render(ctx, http.StatusOK, templates.Queries(data))
	})

	// Handle the form submission and return the purchase confirmation view
	e.POST("/purchase", func(ctx echo.Context) error {
		
		
	quantity, _ := strconv.Atoi(ctx.FormValue("quantity"))
	product := ctx.FormValue("product")
	price := productMap[product].Price
	tax := 1.08
	subtotal :=  (price * float64(quantity))
	total :=  subtotal* tax
		purchase := types.PurchaseInfo{
			First:    ctx.FormValue("first"),
			Last:     ctx.FormValue("last"),
			Email:    ctx.FormValue("email"),
			Product:  product,
			Price:    price,
			Quantity: quantity,
			Donate:   ctx.FormValue("donate"),
			Tax:      tax,
			Subtotal: subtotal,
			Total:    total,
		}
		return Render(ctx, http.StatusOK, templates.Base(templates.PurchaseConfirmation(purchase)))
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
