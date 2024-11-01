// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Base(content templ.Component) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html><head><title>GameHalt</title><link rel=\"stylesheet\" href=\"assets/styles/styles.css\"><link rel=\"icon\" type=\"image/x-icon\" href=\"assets/images/page_icon.ico\"></head><body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = header().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = content.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = footer().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\r\n\t\tdocument.addEventListener('DOMContentLoaded', function () {\r\n\t\t\t//listener for dropdown\r\n\t\t\tconst dropdown = document.getElementById(\"productSelection\");\r\n\t\t\tdropdown.addEventListener('change', function () {\r\n\t\t\t\t//grab the dropdown menu element\r\n\t\t\t\tconst selectedOption = dropdown.options[dropdown.selectedIndex];\r\n\t\t\t\t\r\n\t\t\t\t//change the max amount of quantity to buy so they can't buy more than are in stock\r\n\t\t\t\t//the database already disallows this but this is so the user knows and so the POST works right\r\n\t\t\t\tlet quantity = document.getElementById(\"quantity\");\r\n\t\t\t\tlet purchase = document.getElementById(\"purchase\");\r\n\r\n\t\t\t\t//get the amount of the product in stock\r\n\t\t\t\tconst available = document.getElementById(\"available\");\r\n\t\t\t\tlet instockValue = selectedOption.getAttribute('data-stock');\r\n\r\n\t\t\t\tavailable.value = Number(instockValue);\r\n\r\n\t\t\t\tif (instockValue == 0){\r\n\t\t\t\t\tquantity.disabled = true;\r\n\t\t\t\t\tpurchase.disabled = true;\r\n\t\t\t\t}\r\n\t\t\t\telse {\r\n\t\t\t\t\tquantity.disabled = false;\r\n\t\t\t\t\tquantity.setAttribute('max', instockValue.toString());\r\n\t\t\t\t\tpurchase.disabled = false;\r\n\t\t\t\t}\r\n\t\t\t\t\t;\r\n\t\t\t});\r\n\t\t});\r\n\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
