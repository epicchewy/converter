package parallel

type Opts struct {
	NumWorkers int `long:"workers" description:"Number of concurrent tasks to run" default:"10" env:"NUM_WORKERS"`
}
