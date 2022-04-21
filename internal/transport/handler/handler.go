package handler

import (
	"context"
	"fmt"
	"github.com/egelis/GetBlock/internal/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Service interface {
	GetAddrMostChangedBalance(ctx context.Context) (*common.MostChanged, error)
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/most-changed-balance", h.GetAddrWithMostChangedBalance)

	return router
}

func (h *Handler) GetAddrWithMostChangedBalance(c *gin.Context) {
	start := time.Now()

	res, err := h.service.GetAddrMostChangedBalance(c.Request.Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address": res.Address,
		"value":   res.Value,
	})

	log.Info("response time: ", time.Since(start))

	log.Info(fmt.Sprintf("address: %s, absolute value: %s",
		res.Address, res.Value))
}
