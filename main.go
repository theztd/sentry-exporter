package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var metrics_file, read_token, port string
var interval int
var daemon_mode bool

func main() {
	_daemon_mode, _ := strconv.ParseBool((getEnv("DAEMON_MODE", "false")))

	flag.StringVar(&metrics_file, "o",
		getEnv("METRICS_FILE", "tmp_out.prom"),
		"Path to file where metrics should be saved")

	flag.StringVar(&read_token, "read-token",
		getEnv("READ_TOKEN", ""),
		"Read token for HTTP endpoint, if empty, no autentication is required")

	flag.IntVar(&interval, "i",
		60,
		"Interval in seconds between gathering data from sentry")

	flag.BoolVar(&daemon_mode, "daemon-mode",
		_daemon_mode,
		"Run in daemon mode and server results via HTTP protocol")
	flag.StringVar(&port, "p",
		getEnv("PORT", ":9091"),
		"Listening port")

	flag.Parse()

	if SENTRY_TOKEN == "UNDEFINED" {
		log.Fatal("Enviroment variable SENTRY_TOKEN not found. Exiting...")
	}

	log.Printf("Running sentry exporter (ver: %s)", VERSION)

	if daemon_mode == true {
		log.Println("Running in daemon mode.")

		// poustim jako go rutinu
		go serveMetrics(metrics_file, port, read_token)

		for true {
			log.Println("Gathering metrics...")
			updateMetrics(metrics_file)
			time.Sleep(time.Duration(interval) * time.Second)
		}

	} else {
		updateMetrics(metrics_file)
	}

}

func updateMetrics(metrics_file string) {
	/*
		Gather metrics and save them to file

	*/
	f, err := os.Create(metrics_file)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	projects, err := sentryProjectList()
	if err != nil {
		log.Panicln("Unable to list objects")
	}

	f.WriteString("# HELP sentry_received_events_count Hourly received amount of events grouped by project-slug.\n")
	for _, p := range projects {
		stats, err := sentryProjectStats(p.Organization.Slug, p.Slug)
		if err != nil {
			log.Panicln(err)
		}
		f.WriteString(fmt.Sprintf("sentry_received_events_count{project=\"%s\"} %d\n", p.Slug, stats[0][1]))

	}

}
