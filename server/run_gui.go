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
	if lemon.Gui == "gui" {
		systray.Run(func() {
			systray.SetTitle("lemonade-server")
			systray.SetTooltip("Server started on :" + strconv.Itoa(c.Port))
			systray.SetIcon(lemon.Icon)
			quit := systray.AddMenuItem("Quit", "Quit")
			go func() {
				<-quit.ClickedCh
				systray.Quit()
			}()
			go func() {
				serve(c, logger)
				systray.Quit()
			}()
		}, func() {
			os.Exit(0)
		})
	} else {
		serve(c, logger)
	}
}
