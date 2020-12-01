package builder

import (
	"fmt"
	"net/url"
	"path"

	"github.com/unknwon/com"
	"github.com/kdf5000/pugo/app/theme"
	"github.com/inconshreveable/log15"
)

// ReadTheme read *Theme to *Context
func ReadTheme(ctx *Context) {
	if ctx.Source == nil {
		ctx.Err = fmt.Errorf("theme depends on loaded source data")
		return
	}
	dir, _ := toDir(ctx.ThemeName)
	if !com.IsDir(dir) {
		ctx.Err = fmt.Errorf("theme directory '%s' is missing", dir)
		return
	}
	log15.Info("Theme|%s", dir)
	ctx.Theme = theme.New(dir)
	ctx.Theme.Func("url", func(str ...string) string {
		if len(str) > 0 {
			if ur, _ := url.Parse(str[0]); ur != nil {
				if ur.Host != "" {
					return str[0]
				}
			}
		}
		return path.Join(append([]string{ctx.Source.Meta.Path}, str...)...)
	})
	ctx.Theme.Func("fullUrl", func(str ...string) string {
		return ctx.Source.Meta.Root + path.Join(str...)
	})
	if err := ctx.Theme.Validate(); err != nil {
		log15.Warn("Theme|%s|%s", dir, err.Error())
	}
}
