package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
)

var gJobFile = flag.String("job_file", "-", "job file")
var gN = flag.Int("n", 100, "number of replica")

type job []interface{}

func readJobFromCmdLineArgs() ([]job, error) {
	var jobFile io.Reader
	if *gJobFile == "-" {
		jobFile = os.Stdin
	} else {
		f, err := os.Open(*gJobFile)
		if err != nil {
			return nil, err
		}
		defer f.Close() // nolint:errcheck
		jobFile = f
	}

	var jobs []job
	// read file line by line
	jobScanner := bufio.NewScanner(jobFile)
	for jobScanner.Scan() {
		line := jobScanner.Text()
		if line == "" {
			continue
		}
		jobStr := strings.Split(strings.TrimSpace(line), " ")
		jobInterface := make([]interface{}, len(jobStr))
		for i, str := range jobStr {
			jobInterface[i] = str
		}
		jobs = append(jobs, jobInterface)
	}

	var result []job
	for i := 0; i < *gN; i++ {
		result = append(result, jobs...)
	}
	return result, nil
}
