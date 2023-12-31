package accrual

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/models"
	"github.com/0loff/grade_gophermart/order"
	"go.uber.org/zap"
)

type Accrual struct {
	orderRepo       order.Repository
	accrualEndpoint string
	OrderCh         chan string
}

func NewAccrualClient(
	orderRepo order.Repository,
	endpoint string,
) {
	accrual := &Accrual{
		orderRepo:       orderRepo,
		accrualEndpoint: endpoint,
		OrderCh:         make(chan string, 10),
	}

	go accrual.Run(accrual.OrderCh)
}

func (a Accrual) Run(OrderChan chan string) {
	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case Order := <-OrderChan:
			orderUpdate := a.AccrualRequest(Order)
			a.orderRepo.UpdatePendingOrder(context.Background(), orderUpdate)

		case <-ticker.C:
			a.GetPendingOrders()
		}
	}
}

func (a Accrual) GetPendingOrders() {
	orders, err := a.orderRepo.GetPendingOrders(context.Background())
	if err != nil {
		logger.Log.Error("Unable to get list of pending orders", zap.Error(err))
	}

	var ordersList []string
	for _, order := range orders {
		ordersList = append(ordersList, order.OrderNum)
	}

	ordersChls := a.ChGenerator(ordersList)
	a.MergeChs(ordersChls)
}

func (a Accrual) ChGenerator(OrderList []string) chan string {
	inputCh := make(chan string)

	go func() {
		defer close(inputCh)

		for _, Order := range OrderList {
			inputCh <- Order
		}
	}()

	return inputCh
}

func (a Accrual) MergeChs(resultChan ...chan string) {
	var wg sync.WaitGroup

	for _, ch := range resultChan {
		chClosure := ch
		wg.Add(1)

		go func() {
			for data := range chClosure {
				a.OrderCh <- data
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
	}()
}

func (a Accrual) AccrualRequest(OrderNumber string) models.AccrualResponse {
	accrualOrder := new(models.AccrualResponse)
	endpoint := a.accrualEndpoint + "/api/orders/" + OrderNumber
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logger.Log.Error("Cannot create new request to accrual service", zap.Error(err))
	}

	request.Header.Add("Content-type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		logger.Log.Error("Accrual service request error", zap.Error(err))
		return *accrualOrder
	}
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)

	if err := dec.Decode(accrualOrder); err != nil {
		logger.Log.Error("Cannot decode response JSON body", zap.Error(err))
		return *accrualOrder
	}

	return *accrualOrder
}
