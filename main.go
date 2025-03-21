package main

import (
	"net/http"

	"strconv"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/session"
	"github.com/gorilla/sessions"
	"go-store/templates"
	"go-store/types"
	"database/sql"
	"go-store/db"
	"github.com/go-sql-driver/mysql"
	"fmt"
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

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.Static("assets", "./assets")

	/*
	* DEFAULT PAGE
	*/
	e.GET("/", func(ctx echo.Context) error {
		sess, err := session.Get("session", ctx)
		if err != nil {
			return err // Return an error if session retrieval fails
		}

		security, ok := sess.Values["security"].(int)
		first_name, ok := sess.Values["first_name"].(string)
		if !ok {
			security = 0 //if session not active
			first_name = ""
		}

		if ctx.QueryParam("err") == "invalid_auth" {
			return Render(ctx, http.StatusOK, templates.InitialPage(first_name,security, "You do not have authorization to view this page"))
		} else if ctx.QueryParam("err") == "login_required" {
			return Render(ctx, http.StatusOK, templates.InitialPage(first_name, security, "You need to log in first to view this page"))
		}
		return Render(ctx, http.StatusOK, templates.InitialPage(first_name, security, ""))
	})

	/*
	* LOGIN PAGE + SESSION TOOLS
	*/
	e.GET("/login", func(ctx echo.Context) error {
		connection := connect()
		first_name, security, _ := db.GetUserSecurity(connection, ctx.FormValue("username"), ctx.FormValue("password"))
		if security == 0 {
			return Render(ctx, http.StatusOK, templates.InitialPage("",security, "Incorrect username or password, either try again or continue as guest."))
		} else {
			sess, err := session.Get("session", ctx)
			if err != nil {
				return err
			}
			
			var path string
			if security == 1 {
				path = "/order_entry"
			} else if security == 2 {
				path =  "/products"
			} else {
				path = "/store"
			}
			sess.Options = &sessions.Options{
				Path: "/",
				MaxAge:   86400 * 7,
				HttpOnly: true,
			}
			
			sess.Values["security"] = security
			sess.Values["first_name"] = first_name
			
			if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
				return err
			}
			
			return ctx.Redirect(http.StatusSeeOther, path)
		}
	})

	e.GET("/read-session", func(ctx echo.Context) error {
		sess, err := session.Get("session", ctx)
		if err != nil {
			fmt.Printf("Error getting session: %v\n", err)
			return err
		}
	
		
		if security, ok := sess.Values["security"]; ok {
			return ctx.String(http.StatusOK, fmt.Sprintf("security=%v\n", security))
		}
	
		
		return ctx.String(http.StatusOK, "No session found\n")
	})
	
	e.GET("/logout", func(ctx echo.Context) error {
		sess, _ := session.Get("session", ctx)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		sess.Values["security"] = 0
		sess.Values["first_name"] = ""
		sess.Save(ctx.Request(), ctx.Response())
		return ctx.Redirect(http.StatusSeeOther, "/")
	})

	/*
	*	STORE PAGE
	*/
	e.GET("/store", func(ctx echo.Context) error {
		connection := connect()
		storeProducts, _ := db.GetAllProducts(connection)
		sess, err := session.Get("session", ctx)
		if err != nil {
			return err // Return an error if session retrieval fails
		}

		security, ok := sess.Values["security"].(int)
		first_name, ok := sess.Values["first_name"].(string)
		if !ok {
			security = 0 //if session not active
			first_name = ""
		}
		return Render(ctx, http.StatusOK, templates.Base(first_name, security, templates.Store(storeProducts)))
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
			ProductsViewed: ctx.FormValue("productTracking"),
		}
		
		//add order but only if it isn't already in there (checked inside of AddOrder)
		db.AddOrder(connection, dbProduct.Id, customer.Id, quantity, ctx.FormValue("donate"), (int64)(timestamp))
		sess, err := session.Get("session", ctx)
		if err != nil {
			return err // Return an error if session retrieval fails
		}

		security, ok := sess.Values["security"].(int)
		first_name, ok := sess.Values["first_name"].(string)
		if !ok {
			security = 0 //if session not active
			first_name = ""
		}
		return Render(ctx, http.StatusOK, templates.Base(first_name,security,templates.PurchaseConfirmation(purchase)))
	})


	/*
	*	ADMIN PAGE
	*/
	e.GET("/admin", func(ctx echo.Context) error {
		connection := connect()
		customers, _ := db.GetAllCustomers(connection)
		orders, _ := db.GetAllOrders(connection)
		numOrders, _ := db.NumOfOrders(connection)
		products, _ := db.GetAllProducts(connection)
	
		first_name, security, redirect := ClearanceCheck(ctx, 1)
		
		if redirect != nil {
			return redirect
		} 
		return Render(ctx, http.StatusOK, templates.Admin(first_name, security, customers, orders, numOrders, products))
	})


	/*
	*	ORDER ENTRY
	*/
	e.GET("/order_entry", func(ctx echo.Context) error {
		// Connect to your database and get products
		connection := connect()
		storeProducts, _ := db.GetAllProducts(connection)
	
		first_name, security, redirect := ClearanceCheck(ctx, 1)
		if redirect != nil {
			return redirect
		} 
	
		// Render the order entry page with the security level and products
		return Render(ctx, http.StatusOK, templates.OrderEntry(first_name, security, storeProducts))
	})
	
	
	e.POST("/search_results", func(ctx echo.Context) error {
		connection := connect()
		input := ctx.QueryParam("field")
		searchTerm := ctx.QueryParam("input")
		customerSearch, _ := db.SearchCustomers(connection, input, searchTerm)
		return Render(ctx, http.StatusOK, templates.UserSearch(customerSearch))
	})
	e.POST("/order_placed", func(ctx echo.Context) error {
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


	/*
	*	PRODUCT UPDATES
	*/
	e.GET("/products", func(ctx echo.Context) error {
		connection := connect()
		storeProducts, _ := db.GetAllProducts(connection);

		first_name, security, redirect := ClearanceCheck(ctx, 2)
		if redirect != nil {
			return redirect
		} 
		return Render(ctx, http.StatusOK, templates.Products(first_name,security, storeProducts))
		
	})

	e.POST("/product_change", func(ctx echo.Context) error {
		connection := connect()
		
		
		if ctx.QueryParam("crud") == "create" { //CREATE
			//price checking
			var price float64
			if ctx.QueryParam("price") == "" {
				price = 0
			} else {
				price, _ =  strconv.ParseFloat(ctx.QueryParam("price"), 64)
			}
			var quantity int
			//quantity checking
			if ctx.QueryParam("quantity") == "" {
				quantity = 0
			} else {
				quantity, _ = strconv.Atoi(ctx.QueryParam("quantity"))
			}
			inactive, _ := strconv.Atoi(ctx.QueryParam("inactive"))
			db.CreateProduct(connection, ctx.QueryParam("name"), ctx.QueryParam("image"), quantity, price, inactive)
		} else if ctx.QueryParam("crud") == "update" { //UPDATE
			id, _ := strconv.Atoi(ctx.QueryParam("id"))
			//price checking
			var price float64
			if ctx.QueryParam("price") == "" {
				price = 0
			} else {
				price, _ =  strconv.ParseFloat(ctx.QueryParam("price"), 64)
			}
			var quantity int
			//quantity checking
			if ctx.QueryParam("quantity") == "" {
				quantity = 0
			} else {
				quantity, _ = strconv.Atoi(ctx.QueryParam("quantity"))
			}
			inactive, _ := strconv.Atoi(ctx.QueryParam("inactive"))
			db.UpdateProduct(connection, id, ctx.QueryParam("name"), ctx.QueryParam("image"), quantity, price, inactive)

		} else if ctx.QueryParam("crud") == "delete" { //DELETE
			id, _ := strconv.Atoi(ctx.QueryParam("id"))
			db.DeleteProduct(connection, id)

		} else if ctx.QueryParam("crud") == "deleteRequest" { //DELETE REQUEST CHECK
			id, _ := strconv.Atoi(ctx.QueryParam("id"))
			orders, _ := db.GetOrdersByProduct(connection, id)
			if orders {
				return ctx.String(http.StatusOK, "That product has orders!")
			} else {
				return ctx.String(http.StatusOK, "")
			}
		}

		storeProducts, _ := db.GetAllProducts(connection)
		return Render(ctx, http.StatusOK, templates.ProductTable(storeProducts))
	})

	/*
	*	VARIOUS FUNCTIONS
	*/
	e.GET("/product_quantity", func(ctx echo.Context) error {
		connection := connect()
		product, _ := db.GetProductByName(connection, ctx.QueryParam("product"))
		return ctx.String(http.StatusOK, fmt.Sprintf("%d", product.Instock))
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

func ClearanceCheck(ctx echo.Context, threshold int) (string, int, error){
	sess, err := session.Get("session", ctx)
	if err != nil {
		return "",0, err // Return an error if session retrieval fails
	}

	// Check if session contains the 'security' value
	security, ok := sess.Values["security"].(int)
	first_name, ok := sess.Values["first_name"].(string)
	if !ok {
		return "", 0, ctx.Redirect(http.StatusSeeOther, "/?err=login_required")
	} else if security < threshold {
		// Redirect to an error page if no valid security value found
		return "", 0, ctx.Redirect(http.StatusSeeOther, "/?err=invalid_auth")
	}
	return first_name,security, nil
}