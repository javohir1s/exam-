package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"market_system/models"
	"market_system/pkg/helpers"
)

func (c *Handler) Client(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		c.CreateClient(w, r)
	case "GET":
		var values = r.URL.Query()
		if _, ok := values["id"]; ok {
			c.GetByIDClient(w, r)
		} else {
			c.GetListClient(w, r)
		}
	case "PUT":
		c.UpdateClient(w, r)
	case "DELETE":
		c.DeleteClient(w, r)
	}
}

func (c *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {

	var createClient models.CreateClient
	err := json.NewDecoder(r.Body).Decode(&createClient)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	resp, err := c.storage.Client().Create(&createClient)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, resp)
}

func (c *Handler) GetByIDClient(w http.ResponseWriter, r *http.Request) {

	var id = r.URL.Query().Get("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	resp, err := c.storage.Client().GetByID(&models.ClientPrimaryKey{ID: id})
	if err == sql.ErrNoRows {
		handleResponse(w, http.StatusBadRequest, "no rows in result set")
		return
	}

	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (c *Handler) GetListClient(w http.ResponseWriter, r *http.Request) {

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

	search := r.URL.Query().Get("search")
	if err != nil {
		handleResponse(w, http.StatusBadRequest, "invalid query search")
		return
	}

	resp, err := c.storage.Client().GetList(&models.GetListClientRequest{
		Limit:  limit,
		Offset: offset,
		Search: search,
	})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusOK, resp)
}

func (c *Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {

	var updateClient models.UpdateClient
	err := json.NewDecoder(r.Body).Decode(&updateClient)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	if !helpers.IsValidUUID(updateClient.ID) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	rowsAffected, err := c.storage.Client().Update(&updateClient)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		handleResponse(w, http.StatusBadRequest, "no rows affected")
		return
	}

	resp, err := c.storage.Client().GetByID(&models.ClientPrimaryKey{ID: updateClient.ID})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusAccepted, resp)
}

func (c *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	err := c.storage.Client().Delete(&models.ClientPrimaryKey{ID: id})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusNoContent, nil)
}
