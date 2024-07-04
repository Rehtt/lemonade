//go:build gui

package server

import (
	"os"
	"strconv"

	"github.com/Rehtt/lemonade/lemon"
	"github.com/getlantern/systray"
	log "github.com/inconshreveable/log15"
)

func Serve(c *lemon.CLI, logger log.Logger) {
	systray.Run(func() {
		systray.SetTitle("lemonade-server")
		systray.SetTooltip("Server started on :" + strconv.Itoa(c.Port))
		systray.SetIcon(lemon.Icon)
		restart := systray.AddMenuItem("Restart", "Restart")
		quit := systray.AddMenuItem("Quit", "Quit")
		startServe := make(chan struct{}, 1)
		startServe <- struct{}{}

		go func() {
			for {
				select {
				case <-quit.ClickedCh:
					systray.Quit()
				case <-restart.ClickedCh:
					err := stopServe(logger)
					if err != nil {
						systray.SetTooltip(err.Error())
						continue
					}
					startServe <- struct{}{}
				}
			}
		}()
		go func() {
			for {
				<-startServe
				err := serve(c, logger)
				if err != nil {
					systray.SetTooltip(err.Error())
				}
			}
		}()
	}, func() {
		os.Exit(0)
	})
}
