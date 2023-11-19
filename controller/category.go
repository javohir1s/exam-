package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"market_system/models"
	"market_system/pkg/helpers"
)

func (c *Handler) Category(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		c.CreateCategory(w, r)
	case "GET":
		var values = r.URL.Query()
		if _, ok := values["id"]; ok {
			c.GetByIDCategory(w, r)
		} else {
			c.GetListCategory(w, r)
		}
	case "PUT":
		c.UpdateCategory(w, r)
	case "DELETE":
		c.DeleteCategory(w, r)
	}
}

func (c *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {

	var createCategory models.CreateCategory
	err := json.NewDecoder(r.Body).Decode(&createCategory)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	if createCategory.ParentID != "" {
		if !helpers.IsValidUUID(createCategory.ParentID) {
			handleResponse(w, http.StatusBadRequest, "parent id is not uuid")
			return
		}
	}

	resp, err := c.storage.Category().Create(&createCategory)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, resp)
}

func (c *Handler) GetByIDCategory(w http.ResponseWriter, r *http.Request) {

	var id = r.URL.Query().Get("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	resp, err := c.storage.Category().GetByID(&models.CategoryPrimaryKey{Id: id})
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

func (c *Handler) GetListCategory(w http.ResponseWriter, r *http.Request) {

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

	resp, err := c.storage.Category().GetList(&models.GetListCategoryRequest{
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

func (c *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {

	var updateCategory models.UpdateCategory
	err := json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	if !helpers.IsValidUUID(updateCategory.Id) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	if updateCategory.ParentID != "" {
		if !helpers.IsValidUUID(updateCategory.ParentID) {
			handleResponse(w, http.StatusBadRequest, "parent id is not uuid")
			return
		}
	}

	rowsAffected, err := c.storage.Category().Update(&updateCategory)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		handleResponse(w, http.StatusBadRequest, "no rows affected")
		return
	}

	resp, err := c.storage.Category().GetByID(&models.CategoryPrimaryKey{Id: updateCategory.Id})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusAccepted, resp)
}

func (c *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")

	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	err := c.storage.Category().Delete(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusNoContent, nil)
}
