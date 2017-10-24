package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"text/template"
	"time"

	docopt "github.com/docopt/docopt-go"
	karma "github.com/reconquest/karma-go"
)

var (
	version = "[manual build]"
	usage   = "i3-battery-nagbar " + version + `

Shows nagbar when battery percentage is less then specified value.

Usage:
  i3-battery-nagbar [options]
  i3-battery-nagbar -h | --help
  i3-battery-nagbar --version

Options:
  --threshold <int>      Threshold to show notification.
                          [default: 15]
  --message <msg>        Show specified message in i3-nagbar. This is Go template,
                          battery percentage will be in .Percentage variable.
                          [default: Too low charge of battery: {{ .percentage}}%]
  --interval <duration>  Use specified interval as timer ticker.
                          [default: 1s]
  -h --help              Show this screen.
  --version              Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	interval, err := time.ParseDuration(args["--interval"].(string))
	if err != nil {
		log.Fatalf("unable to parse duration: %s", err)
	}

	tpl, err := template.New("message").Parse(args["--message"].(string))
	if err != nil {
		log.Fatalf("unable to parse template: %s", err)
	}

	threshold, err := strconv.Atoi(args["--threshold"].(string))
	if err != nil {
		log.Fatalf("unable to parse threshold: %s", err)
	}

	var nagbar *os.Process
	var prevPresent bool
	for range time.Tick(interval) {
		percentage, present, err := GetBatteyInfo()
		if err != nil {
			log.Println(err)
			continue
		}

		if present {
			if !prevPresent {
				stopProcess(nagbar)
			}

			prevPresent = present

			continue
		}

		prevPresent = present

		if percentage <= threshold {
			if !isRunning(nagbar) {
				nagbar, err = startNagbar(tpl, percentage)
				if err != nil {
					log.Println(err)

					continue
				}
			}
		}
	}
}

func startNagbar(
	tpl *template.Template,
	percentage int,
) (*os.Process, error) {
	buffer := bytes.NewBuffer(nil)
	err := tpl.Execute(buffer, map[string]interface{}{"percentage": percentage})
	if err != nil {
		return nil, karma.Format(
			err,
			"unable to execute template",
		)
	}

	message := buffer.String()

	args := []string{"i3-nagbar", "-m", message}

	cmd := exec.Command(args[0], args[1:]...)

	err = cmd.Start()
	if err != nil {
		return nil, karma.Describe("args", fmt.Sprintf("%q", args)).Format(
			err,
			"unable to start i3-nagbar process",
		)
	}

	return cmd.Process, nil
}

func stopProcess(process *os.Process) {
	if process == nil {
		return
	}

	defer process.Release()
	process.Signal(os.Interrupt)
}

func isRunning(process *os.Process) bool {
	if process == nil {
		return false
	}

	err := process.Signal(syscall.Signal(0))

	return err == nil
}
