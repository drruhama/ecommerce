package auth

import (
	routerChi "ecommerce/infra/router/chi"
	"ecommerce/utility"
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) Handler {
	return Handler{
		svc: svc,
	}
}

// method register
func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	// proses parsing request dari client ke struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// disini kita sudah menggunakan
		// package `routerChi` yang sudah kita buat sebelumnya
		// untuk membuat sebuah response
		resp := routerChi.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "ERR BAD REQUEST",
			Error:   err.Error(),
		}
		routerChi.WriteJsonResponse(w, resp)
		return
	}

	// membuat object auth
	auth := New(req.Email, req.Password)

	// proses insert
	err = h.svc.Create(auth)
	if err != nil {
		resp := routerChi.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "ERR SERVER",
			Error:   err.Error(),
		}
		routerChi.WriteJsonResponse(w, resp)
		return
	}
	resp := routerChi.APIResponse{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
	}
	routerChi.WriteJsonResponse(w, resp)
}

// method login
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	// proses parsing request dari client ke struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp := routerChi.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "ERR BAD REQUEST",
			Error:   err.Error(),
		}
		routerChi.WriteJsonResponse(w, resp)
		return
	}

	// membuat object auth
	auth := New(req.Email, req.Password)
	// proses login, dan akan me-return object auth yang baru
	newAuth, err := h.svc.Login(auth)
	if err != nil {
		resp := routerChi.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "ERR SERVER",
			Error:   err.Error(),
		}
		routerChi.WriteJsonResponse(w, resp)
		return
	}

	// proses pembuatan token, menggunakan id dari newAuth
	token := utility.NewJWT(newAuth.Id)
	// melakukan generate token
	tokString, err := token.GenerateToken()
	if err != nil {
		resp := routerChi.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "ERR SERVER",
			Error:   err.Error(),
		}
		routerChi.WriteJsonResponse(w, resp)
		return
	}

	resp := routerChi.APIResponse{
		Status:  http.StatusOK,
		Message: "SUCCESS",
		// payloadnya kita custom, krna kita cuma ingin
		// menampilkan access tokennya saja
		Payload: map[string]interface{}{
			"token": tokString,
		},
	}
	routerChi.WriteJsonResponse(w, resp)
}
