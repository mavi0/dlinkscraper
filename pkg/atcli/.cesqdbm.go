//unused.

package atcli

import (
	"fmt"

	"github.com/mavi0/dlinkscraper/pkg/router"
	"github.com/sirupsen/logrus"
)

type CESQdbmResult struct {
	NRBand         int
	EARFCN         int
	DLBandwidthMHz int
}

func CESQdbm(router *router.Router) (*CESQdbmResult, error) {
	if err := router.WriteCommand("atcli at+cesqdbm\n"); err != nil {
		return nil, fmt.Errorf("command failed to be sent to router: %w", err)
	}
	output, _, err := router.Expect("OK")
	if err != nil {
		return nil, fmt.Errorf("failed to get response for command from router: %w", err)
	}
	logrus.WithField("output", output).Debugln("cesqdbm obtained from router")

	result := CESQdbmResult{}

}
