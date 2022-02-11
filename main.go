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

	flag.Parse()

	if SENTRY_TOKEN == "UNDEFINED" {
		log.Fatal("Enviroment variable SENTRY_TOKEN not found. Exiting...")
	}

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

	f.WriteString("# HELP sentry_received_events_count Hourly received amount of events grouped by project-slug.\n")
	for _, p := range projects {
		stats, err := sentryProjectStats(p.Organization.Slug, p.Slug)
		if err != nil {
			log.Panicln(err)
		}
		f.WriteString(fmt.Sprintf("sentry_received_events_count{project=\"%s\"} %d\n", p.Slug, stats[0][1]))

		// issues, err := sentryProjectIssues(p.Organization.Slug, p.Slug)
		// if err != nil {
		// 	log.Panicln(err)
		// }
		// for i, issue := range issues {
		// 	fmt.Println(i, issue.Title)
		// }

	}

}
