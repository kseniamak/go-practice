package handler

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"go_hw_2/models"
)

type Handler struct {
	accounts map[string]*models.Account
	mutex    sync.RWMutex
}

func New() *Handler {
	return &Handler{
		accounts: make(map[string]*models.Account),
	}
}

func (h *Handler) CreateAccount(c echo.Context) error {
	var req models.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid request")
	}
	if req.Name == "" {
		return respondWithError(c, http.StatusBadRequest, "Empty name")
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.accounts[req.Name]; exists {
		return respondWithError(c, http.StatusForbidden, "Account already exists")
	}

	h.accounts[req.Name] = &models.Account{Name: req.Name, Amount: 0}
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) DeleteAccount(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return respondWithError(c, http.StatusBadRequest, "Empty name")
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.accounts[name]; !exists {
		return respondWithError(c, http.StatusNotFound, "Account not found")
	}

	delete(h.accounts, name)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateAmount(c echo.Context) error {
	var req models.UpdateAmountRequest
	if err := c.Bind(&req); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid request")
	}
	if req.Name == "" {
		return respondWithError(c, http.StatusBadRequest, "Empty name")
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	account, exists := h.accounts[req.Name]
	if !exists {
		return respondWithError(c, http.StatusNotFound, "Account not found")
	}

	account.Amount += req.Amount
	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateName(c echo.Context) error {
	var req models.UpdateNameRequest
	if err := c.Bind(&req); err != nil {
		return respondWithError(c, http.StatusBadRequest, "Invalid request")
	}
	if req.Name == "" || req.NewName == "" {
		return respondWithError(c, http.StatusBadRequest, "Empty name")
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	account, exists := h.accounts[req.Name]
	if !exists {
		return respondWithError(c, http.StatusNotFound, "Account not found")
	}
	if _, exists := h.accounts[req.NewName]; exists {
		return respondWithError(c, http.StatusConflict, "Account with new name already exists")
	}

	account.Name = req.NewName
	h.accounts[req.NewName] = account
	delete(h.accounts, req.Name)
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetAccount(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return respondWithError(c, http.StatusBadRequest, "Empty name")
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	account, exists := h.accounts[name]
	if !exists {
		return respondWithError(c, http.StatusNotFound, "Account not found")
	}

	return c.JSON(http.StatusOK, account)
}

func respondWithError(c echo.Context, code int, message string) error {
	return c.String(code, message)
}
