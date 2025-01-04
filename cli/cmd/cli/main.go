package main

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/epicchewy/converter/cli/pkg/executors"
	"github.com/epicchewy/converter/cli/pkg/logger"
	"github.com/epicchewy/converter/cli/pkg/parallel"
	"github.com/epicchewy/converter/cli/pkg/producers"
	"github.com/jessevdk/go-flags"
)

var (
	opts struct {
		LoggerOptions       logger.Opts                `group:"Logger Options"`
		WorkerOptions       parallel.Opts              `group:"Worker Options"`
		FileProducerOptions producers.FileProducerOpts `group:"File Producer Options"`
		QueueSize           int                        `group:"Queue Size" default:"10000" description:"The size of the queue" long:"queue-size" short:"qs"`
	}
)

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		logger.Fatalw("Parsing command line arguments", "error", err)
	}

	zap, err := logger.NewZapLogger(&opts.LoggerOptions)
	if err != nil {
		logger.Fatalw("Failed to initialize logger", "error", err)
	}
	defer zap.Sync()

	logger.Set(zap)

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		logger.Fatalw("ffmpeg is not installed. Please install ffmpeg first.")
		os.Exit(1)
	}

	start := time.Now()
	logger.Infow("Starting conversion job", "time", start.Format(time.RFC3339))

	producer, err := producers.NewFileProducer(opts.FileProducerOptions)
	if err != nil {
		logger.Fatalw("Failed to initialize file producer", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	execs := make([]parallel.Executor, opts.WorkerOptions.NumWorkers)
	for i := range execs {
		execs[i] = executors.NewFileExecutor()
	}

	queue := parallel.NewChanQueue(ctx, opts.QueueSize)
	runner := parallel.NewRunner(producer, execs, queue, nil)
	if err := runner.Run(ctx); err != nil {
		logger.Fatalw("Failed to run parallel runner", "error", err)
		os.Exit(1)
	}

	logger.Infow("Total time", "time", time.Since(start).String())

	os.Exit(0)
}

// workflow
// - validate cli inputs
// - create new producer (file)
// - each job creates a new queue
// - each job creates a new executor (file)
// - each job creates a new worker (file)
// - each job creates a new queue
