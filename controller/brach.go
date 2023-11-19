package controller

import (
	"database/sql"
	"encoding/json"
	"market_system/models"
	"market_system/pkg/helpers"
	"net/http"
)

func (c *Handler) Branch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.CreateBranch(w, r)
	case "GET":
		var values = r.URL.Query()
		if _, ok := values["id"]; ok {
			c.GetByIDBranch(w, r)
		} else {
			c.GetListBranch(w, r)
		}
	case "PUT":
		c.UpdateBranch(w, r)
	case "DELETE":
		c.DeleteBranch(w, r)
	}
}

func (c *Handler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	var createBranch models.CreateBranch
	err := json.NewDecoder(r.Body).Decode(&createBranch)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	resp, err := c.storage.Branch().Create(&createBranch)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusCreated, resp)
}

func (c *Handler) GetByIDBranch(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	resp, err := c.storage.Branch().GetByID(&models.BranchPrimaryKey{ID: id})
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

func (c *Handler) UpdateBranch(w http.ResponseWriter, r *http.Request) {
	var updateBranch models.UpdateBranch
	err := json.NewDecoder(r.Body).Decode(&updateBranch)
	if err != nil {
		handleResponse(w, http.StatusBadRequest, err)
		return
	}

	if !helpers.IsValidUUID(updateBranch.ID) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}
	rowsAffected, err := c.storage.Branch().Update(&updateBranch)
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	if rowsAffected == 0 {
		handleResponse(w, http.StatusBadRequest, "no rows affected")
		return
	}

	resp, err := c.storage.Branch().GetByID(&models.BranchPrimaryKey{ID: updateBranch.ID})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusAccepted, resp)
}

func (c *Handler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	if !helpers.IsValidUUID(id) {
		handleResponse(w, http.StatusBadRequest, "id is not uuid")
		return
	}

	err := c.storage.Branch().Delete(&models.BranchPrimaryKey{ID: id})
	if err != nil {
		handleResponse(w, http.StatusInternalServerError, err)
		return
	}

	handleResponse(w, http.StatusNoContent, nil)
}

func (c *Handler) GetListBranch(w http.ResponseWriter, r *http.Request) {

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

	resp, err := c.storage.Branch().GetList(&models.GetListBranchRequest{
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
