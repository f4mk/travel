package mail

import "errors"

var (
	ErrParseMessage   = errors.New("error mail service parsing message")
	ErrParseHeader    = errors.New("error parsing header")
	ErrAckMessage     = errors.New("error acknowledging message")
	ErrNackReqMessage = errors.New("error requeueing message")
	ErrUsecaseLayer   = errors.New("error from usecase layer")
	ErrNackMessage    = errors.New("error queueing message to dlq")
	ErrMissingHeader  = errors.New("missing header")
	ErrGetCount       = errors.New("could not get count from x-death header")
	ErrHeaderFormat   = errors.New("unexpected x-death header format")
	ErrChanFull       = errors.New("message error channel is full, might be an error")
)
