package config

type Opts struct {
	SourceDir string `long:"source-dir" description:"Directory to search for source files" default:"." env:"SOURCE_DIR"`
	TargetDir string `long:"target-dir" description:"Directory to save converted files" default:"." env:"TARGET_DIR"`
}
