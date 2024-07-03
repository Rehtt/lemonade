//go:build !gui

package server

import (
	"github.com/Rehtt/lemonade/lemon"
	log "github.com/inconshreveable/log15"
)

func Serve(c *lemon.CLI, logger log.Logger) {
	serve(c, logger)
}
