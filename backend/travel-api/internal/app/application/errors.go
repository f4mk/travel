package application

import "errors"

var (
	ErrInitConnDB       = errors.New("api: error initializing connection to db")
	ErrCreateKeyStore   = errors.New("api: error creating keystore")
	ErrCreateBroker     = errors.New("api: error creating message broker")
	ErrCreateQueue      = errors.New("api: error creating message queue")
	ErrCreateMailServer = errors.New("api: error creating mail server")
	ErrConnRedis        = errors.New("api: error connecting to redis")
	ErrConatructAuth    = errors.New("api: error constructing auth")
	ErrRunDebug         = errors.New("debug: error running debug server")
	ErrRunServer        = errors.New("api: error running http2 server")
	ErrStartServer      = errors.New("api: error starting http2 server")
	ErrGracefulShutdown = errors.New("api: error gracefully shutdown http2 server")
)
