package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	karma "github.com/reconquest/karma-go"
)

func GetBatteyInfo(
	ueventPath string,
) (int, bool, error) {
	file, err := os.Open(ueventPath)
	if err != nil {
		return 0, false, karma.Format(
			err,
			"unable to open uevent file: %s",
		)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var full, now float64
	var present bool
	for scanner.Scan() {
		tokens := strings.SplitN(scanner.Text(), "=", 2)
		if len(tokens) != 2 {
			continue
		}

		switch tokens[0] {
		case "POWER_SUPPLY_ENERGY_FULL_DESIGN":
			full, _ = strconv.ParseFloat(tokens[1], 64)
		case "POWER_SUPPLY_CHARGE_FULL":
			full, _ = strconv.ParseFloat(tokens[1], 64)
		case "POWER_SUPPLY_ENERGY_NOW":
			now, _ = strconv.ParseFloat(tokens[1], 64)
		case "POWER_SUPPLY_CHARGE_NOW":
			now, _ = strconv.ParseFloat(tokens[1], 64)
		case "POWER_SUPPLY_STATUS":
			present = tokens[1] != "Discharging"
		}
	}

	return int(now / full * 100), present, nil
}
