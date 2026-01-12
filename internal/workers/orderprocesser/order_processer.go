package orderprocesser

import (
	"context"
	"errors"
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/logger"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/middleware"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/models"
	orderRepo "github.com/acya-skulskaya/go-musthave-diploma/internal/repository/order"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/service/accrual"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

const (
	maxConcurrentInOrdersQueue    = 10
	timeBeforeNextOrderProcessing = 50 * time.Millisecond
	maxOrdersNumByRun             = 500

	ErrMsgCouldNotGetUnprocessedOrders = "could not get unprocessed orders"
)

type OrdersProcessor struct {
	Repo    orderRepo.OrderRepository
	Accrual accrual.Service
}

func New(repo *orderRepo.OrderRepository, a *accrual.Service) *OrdersProcessor {
	return &OrdersProcessor{
		Repo:    *repo,
		Accrual: *a,
	}
}

func (op *OrdersProcessor) Run(ctx context.Context) {
	ctxRun, ctxCancelRunFn := context.WithCancel(ctx)
	defer ctxCancelRunFn()

	retryAfter := timeBeforeNextOrderProcessing

	// создаём переменную errgroup
	g := new(errgroup.Group)

	// наши данные
	orders, err := op.Repo.GetUnprocessed(ctx, maxOrdersNumByRun)
	if err != nil {
		time.AfterFunc(retryAfter, func() {
			op.Run(ctx)
		})
		return
	}

	if len(orders) == 0 {
		time.AfterFunc(retryAfter, func() {
			op.Run(ctx)
		})
		return
	}
	// если получили количество заказов равное тому котоое запрашивали, значит скорее всего есть еще необработанные заказы и можно сразу продолжить их обработку
	if len(orders) == maxOrdersNumByRun {
		retryAfter = 0
	}

	// генератор возвращает канал, через который он отправляет данные
	ordersCh := op.generator(orders)

	for data := range ordersCh {
		if ctxRun.Err() != nil {
			continue
		}

		// тут объявляем новую переменную внутри цикла, чтобы копировать переменную
		// в замыкание каждой горутины, а не использовать одно общее на всех значение.
		o := data

		// потребитель должен возвращать ошибку.
		// сигнатура анонимной функции всегда такая как в примере.
		g.Go(func() error {
			// получаем ошибку
			err = op.processOrder(ctxRun, o)
			if err != nil {
				// возвращаем ошибку
				ctxCancelRunFn()
				return err
			}

			return nil
		})
	}

	// здесь ждём выполнения горутин, и если хотя бы в одной из них возникает ошибка,
	// то присваиваем её err и обрабатываем. В этом случае просто выводим на экран.
	// Обратите внимание, что g.Wait() ждёт завершения всех запущенных горутин, даже
	// если приозошла ошибка.
	if err = g.Wait(); err != nil {
		if errors.Is(err, accrual.ErrTooManyRequests) {
			retryAfter = op.Accrual.RetryAfter
		}

		logger.Log.Error("OrdersProcessor:Run error processing orders", zap.Error(err))
	}

	time.AfterFunc(retryAfter, func() {
		op.Run(ctx)
	})
}

// generator возвращает канал, а затем отправляет в него данные
func (op *OrdersProcessor) generator(orders []models.Order) chan models.Order {
	// создаём канал данных
	ordersCh := make(chan models.Order, maxConcurrentInOrdersQueue)

	// вызываем горутину в которой отправляем данные в канал ordersCh
	go func() {
		// по завершении горутины закрываем канал
		defer close(ordersCh)

		// перебираем данные в слайсе
		for _, order := range orders {
			// отправляем данные из слайса в канал
			ordersCh <- order
		}
	}()

	// возвращаем канал с данными
	return ordersCh
}

// просто возвращает ошибку
func (op *OrdersProcessor) processOrder(ctx context.Context, order models.Order) error {
	//nolint:nilerr //no need to do anything
	if ctx.Err() != nil {
		return nil
	}

	ctx = context.WithValue(ctx, middleware.ContextKeyTraceID, "processOrder"+order.OrderNumber)

	orderResponse, err := op.Accrual.GetOrderStatus(ctx, order.OrderNumber)
	if err != nil {
		//nolint:gocritic // неудобно
		if errors.Is(err, accrual.ErrTooManyRequests) {
			return accrual.ErrTooManyRequests
		} else if errors.Is(err, accrual.ErrNoContent) {
			return nil
		} else {
			return fmt.Errorf("could not get order status from accrual service %w", err)
		}
	}

	if (orderResponse.Status == accrual.StatusRegistered && order.Status == models.OrderStatusNew) || (orderResponse.Status == accrual.StatusProcessing && order.Status == models.OrderStatusProcessing) {
		return nil
	}

	status := string(orderResponse.Status)
	if orderResponse.Status == accrual.StatusRegistered {
		status = string(models.OrderStatusNew)
	}

	found := false
	for _, v := range models.OrderStatuses {
		if string(v) == status {
			found = true
			break
		}
	}
	if !found {
		return accrual.ErrUnknownOrderStatus
	}

	if status == string(accrual.StatusInvalid) || status == string(accrual.StatusProcessed) {
		err = op.Repo.Finish(ctx, order.OrderNumber, status, orderResponse.Accrual, order.UserID)
		if err != nil {
			return fmt.Errorf("OrdersProcessor:processOrder "+"could not finish order: %w", err)
		}
	} else {
		err = op.Repo.UpdateOrderStatus(ctx, order.OrderNumber, status)
		if err != nil {
			return fmt.Errorf("OrdersProcessor:processOrder "+"could not update order: %w", err)
		}
	}

	return nil
}
