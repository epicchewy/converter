package logger

type Opts struct {
	Level      string `long:"log-level" description:"Log level for the logger" default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" env:"LOG_LEVEL"`
	Output     string `long:"log-output" description:"File location to write logs" default:"stdout" env:"LOG_FILE"`
	Encoding   string `long:"log-encoding" description:"Log encoding i.e. json" default:"json" env:"LOG_ENCODING"`
	SampleRate int    `long:"log-sample-rate" description:"Logging sample rate" default:"100" env:"LOG_SAMPLE_RATE"`
}
