package atcli

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
)

func getATCLIValue(input string, startString string, index int) (int, error) {

	// Define the regular expression pattern
	regexString := fmt.Sprintf(`\b%s\s*:?\s*(-?\d+\.?\d*)`, startString)
	re, err := regexp.Compile(regexString)
	if err != nil {
		logrus.WithError(err).
			WithField("regex_string", regexString).
			Errorln("failed to compile a regular expression")
		return 0, fmt.Errorf("failed to compile regex: %w", err)
	}

	// Find the substring that matches the pattern
	match := re.FindAllStringSubmatch(input, -1)
	if len(match) <= index {
		logrus.
			WithField("input", input).
			WithField("regex_string", regexString).
			Warnln("no match was found during regex lookup")
		return 0, fmt.Errorf("no match found in %s", regexString)
	}
	intValue, err := strconv.Atoi(match[index][1])
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

// func getATCLIIntValue(input string, startString string) (int, error) {
// 	stringValue, err := getATCLIValue(input, startString)
// 	if err != nil {
// 		return 0, err
// 	}
// 	intValue, err := strconv.Atoi(stringValue)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return intValue, nil
// }
