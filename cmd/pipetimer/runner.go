package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/vitraum/golang-pipedrive"

	"github.com/vitraum/pipetimer"
)

func main() {
	var pipeName = ""
	flag.StringVar(&pipeName, "pipe", "", "name of pipeline to be used (mandatory)")

	var token = ""
	flag.StringVar(&token, "token", "", "API token to be used (mandatory)")

	var fname = "-"
	flag.StringVar(&fname, "fname", fname, "filename to write to, defaults to stdout")

	var filterName = ""
	flag.StringVar(&filterName, "filter", filterName, "filter name to use (optional)")

	var filterID = 0
	flag.IntVar(&filterID, "filterID", 0, "filter id to use (optional)")

	var sample = 0
	flag.IntVar(&sample, "sample", sample, "number of random samples to take")

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

	alldeals, err := timer.FetchDeals(filterName, filterID)
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

	var file io.WriteCloser
	if fname != "-" {
		file, err = os.Create(fname)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else {
		file = os.Stdout
	}

	updates, err := timer.FilterPipelineChanges(deals)
	if err != nil {
		panic(err)
	}

	csv := pipetimer.NewPipeWriter(file, timer.Stages)
	csv.WriteHeader()
	defer csv.Flush()
	for _, dealFlow := range updates {
		csv.Write(dealFlow)

		if !verbose {
			continue
		}
		fmt.Printf("+ %d %s\n", dealFlow.Deal.Id, dealFlow.Deal.Title)
		for _, stage := range timer.Stages {
			for _, update := range dealFlow.Updates {
				if stage.Name == update.Phase {
					when := update.PiT.Local().Format("2006-01-02 15:04")
					fmt.Printf("  %s -> %s %d", stage.Name, when, int(update.Duration/86400))
					if update.PhaseTouchdowns > 1 {
						fmt.Printf(" (%dx)", update.PhaseTouchdowns)
					}
					fmt.Println()
					break
				}
			}
		}
	}
}
