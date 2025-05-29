package api

import (
	"errors"
	"net/http"
	"parte3/internal/sale"
	"parte3/internal/user"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// handler holds the user service and implements HTTP handlers for user CRUD.
type handler struct {
	userService *user.Service
	saleService *sale.Service
	logger      *zap.Logger
}

// handleCreate handles POST /users
func (h *handler) handleCreate(ctx *gin.Context) {
	// request payload
	var req struct {
		Name     string `json:"name" binding:"required,regexp"` //anotations; si el content type es json, el nombre del campo es name
		Address  string `json:"address" binding:"required"`
		NickName string `json:"nickname" binding:"omitempty,regexp"` //solo letras
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &user.User{
		Name:     req.Name,
		Address:  req.Address,
		NickName: req.NickName,
	}
	if err := h.userService.Create(u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("user created", zap.Any("user", u))
	ctx.JSON(http.StatusCreated, u)
}

// handleRead handles GET /users/:id
func (h *handler) handleRead(ctx *gin.Context) {
	id := ctx.Param("id")

	u, err := h.userService.Get(id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) { //compara si el error es del tipo ErrNotFound
			// si el error es del tipo ErrNotFound, devuelve un 404
			h.logger.Warn("user not found", zap.String("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("error trying to get user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("get user succeed", zap.Any("user", u))
	ctx.JSON(http.StatusOK, u)
}

// handleUpdate handles PUT /users/:id
func (h *handler) handleUpdate(ctx *gin.Context) {
	id := ctx.Param("id")

	// bind partial update fields
	var fields *user.UpdateFields
	var user_estado user.User
	if err := ctx.ShouldBindJSON(&fields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.userService.Update(id, fields, user_estado)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			h.logger.Warn("user not found", zap.String("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("error trying to get user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("update user succeed", zap.Any("user", u))
	ctx.JSON(http.StatusOK, u)
}

// handleDelete handles DELETE /users/:id
func (h *handler) handleDelete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.userService.Delete(id); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			h.logger.Warn("user not found", zap.String("id", id))
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("error trying to get user", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("delete user succeed", zap.Any("user", id))
	ctx.Status(http.StatusNoContent)
}

func (h *handler) handleListActive(ctx *gin.Context) {
	users, err := h.userService.ListActive()
	if err != nil {
		h.logger.Error("error trying to get users", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("list users succeed", zap.Any("user", users))
	ctx.JSON(http.StatusOK, users)
}

//HANDLER PARA VENTAS

// handleCreateSale handles POST /sales
func (h *handler) handleCreateSale(ctx *gin.Context) {
	var req sale.CreateSaleRequest // Usa la request struct de tu paquete sale
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("error binding request for create sale", zap.Error(err)) // LOG AÑADIDO
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Llama al servicio de ventas
	newSale, err := h.saleService.Create(req.UserID, req.Amount)
	if err != nil {
		// Maneja errores específicos como user not found, etc.
		if errors.Is(err, sale.ErrUserNotFound) {
			h.logger.Warn("user not found for sale creation", // LOG AÑADIDO
				zap.String("user_id", req.UserID),
				zap.Float64("amount", req.Amount),
				zap.Error(err),
			)
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, sale.ErrInvalidAmount) {
			h.logger.Warn("invalid amount for sale creation", // LOG AÑADIDO
				zap.String("user_id", req.UserID),
				zap.Float64("amount", req.Amount),
				zap.Error(err),
			)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Error genérico
		h.logger.Error("error creating sale", // LOG AÑADIDO
			zap.String("user_id", req.UserID),
			zap.Float64("amount", req.Amount),
			zap.Error(err),
		)
		h.logger.Info("sale created successfully", zap.Any("sale", newSale)) // LOG AÑADIDO
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, newSale)
}

func (h *handler) handleReadSales(ctx *gin.Context) {
	id := ctx.Param("id")
	status := ctx.Param("status")

	sales, metadata, err := h.saleService.Get(id, &status)
	if err != nil {
		if errors.Is(err, sale.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	respuesta := gin.H{
		"metadata": metadata,
		"results":  sales,
	}

	ctx.JSON(http.StatusOK, respuesta)
}

func (h *handler) handleUpdateSaleStatus(ctx *gin.Context) {
	id := ctx.Param("id") //

	var req sale.UpdateSale
	if err := ctx.ShouldBindJSON(&req); err != nil { //
		h.logger.Error("error binding request for update sale status", // LOG AÑADIDO
			zap.String("sale_id", id),
			zap.Error(err),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSale, err := h.saleService.Update(id, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, sale.ErrNotFound): //
			h.logger.Warn("sale not found for status update", // LOG AÑADIDO
				zap.String("sale_id", id),
				zap.String("requested_status", req.Status),
				zap.Error(err),
			)
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, sale.ErrSaleNotActive):
			h.logger.Warn("sale not active for status update", // LOG AÑADIDO
				zap.String("sale_id", id),
				zap.String("requested_status", req.Status),
				zap.Error(err),
			)
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // O StatusConflict
		case errors.Is(err, sale.ErrSaleMustBePending):
			h.logger.Warn("sale not pending for status update", // LOG AÑADIDO
				zap.String("sale_id", id),
				zap.String("requested_status", req.Status),
				zap.Error(err),
			)
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict es apropiado aquí
		case errors.Is(err, sale.ErrInvalidSaleStateTransition):
			h.logger.Warn("invalid state transition for sale status update", // LOG AÑADIDO
				zap.String("sale_id", id),
				zap.String("requested_status", req.Status),
				zap.Error(err),
			)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			h.logger.Error("error updating sale status", // LOG AÑADIDO
				zap.String("sale_id", id),
				zap.String("requested_status", req.Status),
				zap.Error(err),
			)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //
		}
		return
	}
	h.logger.Info("sale status updated successfully", zap.Any("sale", updatedSale)) // LOG AÑADIDO
	ctx.JSON(http.StatusOK, updatedSale)                                            //
}
