package http

import (
	"encoding/json"
	"errors"
	"hitalent-test-task/internal/services"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *services.DepartmentService
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /_info", h.getInfo)
	mux.HandleFunc("POST /departments", h.createDepartment)
	mux.HandleFunc("POST /departments/{id}/employees/", h.addEmployee)
	mux.HandleFunc("GET /departments/{id}", h.getDepartment)
	mux.HandleFunc("PATCH /departments/{id}", h.moveDepartment)
	mux.HandleFunc("DELETE /departments/{id}", h.removeDepartment)
}

func (h *Handler) getInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}

func (h *Handler) createDepartment(w http.ResponseWriter, r *http.Request) {

	var req CreateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "name is required")
	}

	department, err := h.service.CreateDepartment(r.Context(), req.Name, req.ParentID)

	log.Println(err)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDepartmentNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "department not found")
		case errors.Is(err, services.ErrDepartmentAlreadyExists):
			writeError(w, http.StatusConflict, "CONFLICT", "department exists")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusCreated, department)
}

func (h *Handler) addEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid department id")
		return
	}

	var req CreateEmployeeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	if req.FullName == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "full_name is required")
	}
	if req.Position == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "position is required")
	}

	employee, err := h.service.CreateEmployee(r.Context(), id, req.FullName, req.Position, req.HiredAt)

	log.Println(err)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDepartmentNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "department not found")
		case errors.Is(err, services.ErrInvalidEmployeeName):
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid employee name")
		case errors.Is(err, services.ErrInvalidPosition):
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid position")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusCreated, ToEmployeeResponse(employee))
}

func (h *Handler) getDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid department id")
		return
	}

	depth := 1
	if value := r.URL.Query().Get("depth"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid depth")
			return
		}
		depth = parsed
	}

	includeEmployees := true
	if value := r.URL.Query().Get("include_employees"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid include_employees")
			return
		}
		includeEmployees = parsed
	}

	tree, err := h.service.GetDepartmentTree(
		r.Context(),
		id,
		depth,
		includeEmployees,
	)
	log.Println(err)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDepartmentNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "department not found")
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, ToDepartmentTreeResponse(tree))
}

func (h *Handler) moveDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid department id")
		return
	}

	var req UpdateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	department, err := h.service.UpdateDepartment(
		r.Context(),
		id,
		req.Name,
		req.ParentID,
	)

	log.Println(err)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDepartmentNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "department not found")

		case errors.Is(err, services.ErrDepartmentCycle):
			writeError(w, http.StatusConflict, "CONFLICT", "cycle detected")

		case errors.Is(err, services.ErrDepartmentAlreadyExists):
			writeError(w, http.StatusConflict, "DEPARTMENT_ALREADY_EXISTS", "department already exists")

		case errors.Is(err, services.ErrSelfParent):
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid parent")

		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}

		return
	}
	writeJSON(w, http.StatusOK, ToDepartmentResponse(department))
}

func (h *Handler) removeDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)

	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid department id")
		return
	}

	mode := r.URL.Query().Get("mode")

	switch mode {

	case "cascade":
		err = h.service.DeleteCascade(r.Context(), id)

	case "reassign":
		raw := r.URL.Query().Get("reassign_to_department_id")
		if raw == "" {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "reassign_to_department_id is required")
			return
		}
		targetID, errParse := strconv.ParseInt(r.URL.Query().Get("reassign_to_department_id"), 10, 64)

		if errParse != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid reassign_to_department_id")
			return
		}

		err = h.service.DeleteReassign(r.Context(), id, targetID)

	default:
		writeError(w, http.StatusBadRequest, "INTERNAL_ERROR", "mode must be cascade or reassign")
		return
	}

	log.Println(err)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDepartmentNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "department not found")

		case errors.Is(err, services.ErrParentDepartmentNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "target department not found")

		case errors.Is(err, services.ErrDepartmentCycle):
			writeError(w, http.StatusConflict, "CONFLICT", "department cycle")

		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func New(service *services.DepartmentService) *Handler {
	return &Handler{
		service: service,
	}
}
