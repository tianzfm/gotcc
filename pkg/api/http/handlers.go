package http

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gotcc/internal/engine"
    "gotcc/internal/model"
    "gotcc/internal/dao"
)

// Handler struct to hold dependencies
type Handler struct {
    Engine *engine.TransactionEngine
    TaskDAO dao.TaskDAO
}

// NewHandler creates a new instance of Handler
func NewHandler(engine *engine.TransactionEngine, taskDAO dao.TaskDAO) *Handler {
    return &Handler{
        Engine: engine,
        TaskDAO: taskDAO,
    }
}

// StartTransaction handles the initiation of a new transaction
func (h *Handler) StartTransaction(c *gin.Context) {
    var flow model.Flow
    if err := c.ShouldBindJSON(&flow); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Start the transaction using the engine
    instance, err := h.Engine.StartTransaction(flow)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, instance)
}

// GetTransactionStatus handles fetching the status of a transaction
func (h *Handler) GetTransactionStatus(c *gin.Context) {
    id := c.Param("id")

    status, err := h.Engine.GetTransactionStatus(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, status)
}

// CancelTransaction handles the cancellation of a transaction
func (h *Handler) CancelTransaction(c *gin.Context) {
    id := c.Param("id")

    err := h.Engine.CancelTransaction(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Transaction cancelled successfully"})
}