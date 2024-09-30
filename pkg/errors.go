package nex

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	NotFoundPrefix  = "not found"
	TxnFailedPrefix = "txn failed"
)

func NotFound(message string) error {
	err := fmt.Errorf("%s: %s", NotFoundPrefix, message)
	return err
}

func IsNotFound(err error) bool {
	return strings.HasPrefix(err.Error(), NotFoundPrefix)
}

func TxnFailed(message string) error {
	err := fmt.Errorf("%s: %s", TxnFailedPrefix, message)
	log.Error(err)
	return err
}

func IsTxnFailed(err error) bool {
	return strings.HasPrefix(err.Error(), TxnFailedPrefix)
}
