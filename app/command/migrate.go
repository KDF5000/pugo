package command

import (
	"github.com/kdf5000/pugo/app/migrator"
	"github.com/urfave/cli"
)

var (
	// Migrate is command of 'migrate'
	Migrate = cli.Command{
		Name:  "migrate",
		Usage: "migrate hexo to pugo",
		Flags: []cli.Flag{
			migrateSrcFlag,
			migrateDestFlag,
			debugFlag,
		},
		Before: Before,
		Action: func(ctx *cli.Context) error {
			migrator := migrator.NewMigrator(ctx.String("from"), ctx.String("to"))
			return migrator.Migrate()
		},
	}
)
