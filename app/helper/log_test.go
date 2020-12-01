package helper

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/inconshreveable/log15"
)

func TestLog(t *testing.T) {
	Convey("Log", t, func() {
		var buf bytes.Buffer
		l := log15.New()
		l.SetHandler(log15.StreamHandler(&buf, LogfmtFormat()))
		l.Debug("ABC|%s|%s|%s", "a", "b", "c")

		So(buf.String(), ShouldContainSubstring, "ABC|a|b|c")
	})
}
