package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rent-car/models"

	"github.com/google/uuid"
)

func (c Controller) Customer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.CreateCustomer(w, r)
	case http.MethodGet:
		values := r.URL.Query()
		_, ok := values["id"]
		if !ok {
			c.GetAllCustomers(w, r)
		} else {
			//get one
		}
	}
}

func (c Controller) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	customer := models.Customer{}

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		errStr := fmt.Sprintf("error while decoding customer request body, err: %v\n", err)
		fmt.Println(errStr)
		handleResponse(w, http.StatusBadRequest, errStr)
		return
	}

	id, err := c.Store.Customer().Create(customer)
	if err != nil {
		fmt.Println("error while creating customer, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, id)

}

func (c Controller) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	customer := models.Customer{}

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		errStr := fmt.Sprintf("error whilde decoding update customer request body, err: %v\n", err)
		fmt.Println(errStr)
		handleResponse(w, http.StatusBadRequest, errStr)
		return
	}
	customer.Id = r.URL.Query().Get("id")
	err := uuid.Validate(customer.Id)
	if err != nil {
		fmt.Println("error while validating, err: ", err)
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.Store.Customer().Update(customer)

	if err != nil {
		fmt.Println("error while updating customer, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}
	handleResponse(w, http.StatusOK, id)
}

func (c Controller) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	var (
		values  = r.URL.Query()
		search  string
		request = models.GetAllCustomerRequest{}
	)
	if _, ok := values["search"]; ok {
		search = values["search"][0]
	}

	request.Search = search

	page, err := ParsePageQueryParam(r)
	if err != nil {
		fmt.Println("error while parsing page, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	limit, err := ParseLimitQueryParam(r)
	if err != nil {
		fmt.Println("error while parsing limit, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("page: ", page)
	fmt.Println("limit: ", limit)

	request.Page = page
	request.Limit = limit
	customers, err := c.Store.Customer().GetAll(request)
	if err != nil {
		fmt.Println("error while getting customers, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, customers)
}
func (c Controller) DeleteCustomer(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	fmt.Println("id: ", id)

	err := uuid.Validate(id)
	if err != nil {
		fmt.Println("error while validating id, err: ", err)
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Store.Customer().Delete(id)
	if err != nil {
		fmt.Println("error while deleting customer, err: ", err)
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, id)
}
