package handler

import (
	"encoding/json"
	"net/http"

	"api/internal/service"
)

type CsHandler struct {
	csService *service.CsService
}

type RegisterServerHandlerResponse struct {
	Port int `json:"port"`
}

func NewCsHandler(csService *service.CsService) *CsHandler {
	return &CsHandler{
		csService: csService,
	}
}

func (h *CsHandler) RegisterServerHandler(w http.ResponseWriter, r *http.Request) {
	var server service.CsServerPayload
	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	port, err := h.csService.RegisterServer(r.Context(), server)
	if err != nil {
		http.Error(w, "Failed to create server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&RegisterServerHandlerResponse{Port: port})
}

func (h *CsHandler) GetServerList(w http.ResponseWriter, r *http.Request) {
	serverList, err := h.csService.GetServerList(r.Context())
	if err != nil {
		http.Error(w, "Failed to get server list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(serverList)
}

func (h *CsHandler) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.URL.Query().Get("ipAddress")

	if ipAddress == "" {
		http.Error(w, "Missing 'ipAddress' parameter", http.StatusBadRequest)
		return
	}

	info, err := h.csService.GetServerStatus(r.Context(), service.CsServerStatusPayload{IpAddress: ipAddress})
	if err != nil {
		http.Error(w, "Failed to get server status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(info)
}
