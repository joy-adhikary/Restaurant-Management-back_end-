package routes

import (
	controller "Restaurant-Management-back_end/controllers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/invoices", controller.GetInvoices())
	incomingRoutes.GET("/invoices/:invoice_id", controller.GetInvoice())
	incomingRoutes.POST("/invoices", controller.CreateInvoice()) // create new item
	incomingRoutes.PATCH("/invoices/:invoice_id", controller.UpdateInvoice())

}
