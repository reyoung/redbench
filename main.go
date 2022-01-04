package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"emperror.dev/emperror"
)

var gConcurrency = flag.Int("concurrency", 1, "number of concurrent requests")

func main() {
	flag.Parse()
	cli, err := newRedisClientFromArgs()
	emperror.Panic(err)
	defer func() {
		emperror.Panic(cli.Close())
	}()

	jobs, err := readJobFromCmdLineArgs()
	emperror.Panic(err)

	var complete sync.WaitGroup
	complete.Add(*gConcurrency)
	var tc timeCollection
	var tcMtx sync.Mutex
	for i := 0; i < *gConcurrency; i++ {
		go func(pos int) {
			defer complete.Done()
			var local timeCollection
			for i := pos; i < len(jobs); i += *gConcurrency {
				job := jobs[i]
				local.do(func() {
					cli.Do(context.Background(), job...)
				})
			}
			tcMtx.Lock()
			tc.merge(&local)
			tcMtx.Unlock()
		}(i)
	}
	complete.Wait()
	tc.prettyPrint(os.Stdout)
}
