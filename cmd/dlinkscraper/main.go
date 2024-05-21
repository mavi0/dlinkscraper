package main

import (
	"fmt"
	"net/http"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mavi0/dlinkscraper/pkg/atcli"
	"github.com/mavi0/dlinkscraper/pkg/router"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dlinkscraper",
		Short: "dlinkscraper scrapes data from a D-Link CPE",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					logrus.WithError(err).Warnln("failed to read config file")
				}
			}
			if viper.GetBool("verbose") {
				logrus.SetLevel(logrus.DebugLevel)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			for {
				// curl the telnet-enable page on the router
				response, err := http.Get("http://" + viper.GetString("address") + ":8000/atsq.txt")
				if err != nil {
					logrus.WithError(err).Fatalln("error connecting to router telnet page")
				} else {
					logrus.WithField("code", response).Infoln("Tenlet enable response code")
				}

				router, err := router.NewRouter(viper.GetString("address") + ":23")
				if err != nil {
					logrus.WithError(err).Fatalln("error dialing into router via telnet")
				}

				// Login to router
				if err := router.Login(viper.GetString("username"), viper.GetString("password")); err != nil {
					logrus.WithError(err).Fatalln("failed to login to router")
				}

				// Get BNR info
				bnrInfo, err := atcli.BNRInfo(router)
				if err != nil {
					logrus.WithError(err).Fatalln("failed to get bnr info from router")
				} else {
					logrus.WithField("bnrInfo", bnrInfo).Debugln("found bnrInfo")
					router.Close()

					// Create the InfluxDB client
					client := influxdb2.NewClient(viper.GetString("influx_url"), viper.GetString("influx_token"))

					// Get the write API
					writeAPI := client.WriteAPI(viper.GetString("influx_org"), viper.GetString("influx_bucket"))

					measurements := influxdb2.NewPoint("measurement",
						map[string]string{"host": viper.GetString("hostname")},
						map[string]interface{}{
							"nr_band":                bnrInfo.NRBand,
							"earfcn":                 bnrInfo.EARFCN,
							"dl_bandwidth_mhz":       bnrInfo.DLBandwidthMHz,
							"physical_cell_id":       bnrInfo.PhysicalCellID,
							"average_pusch_power_tx": bnrInfo.AveragePUSCHPowerTx,
							"average_pucch_power_tx": bnrInfo.AveragePUCCHPowerTx,
							"rsrq":                   bnrInfo.PowerInfo.RSRQ,
							"rsrp":                   bnrInfo.PowerInfo.RSRP,
							"sinr":                   bnrInfo.PowerInfo.SINR,
							"nrcqi":                  bnrInfo.NRCQI,
							"rank":                   bnrInfo.RANK,
							"serving_beam_ssb_index": bnrInfo.ServingBeamSSBIndex,
							"fr2_serving_beam":       bnrInfo.FR2ServingBeam,
						},
						time.Now(),
					)

					// Write the main data point
					writeAPI.WritePoint(measurements)

					// Create and write points for each RXInfo
					for i, rx := range bnrInfo.RXInfo {
						rxPoint := influxdb2.NewPoint("rx_info",
							map[string]string{"host": "default", "rx_index": fmt.Sprintf("%d", i)},
							map[string]interface{}{
								"power": rx.Power,
								"ecio":  rx.ECIO,
								"rsrp":  rx.RSRP,
								"phase": rx.Phase,
								"sinr":  rx.SINR,
							},
							time.Now(),
						)
						writeAPI.WritePoint(rxPoint)
					}

					// Flush the write API to ensure the data is sent
					writeAPI.Flush()
					client.Close()

					// if interval set to over 20 secs allow it to loop forever - hard low limit of 20 secs
					if viper.GetDuration("interval") < (20 * time.Second) {
						logrus.WithField("sleep interval", viper.GetDuration("interval")).Infoln("Data written to InfluxDB. Interval less than 20 secs, exiting.")
						break
					}
					logrus.WithField("sleep interval", viper.GetDuration("interval")).Infoln("Data written to InfluxDB. Sleeping...")
					time.Sleep(viper.GetDuration("interval"))
				}
			}

		},
	}
)

func main() {
	rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "output debug logs")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().StringP("address", "a", "router_url", "telnet server address of the dlink router")
	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))

	rootCmd.PersistentFlags().StringP("username", "u", "admin", "username for the dlink router")
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))

	rootCmd.PersistentFlags().StringP("password", "p", "password", "password for the dlink router")
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	rootCmd.PersistentFlags().DurationP("interval", "i", 0, "time between each poll")
	viper.BindPFlag("interval", rootCmd.PersistentFlags().Lookup("interval"))

	rootCmd.PersistentFlags().StringP("hostname", "n", "hostname", "hostname of the system")
	viper.BindPFlag("hostname", rootCmd.PersistentFlags().Lookup("hostname"))

	rootCmd.PersistentFlags().StringP("influx_url", "s", "influxdb", "url of influxdb server")
	viper.BindPFlag("influx_url", rootCmd.PersistentFlags().Lookup("influx_url"))

	rootCmd.PersistentFlags().StringP("influx_token", "t", "token", "api token for influxdb server")
	viper.BindPFlag("influx_token", rootCmd.PersistentFlags().Lookup("influx_token"))

	rootCmd.PersistentFlags().StringP("influx_org", "o", "organisation", "organisation on influxdb server to utilise")
	viper.BindPFlag("influx_org", rootCmd.PersistentFlags().Lookup("influx_org"))

	rootCmd.PersistentFlags().StringP("influx_bucket", "b", "bucket", "bucket on influxdb server to utilise")
	viper.BindPFlag("influx_bucket", rootCmd.PersistentFlags().Lookup("influx_bucket"))

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.dlinkscraper")
	viper.AddConfigPath("/etc/dlinkscraper")

}
