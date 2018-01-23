package agent

import (
	"github.com/nlnwa/sigridr/log"
)

type Config struct {
	WorkerAddress string
	DatabaseHost  string
	DatabasePort  int
	DatabaseName  string
	Logger        log.Logger
}
