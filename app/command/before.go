package command

import (
	"os"

	"github.com/inconshreveable/log15"
	"github.com/inconshreveable/log15/ext"
	"github.com/kdf5000/pugo/app/helper"
	"github.com/urfave/cli"
)

// Before set before handler when start run cli.App
func Before(ctx *cli.Context) error {
	lv := log15.LvlInfo
	if ctx.Bool("debug") {
		lv = log15.LvlDebug
	}
	log15.Root().SetHandler(log15.LvlFilterHandler(lv, ext.FatalHandler(log15.StreamHandler(os.Stderr, helper.LogfmtFormat()))))
	return nil
}
