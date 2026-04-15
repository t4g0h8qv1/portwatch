package main

import (
	"fmt"
	"log"
	"os"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/config"
	"github.com/example/portwatch/internal/portrange"
	"github.com/example/portwatch/internal/report"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/schedule"
)

func main() {
	cfgPath := "portwatch.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ports, err := portrange.Parse(cfg.Ports)
	if err != nil {
		log.Fatalf("failed to parse port range: %v", err)
	}

	bl, err := baseline.Load(cfg.BaselineFile)
	if err != nil {
		log.Printf("no existing baseline found, creating new: %v", err)
		bl = baseline.New(cfg.Target, ports)
		if saveErr := bl.Save(cfg.BaselineFile); saveErr != nil {
			log.Fatalf("failed to save baseline: %v", saveErr)
		}
		fmt.Println("baseline created, exiting")
		return
	}

	notifier := alert.NewStdoutNotifier(os.Stdout)

	interval, err := schedule.ParseDuration(cfg.Interval)
	if err != nil {
		log.Fatalf("failed to parse interval: %v", err)
	}

	runFn := func() error {
		open, err := scanner.OpenPorts(cfg.Target, ports, cfg.Timeout)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		result := alert.Evaluate(bl, open)
		rw := report.NewWriter(os.Stdout, cfg.Format)
		rpt := report.FromAlert(cfg.Target, result, open)
		if err := rw.Write(rpt); err != nil {
			return fmt.Errorf("report write failed: %w", err)
		}

		if err := notifier.Notify(result); err != nil {
			return fmt.Errorf("notification failed: %w", err)
		}
		return nil
	}

	errFn := func(err error) {
		log.Printf("error during scan cycle: %v", err)
	}

	schedule.Run(interval, runFn, errFn)
}
