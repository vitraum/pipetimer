package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/vitraum/golang-pipedrive"

	"github.com/rickar/cal"
	"github.com/vitraum/pipetimer"
)

func main() {
	var pipeName = ""
	flag.StringVar(&pipeName, "pipe", "", "name of pipeline to be used (mandatory)")

	var token = ""
	flag.StringVar(&token, "token", "", "API token to be used (mandatory)")

	var fname = "-"
	flag.StringVar(&fname, "fname", fname, "filename to write to, defaults to stdout")

	var filterID = 0
	flag.IntVar(&filterID, "filterID", 0, "filter id to use (optional)")

	var sample = 0
	flag.IntVar(&sample, "sample", sample, "number of random samples to take")

	var activity = false
	flag.BoolVar(&activity, "activity", activity, "use last activity instead of creation")

	var verbose = false
	flag.BoolVar(&verbose, "verbose", verbose, "enable verbose output")

	var seed int64
	flag.Int64Var(&seed, "seed", 0, "to be used for random sampling")

	flag.Parse()

	var ()

	if pipeName == "" {
		fmt.Println("pipe is mandatory")
		flag.Usage()
		os.Exit(1)
	}

	if token == "" {
		fmt.Println("token is mandatory")
		flag.Usage()
		os.Exit(1)
	}

	apiOptions := []pipedrive.Option{
		pipedrive.HTTPFetcher,
		pipedrive.FixedToken(token),
	}

	if verbose {
		apiOptions = append(apiOptions, pipedrive.LogURLs)
	}

	if sample > 0 {
		if seed == 0 {
			seed = time.Now().UTC().UnixNano()
		}
		if verbose {
			fmt.Printf("Using seed %d\n", seed)
		}
		rand.Seed(seed)
	}

	timerOptions := []pipetimer.Option{
		pipetimer.PipeName(pipeName),
	}

	timer, err := pipetimer.NewPipeTimer(apiOptions, timerOptions...)
	if err != nil {
		panic(err)
	}

	alldeals, err := timer.FetchDeals("", filterID)
	if err != nil {
		panic(err)
	}

	if verbose {
		fmt.Printf("fetched %d deals from pipeline %s in %d stages\n",
			len(alldeals), pipeName, len(timer.Stages))
	}

	var deals pipedrive.Deals
	if sample > 0 {
		for i := 0; i < sample; i++ {
			deals = append(deals, alldeals[rand.Intn(len(alldeals))])
		}
	} else {
		deals = alldeals
	}

	stageByID := make(map[int]*pipedrive.Stage, len(timer.Stages))

	for i, stage := range timer.Stages {
		stageByID[stage.Id] = &timer.Stages[i]
	}

	c := cal.NewCalendar()
	cal.AddGermanHolidays(c)
	c.Observed = cal.ObservedExact

	for _, deal := range deals {
		if deal.Status != "open" {
			continue
		}
		when := deal.Added.Time
		/* if !activity && deal.StageChangetime != nil {
			when = deal.StageChangetime.Time
		} else */
		if activity && deal.LastActivity != nil {
			when = deal.LastActivity.Time
		}
		diff := c.CountWorkdays(when, time.Now())
		fmt.Printf("%v;%v\n", stageByID[deal.Stage].Name, diff)
	}
}
