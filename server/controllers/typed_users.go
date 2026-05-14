package controllers

import "x-tract/models"

type (
	// AdminController handles admin CRUD endpoints.
	AdminController = CRUDController[models.Admin, models.NewAdmin, models.UpdateAdmin]
	// CustomerController handles customer CRUD endpoints.
	CustomerController = CRUDController[models.Customer, models.NewCustomer, models.UpdateCustomer]
	// SalesmanController handles salesman CRUD endpoints.
	SalesmanController = CRUDController[models.Salesman, models.NewSalesman, models.UpdateSalesman]
	// GuestEmployeeController handles guest employee CRUD endpoints.
	GuestEmployeeController = CRUDController[models.GuestEmployee, models.NewGuestEmployee, models.UpdateGuestEmployee]
)

func NewAdminController(store CRUDStore[models.Admin, models.NewAdmin, models.UpdateAdmin]) *AdminController {
	return NewCRUDController[models.Admin, models.NewAdmin, models.UpdateAdmin](store)
}

func NewCustomerController(store CRUDStore[models.Customer, models.NewCustomer, models.UpdateCustomer]) *CustomerController {
	return NewCRUDController[models.Customer, models.NewCustomer, models.UpdateCustomer](store)
}

func NewSalesmanController(store CRUDStore[models.Salesman, models.NewSalesman, models.UpdateSalesman]) *SalesmanController {
	return NewCRUDController[models.Salesman, models.NewSalesman, models.UpdateSalesman](store)
}

func NewGuestEmployeeController(store CRUDStore[models.GuestEmployee, models.NewGuestEmployee, models.UpdateGuestEmployee]) *GuestEmployeeController {
	return NewCRUDController[models.GuestEmployee, models.NewGuestEmployee, models.UpdateGuestEmployee](store)
}
