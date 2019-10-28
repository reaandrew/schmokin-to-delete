package utils

import "time"

type Timer interface {
	Start()
	Stop() time.Duration
}
