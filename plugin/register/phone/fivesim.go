package phone

import (
	"github.com/Fef0/go-5sim/fivesim"
	"golang.org/x/xerrors"
	"time"
)

type FiveSim struct {
	cl         *fivesim.Client
	orderId    int
	fetchCount int
}

func NewFiveSim(token string) *FiveSim {
	return &FiveSim{cl: fivesim.NewClient(token)}
}

func (i *FiveSim) GetNumber() (string, error) {
	number, err := i.cl.BuyActivationNumber("canada", "virtual12", "line", "")
	if err != nil {
		return "", err
	}
	i.orderId = number.ID
	return number.Phone, nil
}

func (i *FiveSim) BanNumber() {
	_, _ = i.cl.BanOrder(i.orderId)
}
func (i *FiveSim) CancelNumber() {
	_, _ = i.cl.CancelOrder(i.orderId)
}

func (i *FiveSim) FinishOrder() {
	_, _ = i.cl.FinishOrder(i.orderId)
}

func (i *FiveSim) GetCode() (string, error) {
	order, err := i.cl.CheckOrder(i.orderId)
	if err != nil {
		return "", err
	}
	if len(order.SMS) > 0 {
		i.FinishOrder()
		return order.SMS[0].Code, nil
	}
	i.fetchCount++
	if i.fetchCount > 10 {
		i.CancelNumber()
		return "", xerrors.New("time out")
	}
	time.Sleep(time.Second * 10)
	return i.GetCode()
}
