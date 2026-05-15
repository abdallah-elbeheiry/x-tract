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
	Auth           *controllers.AuthController
	Admins         *controllers.AdminController
	Customers      *controllers.CustomerController
	Groups         *controllers.GroupController
	Salesmen       *controllers.SalesmanController
	GuestEmployees *controllers.GuestEmployeeController
}

// Register attaches all HTTP endpoints to the provided Gin router.
func Register(router gin.IRoutes, handlers Handlers) {
	registerHealthRoute(router)
	registerAuthRoutes(router, handlers.Auth)
	registerAdminRoutes(router, handlers.Admins)
	registerCustomerRoutes(router, handlers.Customers)
	registerGroupRoutes(router, handlers.Groups)
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

func registerAuthRoutes(router gin.IRoutes, controller *controllers.AuthController) {
	if controller == nil {
		return
	}
	router.POST("/auth/login", controller.Login)
}

func registerCustomerRoutes(router gin.IRoutes, controller crudHandler) {
	if controller == nil {
		return
	}
	registerCRUDRoutes(router, "/customers", controller)
}

func registerGroupRoutes(router gin.IRoutes, controller crudHandler) {
	if controller == nil {
		return
	}
	registerCRUDRoutes(router, "/groups", controller)
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
