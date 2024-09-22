package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/higansama/xyz-multi-finance/customer"
	"github.com/higansama/xyz-multi-finance/customer/dto"
	"github.com/higansama/xyz-multi-finance/internal/auth"
	"github.com/higansama/xyz-multi-finance/internal/infrastructure"
	"github.com/higansama/xyz-multi-finance/internal/response"
	"github.com/higansama/xyz-multi-finance/internal/utils"
)

type Handler struct {
	usecaseV1 domain.Usecase
	infra     infrastructure.Infrastructure
}

func NewCustomerHandler(e *gin.Engine, infra infrastructure.Infrastructure, u domain.Usecase) error {
	handler := &Handler{
		usecaseV1: u,
		infra:     infra,
	}

	public := e.Group("customer")
	public.POST("/register", handler.Register)
	public.POST("/login", handler.Login)

	customerPrivate := public.Group("transaksi")
	customerPrivate.Use(infra.Middleware.AuthMiddleware)
	customerPrivate.POST("create/", handler.CreateTransaction)
	customerPrivate.POST("mylimit/", handler.MyLimit)
	customerPrivate.POST("check/payment/:idcontract", handler.CekMyContract)
	customerPrivate.POST("create/payment", handler.CreatePayment)

	return nil
}

func (h *Handler) Register(c *gin.Context) {
	// h.usecaseV1.Create()
	var payload dto.ReqCreateCustomer
	err := c.ShouldBind(&payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	err = payload.Validate()
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	if payload.NIK == "random" {
		payload.NIK = utils.GenRandomStr(10)
	}

	err = h.usecaseV1.Create(payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "success"})
}

func (h *Handler) Login(c *gin.Context) {
	// h.usecaseV1.Create()
	var payload dto.LoginRequest
	err := c.ShouldBind(&payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	err = payload.Validate()
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	result, err := h.usecaseV1.Login(payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var payload dto.CreateTransaksiRequest
	err := c.ShouldBind(&payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	authData := auth.GetAuthData(c)
	result, err := h.usecaseV1.BuatTransaksi(authData.UserId, payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) MyLimit(c *gin.Context) {
	authData := auth.GetAuthData(c)
	result, err := h.usecaseV1.GetMyLimit(authData.UserId)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CekMyContract(c *gin.Context) {
	authdata := auth.GetAuthData(c)
	idcontract := c.Param("idcontract")

	result, err := h.usecaseV1.CekKontrak(authdata.UserId, idcontract)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreatePayment(c *gin.Context) {
	authdata := auth.GetAuthData(c)

	var payload dto.CreatePaymentRequest
	err := c.ShouldBind(&payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	err = h.usecaseV1.CreatePayment(authdata.UserId, payload)
	if err != nil {
		response.NewResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
