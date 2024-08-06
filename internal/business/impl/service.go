package impl

import (
	"sync"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"
	"github.com/chernyshevuser/gopfermart/internal/business/impl/crypto-svc"
	cryptoimpl "github.com/chernyshevuser/gopfermart/internal/business/impl/crypto-svc/impl"
	sessionsvc "github.com/chernyshevuser/gopfermart/internal/business/impl/session-svc"

	accrualimpl "github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc/impl"
	sessionimpl "github.com/chernyshevuser/gopfermart/internal/business/impl/session-svc/impl"

	"github.com/chernyshevuser/gopfermart/internal/db"
	"github.com/chernyshevuser/gopfermart/tools/config"
	"github.com/chernyshevuser/gopfermart/tools/logger"
)

type gophermart struct {
	db         db.DBService
	logger     logger.Logger
	sessionSvc sessionsvc.Svc
	cryptoSvc  crypto.Svc

	accrualSvc           accrual.Svc
	accrualOutUpdated    chan accrual.Order
	accrualOutNotUpdated chan accrual.Order

	in chan accrual.Order
	// wg is for handling with writing to in chan with new api reqs
	wgIn *sync.WaitGroup

	// wg is for handling with reading from accrual system goroutines
	wgOut *sync.WaitGroup
}

func New(logger logger.Logger, db db.DBService) business.Svc {
	const sz = 10
	in := make(chan accrual.Order, sz)

	accrualSvc, outUpdated, outNotUpdated := accrualimpl.New(logger, in)

	gophermart := &gophermart{
		logger:     logger,
		db:         db,
		sessionSvc: sessionimpl.New(config.JwtSecretKey),
		cryptoSvc:  cryptoimpl.New(config.CryptoKey),

		accrualSvc:           accrualSvc,
		accrualOutNotUpdated: outNotUpdated,
		accrualOutUpdated:    outUpdated,

		in:    in,
		wgIn:  &sync.WaitGroup{},
		wgOut: &sync.WaitGroup{},
	}

	gophermart.wgOut.Add(2)
	go gophermart.handleNotUpdatedOrders()
	go gophermart.handleUpdatedOrders()

	return gophermart
}
