package paymentHandler

import (
	"net/http"
	httpUtil "simbirGo/internal/httputil"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentUsecase interface {
	IncreaseBalance(balanceId, userId uint, isAdmin bool) (int, error)
}

type PaymentHandler struct {
	pu PaymentUsecase
}

func New(pu PaymentUsecase) PaymentHandler {
	return PaymentHandler{pu: pu}
}

// @Summary Пополнение баланса
// @Tags 3. PaymentController
// @Description Добавляет на баланс пользователся с id = {id} 250 000. Администраторы могут изменять баланс любому пользователю, обычные пользователи только себе
// @Security ApiKeyAuth
// @Param id path uint true "Account id"
// @Success 200
// @Failure 400 {object} httpUtil.ResponseError
// @Failure 401 {object} httpUtil.ResponseError
// @Failure 403 {object} httpUtil.ResponseError
// @Router /api/Payment/Hesoyam/{id} [post]
func (ph PaymentHandler) IncreaseBalance(ctx *gin.Context) {
	balanceIdStr := ctx.Param("id")
	balanceId, err := strconv.Atoi(balanceIdStr)
	if err != nil || balanceId < 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": "invalid value of id param"})
		return
	}
	userId := ctx.GetUint("id")
	isAdmin := ctx.GetBool("isAdmin")
	code, err := ph.pu.IncreaseBalance(uint(balanceId), userId, isAdmin)
	if err != nil {
		httpUtil.NewResponseError(ctx, code, err)
		return
	}

	ctx.Status(code)
}
