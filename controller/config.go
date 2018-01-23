package controller

import (
	"github.com/nlnwa/sigridr/log"
)

type Config struct {
	AgentAddress string
	DatabaseHost string
	DatabasePort int
	DatabaseName string
	Logger       log.Logger
}
