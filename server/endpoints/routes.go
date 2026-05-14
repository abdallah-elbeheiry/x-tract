package endpoints

import (
	"x-tract/controllers"

	"github.com/gin-gonic/gin"
)

type crudHandler interface {
	List(*gin.Context)
	Create(*gin.Context)
	GetByID(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

// Handlers collects the HTTP controllers registered by this package.
type Handlers struct {
	Admins         *controllers.AdminController
	Customers      *controllers.CustomerController
	Salesmen       *controllers.SalesmanController
	GuestEmployees *controllers.GuestEmployeeController
}

// Register attaches all HTTP endpoints to the provided Gin router.
func Register(router gin.IRoutes, handlers Handlers) {
	registerHealthRoute(router)
	registerAdminRoutes(router, handlers.Admins)
	registerCustomerRoutes(router, handlers.Customers)
	registerSalesmanRoutes(router, handlers.Salesmen)
	registerGuestEmployeeRoutes(router, handlers.GuestEmployees)
}

func registerHealthRoute(router gin.IRoutes) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func registerAdminRoutes(router gin.IRoutes, controller crudHandler) {
	if controller == nil {
		return
	}
	registerCRUDRoutes(router, "/admins", controller)
}

func registerCustomerRoutes(router gin.IRoutes, controller crudHandler) {
	if controller == nil {
		return
	}
	registerCRUDRoutes(router, "/customers", controller)
}

func registerSalesmanRoutes(router gin.IRoutes, controller crudHandler) {
	if controller == nil {
		return
	}
	registerCRUDRoutes(router, "/salesmen", controller)
}

func registerGuestEmployeeRoutes(router gin.IRoutes, controller crudHandler) {
	if controller == nil {
		return
	}
	registerCRUDRoutes(router, "/guest-employees", controller)
}

func registerCRUDRoutes(router gin.IRoutes, basePath string, handler crudHandler) {
	resourcePath := basePath + "/:id"

	router.GET(basePath, handler.List)
	router.POST(basePath, handler.Create)
	router.GET(resourcePath, handler.GetByID)
	router.PATCH(resourcePath, handler.Update)
	router.DELETE(resourcePath, handler.Delete)
}
