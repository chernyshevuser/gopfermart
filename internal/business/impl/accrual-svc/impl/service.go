package impl

import (
	"sync"

	"github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"
	"github.com/chernyshevuser/gopfermart/tools/config"
	"github.com/chernyshevuser/gopfermart/tools/logger"
)

type svc struct {
	logger logger.Logger
	url    string

	in            chan accrual.Order
	outUpdated    chan accrual.Order
	outNotUpdated chan accrual.Order

	// wg is for handling with workers
	wgWorkers *sync.WaitGroup
	// wg is for handling with main updating orders goroutine
	wgProcess *sync.WaitGroup
}

func New(logger logger.Logger, in chan (accrual.Order)) (service accrual.Svc, outUpdated, outNotUpdated chan (accrual.Order)) {
	const sz = 5

	outUpdated = make(chan (accrual.Order), sz)
	outNotUpdated = make(chan (accrual.Order), sz)

	s := &svc{
		logger: logger,
		url:    config.AccrualSystemAddr,

		in:            in,
		outUpdated:    outUpdated,
		outNotUpdated: outNotUpdated,

		wgWorkers: &sync.WaitGroup{},
		wgProcess: &sync.WaitGroup{},
	}

	s.wgProcess.Add(1)
	go s.UpdateOrders()

	return s, outUpdated, outNotUpdated
}
