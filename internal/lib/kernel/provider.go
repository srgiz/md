package kernel

import (
	"md/internal/lib/kernel/auth"
	"md/internal/lib/kernel/cmdbus"
	"md/internal/lib/kernel/db"
	"md/internal/lib/kernel/logger"
)

const ServiceNameValidatorMiddleware = "bus_validator_middleware"
const ServiceNameTransactionMiddleware = "db_tx_middleware"
const ServiceNameVoterMiddleware = "auth_voter_middleware"
const ServiceNameTraceMiddleware = "logger_trace_middleware"

func NewProviderValidatorMiddleware() ServiceProvider {
	return func(di Di) {
		di.AddService(ServiceNameValidatorMiddleware, cmdbus.NewValidatorMiddleware())
	}
}

func NewProviderTransactionMiddleware(conn db.Conn) ServiceProvider {
	return func(di Di) {
		di.AddService(ServiceNameTransactionMiddleware, db.NewTxMiddleware(conn))
	}
}

func NewProviderVoterMiddleware() ServiceProvider {
	return func(di Di) {
		di.AddService(ServiceNameVoterMiddleware, auth.NewVoterMiddleware())
	}
}

func NewProviderTraceMiddleware() ServiceProvider {
	return func(di Di) {
		di.AddService(ServiceNameTraceMiddleware, logger.NewTraceMiddleware())
	}
}
