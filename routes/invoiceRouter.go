package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joy-adhikary/Restaurant-Management-back_end/controllers"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/invoices", controllers.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id", controllers.GetInvoice())
	incomingRoutes.POST("/invoices", controllers.CreateInvoice()) // create new item
	incomingRoutes.PATCH("/invoices/:invoice_id", controllers.UpdateInvoice())

}
