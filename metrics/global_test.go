package metrics

import (
	"testing"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	le *logrus.Entry
)


func TestMain(m *testing.M) {
	le = logrus.New().WithField("service","testing")
	os.Exit(m.Run())
}
