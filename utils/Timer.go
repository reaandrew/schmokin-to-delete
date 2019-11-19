package utils

import "time"

type StoppableTimer interface {
	Stop() time.Duration
}

type Timer interface {
	Start() StoppableTimer
}
