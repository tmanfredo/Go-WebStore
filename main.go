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

	err = connection.Ping()  // Check if the connection is live
	if err != nil {
		e.Logger.Fatal("Error pinging database: ", err)
	}

	productMap := map[string]struct{
		Price float64
		Image string
	}{
		"Super Mario Odyssey" : {Price: 39.99, Image: "assets/images/odyssey.png"},
		"Undertale" : {Price: 9.99, Image: "assets/images/undertale.png"},
		"Mario Kart 8 Deluxe" : {Price: 39.99, Image: "assets/images/kart.png"},
	}

	
	e.Use(etag.Etag())

	// INFO: If you wanted to load a CSS file, you'd do something like this:
	// `<link rel="stylesheet" href="assets/styles/styles.css">`
	e.Static("assets", "./assets")

	// TODO: Render your base store page here
	e.GET("/store", func(ctx echo.Context) error {
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(productMap)))
	})
	e.GET("/dbQueries", func(ctx echo.Context) error {
		products, _ := db.GetAllProducts(connection)
		customers, _ := db.GetAllCustomers(connection)
		return Render(ctx, http.StatusOK, templates.Queries(products, customers))
	})

	// TODO: Handle the form submission and return the purchase confirmation view
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
