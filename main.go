package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	etag "github.com/pablor21/echo-etag/v4"

	"go-store/templates"
	//"go-store/types"
)

func main() {
	// TODO: Fill in your products here with name -> price as the key -> value pair.
	products := map[string]struct{
		Price float64
		Image string
	}{
		"Super Mario Odyssey" : {Price: 39.99, Image: "assets/images/odyssey.png"},
		"Undertale" : {Price: 9.99, Image: "assets/images/undertale.png"},
		"Mario Kart 8 Deluxe" : {Price: 39.99, Image: "assets/images/kart.png"},
	}

	/* images :=map[string]string {
		"Super Mario Odyssey" : "assets/images/odyssey.png",
		"Undertale" : "assets/images/undertale.png",
		"Mario Kart 8 Deluxe" : "assets/images/kart.png",
	} */

	e := echo.New()
	e.Use(etag.Etag())

	// INFO: If you wanted to load a CSS file, you'd do something like this:
	// `<link rel="stylesheet" href="assets/styles/styles.css">`
	e.Static("assets", "./assets")

	// TODO: Render your base store page here
	e.GET("/store", func(ctx echo.Context) error {
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(products)))
	})

	// TODO: Handle the form submission and return the purchase confirmation view
	e.POST("/purchase", func(ctx echo.Context) error {
		// TODO: Grab the form details from ctx.FormValue("...")
		// purchaseInfo := types.PurchaseInfo{
		// 	// TODO: Maybe use this structure to pass the data to your purchase confirmation page
		// 	// ...
		// }
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(products)))
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
