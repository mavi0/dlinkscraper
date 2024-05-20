package atcli

import (
	"fmt"

	"github.com/mavi0/dlinkscraper/pkg/router"
	"github.com/sirupsen/logrus"
)

type BNRInfoResult struct {
	NRBand              int
	EARFCN              int
	DLBandwidthMHz      int
	PhysicalCellID      int
	AveragePUSCHPowerTx int
	AveragePUCCHPowerTx int
	PowerInfo           struct {
		RSRQ int
		RSRP int
		SINR int
	}
	RXInfo []struct {
		Power int
		ECIO  int
		RSRP  int
		Phase int
		SINR  int
	}
	NRCQI               int
	RANK                int
	ServingBeamSSBIndex int
	FR2ServingBeam      int
}

func BNRInfo(router *router.Router) (*BNRInfoResult, error) {
	if err := router.WriteCommand("atcli at+bnrinfo\n"); err != nil {
		return nil, fmt.Errorf("command failed to be sent to router: %w", err)
	}
	output, _, err := router.Expect("OK")
	if err != nil {
		return nil, fmt.Errorf("failed to get response for command from router: %w", err)
	}
	logrus.WithField("output", output).Debugln("bnr info obtained from router")

	result := BNRInfoResult{}

	// nrband
	nrBand, err := getATCLIValue(output, "NR BAND", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get nr band via atcli commands")
	} else {
		logrus.WithField("nr_band", nrBand).Infoln("found nr band info")
	}
	result.NRBand = nrBand

	// earfcn
	earfcn, err := getATCLIValue(output, "EARFCN", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get earfcn via atcli commands")
	} else {
		logrus.WithField("earfcn", earfcn).Infoln("found earfcn info")
	}
	result.EARFCN = earfcn

	// DL bandwidth
	dlBandwidth, err := getATCLIValue(output, "DL_bandwidth", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get dlBandwidth via atcli commands")
	} else {
		logrus.WithField("dlBandwidth", dlBandwidth).Infoln("found dlBandwidth info")
	}

	// physical cell ID
	phy, err := getATCLIValue(output, "physical cell ID", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get physical cell ID via atcli commands")
	} else {
		logrus.WithField("physical cell ID", phy).Infoln("found physical cell ID info")
	}
	result.PhysicalCellID = phy

	// averaged PUSCH TX power
	PUSCH_TX, err := getATCLIValue(output, "averaged PUSCH TX power ", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get averaged PUSCH TX power  via atcli commands")
	} else {
		logrus.WithField("averaged PUSCH TX power", PUSCH_TX).Infoln("found averaged PUSCH TX power info")
	}
	result.AveragePUSCHPowerTx = PUSCH_TX

	// NR CQI
	NRCQI, err := getATCLIValue(output, "NR CQI", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get averaged NR CQI via atcli commands")
	} else {
		logrus.WithField("NRCQI", NRCQI).Infoln("found NRCQI info")
	}
	result.NRCQI = NRCQI

	// RANK
	RANK, err := getATCLIValue(output, "RANK", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get averaged RANK via atcli commands")
	} else {
		logrus.WithField("RANK", RANK).Infoln("found RANK info")
	}
	result.RANK = RANK

	// ServingBeamSSBIndex
	ServingBeamSSBIndex, err := getATCLIValue(output, "Serving Beam SSB index ", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get averaged Serving beam SSB index via atcli commands")
	} else {
		logrus.WithField("Serving Beam SSB index ", ServingBeamSSBIndex).Infoln("found Serving Beam SSB index info")
	}
	result.ServingBeamSSBIndex = ServingBeamSSBIndex

	// FR2 serving Beam
	FR2ServingBeam, err := getATCLIValue(output, "FR2 serving Beam", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get FR2 serving Beam via atcli commands")
	} else {
		logrus.WithField("FR2 serving Beam", FR2ServingBeam).Infoln("found FR2 serving Beam info")
	}
	result.FR2ServingBeam = FR2ServingBeam

	// averaged PUSCH TX power
	PUCCH_TX, err := getATCLIValue(output, "averaged PUCCH TX power ", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get averaged PUCCH TX power  via atcli commands")
	} else {
		logrus.WithField("averaged PUSCH TX power", PUCCH_TX).Infoln("found averaged PUCCH TX power info")
	}
	result.AveragePUCCHPowerTx = PUCCH_TX

	// RSRQ
	RSRQ, err := getATCLIValue(output, "RSRQ", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get RSRQ via atcli commands")
	} else {
		logrus.WithField("RSRQ", RSRQ).Infoln("found RSRQ info")
	}
	result.PowerInfo.RSRQ = RSRQ

	// RSRP
	RSRP, err := getATCLIValue(output, "RSRP", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get RSRP via atcli commands")
	} else {
		logrus.WithField("RSRP", RSRP).Infoln("found RSRP info")
	}
	result.PowerInfo.RSRP = RSRP

	// SINR
	SINR, err := getATCLIValue(output, "SINR", 0)
	if err != nil {
		logrus.WithError(err).Errorln("failed to get SINR via atcli commands")
	} else {
		logrus.WithField("SINR", SINR).Infoln("found SINR info")
	}
	result.PowerInfo.SINR = SINR

	result.RXInfo = make([]struct {
		Power int
		ECIO  int
		RSRP  int
		Phase int
		SINR  int
	}, 4)
	for i := 0; i < 4; i++ {
		// ecio
		ecio, err := getATCLIValue(output, "ecio", i)
		if err != nil {
			logrus.WithError(err).Errorln("failed to get ecio via atcli commands")
		} else {
			logrus.WithField("ecio", ecio).WithField("rx_info", i).Infoln("found ecio info")
			result.RXInfo[i].ECIO = ecio
		}
		// rsrp
		rsrp, err := getATCLIValue(output, "rsrp", i)
		if err != nil {
			logrus.WithError(err).Errorln("failed to get rsrp via atcli commands")
		} else {
			logrus.WithField("rsrp", rsrp).WithField("rx_info", i).Infoln("found rsrp info")
			result.RXInfo[i].RSRP = rsrp
		}
		// power
		power, err := getATCLIValue(output, "power", i+2)
		if err != nil {
			logrus.WithError(err).Errorln("failed to get power via atcli commands")
		} else {
			logrus.WithField("power", power).WithField("rx_info", i).Infoln("found power info")
			result.RXInfo[i].Power = power
		}
		// phase
		phase, err := getATCLIValue(output, "phase", i)
		if err != nil {
			logrus.WithError(err).Errorln("failed to get phase via atcli commands")
		} else {
			logrus.WithField("phase", phase).WithField("rx_info", i).Infoln("found phase info")
			result.RXInfo[i].Phase = phase
		}
		// sinr
		sinr, err := getATCLIValue(output, "sinr", i)
		if err != nil {
			logrus.WithError(err).Errorln("failed to get sinr via atcli commands")
		} else {
			logrus.WithField("sinr", sinr).WithField("rx_info", i).Infoln("found sinr info")
			result.RXInfo[i].SINR = sinr
		}
	}

	return &result, nil
}
