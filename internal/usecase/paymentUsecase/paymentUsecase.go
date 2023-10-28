package paymentUsecase

import (
	"fmt"
	"simbirGo/internal/database/models"
)

type PaymentRepository interface {
	FindUserById(id uint) models.User
	SaveUser(user models.User)
}

type PaymentUsecase struct {
	r PaymentRepository
}

func New(r PaymentRepository) PaymentUsecase {
	return PaymentUsecase{r: r}
}

func (pu PaymentUsecase) IncreaseBalance(balanceId, userId uint, isAdmin bool) (int, error) {
	if isAdmin {
		user := pu.r.FindUserById(balanceId)
		if user.Id == 0 {
			return 400, fmt.Errorf("user is not exist")
		}

		user.Balance += 250000
		pu.r.SaveUser(user)

		return 200, nil
	}

	if balanceId != userId {
		return 403, fmt.Errorf("user can increase only his balance")
	}

	user := pu.r.FindUserById(balanceId)
	if user.Id == 0 {
		return 400, fmt.Errorf("user is not exist")
	}

	user.Balance += 250000
	pu.r.SaveUser(user)

	return 200, nil
}
