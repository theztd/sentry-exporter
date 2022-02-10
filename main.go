package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	metrics_file := flag.String("o",
		getEnv("METRICS_FILE", "tmp_out.prom"),
		"Path to file where metrics should be saved.")

	if SENTRY_TOKEN == "UNDEFINED" {
		log.Fatal("Enviroment variable SENTRY_TOKEN not found. Exiting...")
	}

	flag.Parse()

	log.Printf("Running sentry exporter (ver: %s)", VERSION)

	f, err := os.Create(*metrics_file)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	projects, err := sentryProjectList()
	if err != nil {
		log.Panicln("Unable to list objects")
	}

	for _, p := range projects {
		stats, err := sentryProjectStats(p.Organization.Slug, p.Slug)
		if err != nil {
			log.Panicln(err)
		}
		f.WriteString("# HELP hourly received amount of events grouped by project-slug. Metric is counter\n")
		f.WriteString(fmt.Sprintf("sentry_received_events{\"project\"=\"%s\"} %d\n", p.Slug, stats[0][1]))

		// issues, err := sentryProjectIssues(p.Organization.Slug, p.Slug)
		// if err != nil {
		// 	log.Panicln(err)
		// }
		// for i, issue := range issues {
		// 	fmt.Println(i, issue.Title)
		// }

	}

}
