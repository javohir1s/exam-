package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"market_system/config"
	"market_system/storage"
)

type Handler struct {
	cfg     *config.Config
	storage storage.StorageI
}

// ErrorResponse - Json model response
type Response struct {
	Status      int         `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

func NewController(cfg *config.Config, strg storage.StorageI) *Handler {
	return &Handler{cfg: cfg, storage: strg}
}

func getIntegerOrDefaultValue(value string, defaultValue int64) (int64, error) {

	if len(value) <= 0 {
		return defaultValue, nil
	}

	number, err := strconv.Atoi(value)
	return int64(number), err
}

func handleResponse(w http.ResponseWriter, status int, data interface{}) {
	var description string
	switch code := status; {
	case code < 400:
		description = "success"
	default:
		description = "error"
		log.Println(config.Error, "error while:", Response{
			Status:      status,
			Description: description,
			Data:        data,
		})

		if code == 500 {
			data = "Internal Server Error"
		}
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{
		Status:      status,
		Description: description,
		Data:        data,
	})
}
