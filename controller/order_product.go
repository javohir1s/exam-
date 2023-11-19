package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"market_system/models"
	"market_system/pkg/helpers"
)

func (c *Handler) OrderProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.CreateOrderProduct(w, r)
	case "GET":
		var values = r.URL.Query()
		if _, ok := values["order_product_id"]; ok {
			c.GetOrderProduct(w, r)
		} else {
			c.GetListOrderProduct(w, r)
		}
	case "PUT":
		c.UpdateOrderProduct(w, r)
	case "DELETE":
		c.DeleteOrderProduct(w, r)
	default:
		handleResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (c *Handler) CreateOrderProduct(w http.ResponseWriter, r *http.Request) {
	var createOrderProduct models.CreateOrderProduct
	err := json.NewDecoder(r.Body).Decode(&createOrderProduct)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := c.storage.OrderProduct().Create(&createOrderProduct)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusCreated, resp)
}

func (c *Handler) GetOrderProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("order_product_id")
	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "Invalid or missing order_product_id")
		return
	}

	resp, err := c.storage.OrderProduct().GetByID(&models.OrderProductPrimaryKey{OrderProductID: id})
	if err == sql.ErrNoRows {
		handleResponse(w, http.StatusNotFound, "Order product not found")
		return
	}

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (c *Handler) GetListOrderProduct(w http.ResponseWriter, r *http.Request) {
	limit, err := getIntegerOrDefaultValue(r.URL.Query().Get("limit"), 10)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, "invalid query limit")
		return
	}

	offset, err := getIntegerOrDefaultValue(r.URL.Query().Get("offset"), 0)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, "invalid query offset")
		return
	}

	resp, err := c.storage.OrderProduct().GetList(&models.GetListOrderProductRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (c *Handler) UpdateOrderProduct(w http.ResponseWriter, r *http.Request) {
	var updateOrderProduct models.UpdateOrderProduct
	err := json.NewDecoder(r.Body).Decode(&updateOrderProduct)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !helpers.IsValidUUID(updateOrderProduct.OrderProductID) {
		handleResponse(w, http.StatusBadRequest, "Invalid order_product_id")
		return
	}

	rowsAffected, err := c.storage.OrderProduct().Update(&updateOrderProduct)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected == 0 {
		handleResponse(w, http.StatusBadRequest, "No rows affected")
		return
	}

	resp, err := c.storage.OrderProduct().GetByID(&models.OrderProductPrimaryKey{OrderProductID: updateOrderProduct.OrderProductID})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusAccepted, resp)
}

func (c *Handler) DeleteOrderProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("order_product_id")
	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "Invalid or missing order_product_id")
		return
	}

	err := c.storage.OrderProduct().Delete(&models.OrderProductPrimaryKey{OrderProductID: id})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(w, http.StatusNoContent, nil)
}
