package handler

import (
	"fmt"
	"net/http"
	"rent-car/api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h Handler) CreateCustomer(c *gin.Context) {
	customer := models.Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		handleResponse(c, "error while reading request body", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.Store.Customer().Create(customer)
	if err != nil {
		handleResponse(c, "error while creating customer, err: ", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "Succes", http.StatusOK, id)

}

func (h Handler) UpdateCustomer(c *gin.Context) {
	customer := models.Customer{}

	if err := c.ShouldBindJSON(&customer); err != nil {
		handleResponse(c, "error while updating request body", http.StatusBadRequest, err.Error())
		return
	}

	customer.Id = c.Param("id")

	err := uuid.Validate(customer.Id)
	if err != nil {
		handleResponse(c, "error while validating, err: ", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.Store.Customer().Update(customer)

	if err != nil {
		handleResponse(c, "error while updating customer, err: ", http.StatusInternalServerError, err.Error())
		return
	}
	handleResponse(c, "Success", http.StatusOK, id)
}

func (h Handler) GetAllCustomers(c *gin.Context) {
	var (
		request = models.GetAllCustomerRequest{}
	)

	request.Search = c.Query("search")
	page, err := ParsePageQueryParam(c)
	if err != nil {
		handleResponse(c, "error while parsing page, err: ", http.StatusInternalServerError, err.Error())
		return
	}

	limit, err := ParseLimitQueryParam(c)

	if err != nil {
		handleResponse(c, "error while parsing limit, err: ", http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("page: ", page)
	fmt.Println("limit: ", limit)

	request.Page = page
	request.Limit = limit
	customers, err := h.Store.Customer().GetAll(request)
	if err != nil {
		handleResponse(c, "error while getting customers, err: ", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "Success", http.StatusOK, customers)
}
func (h Handler) DeleteCustomer(c *gin.Context) {

	id := c.Param("id")
	fmt.Println("id: ", id)

	err := uuid.Validate(id)
	if err != nil {
		fmt.Println(err)
		handleResponse(c, "error while validating id, err: ", http.StatusBadRequest, err.Error())
		return
	}

	err = h.Store.Customer().Delete(id)
	if err != nil {
		fmt.Println(err)
		handleResponse(c, "error while deleting customer, err: ", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "Success", http.StatusOK, id)
}
