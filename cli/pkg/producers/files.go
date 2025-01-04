package producers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/epicchewy/converter/cli/pkg/logger"
	"github.com/epicchewy/converter/cli/pkg/messages"
	"github.com/epicchewy/converter/cli/pkg/parallel"
)

var (
	ValidSourceExts = map[string]bool{
		".m4a":  true,
		".wav":  true,
		".flac": true,
		".aac":  true,
		".ogg":  true,
		".opus": true,
		".m4r":  true,
	}

	ValidTargetExts = map[string]bool{
		".mp3": true,
	}
)

func IsValidSourceExt(ext string) bool {
	return ValidSourceExts[strings.ToLower(ext)]
}

func IsValidTargetExt(ext string) bool {
	return ValidTargetExts[strings.ToLower(ext)]
}

type FileProducerOpts struct {
	SourceFile     string   `short:"sf" long:"source-file" description:"Source file to convert" env:"SOURCE_FILE"`
	SourceDir      string   `short:"sd" long:"source-dir" description:"Source directory to convert" env:"SOURCE_DIR"`
	TargetDir      string   `short:"td" long:"target-dir" description:"Target directory to save converted files" env:"TARGET_DIR"`
	SourceFileExts []string `short:"sfe" long:"source-file-exts" description:"Source file extensions to convert" env:"SOURCE_FILE_EXTS"`
	TargetFileExts []string `short:"tfe" long:"target-file-exts" description:"Target file extensions to save" env:"TARGET_FILE_EXTS"`
}

func (opts *FileProducerOpts) Validate() error {
	if opts.SourceFile != "" && (opts.SourceDir != "" || opts.SourceFileExts != nil) {
		return errors.New("source file and source dir/exts cannot be used together")
	}

	if opts.SourceFileExts == nil || len(opts.SourceFileExts) == 0 {
		return errors.New("source file extensions are required")
	}

	if opts.TargetFileExts == nil || len(opts.TargetFileExts) == 0 {
		return errors.New("target file extensions are required")
	}

	if opts.SourceDir == "" {
		return errors.New("source dir is required")
	}

	if opts.TargetDir == "" {
		return errors.New("target dir is required")
	}

	for _, ext := range opts.SourceFileExts {
		if !IsValidSourceExt(ext) {
			return fmt.Errorf("invalid source file extension: %s", ext)
		}
	}

	for _, ext := range opts.TargetFileExts {
		if !IsValidTargetExt(ext) {
			return fmt.Errorf("invalid target file extension: %s", ext)
		}
	}

	return nil
}

type FileProducer struct {
	SourceFile     string
	SourceDir      string
	TargetDir      string
	SourceFileExts map[string]bool
	TargetFileExts map[string]bool
}

func NewFileProducer(opts FileProducerOpts) (*FileProducer, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	sourceFileExts := make(map[string]bool)
	for _, ext := range opts.SourceFileExts {
		sourceFileExts[ext] = true
	}

	targetFileExts := make(map[string]bool)
	for _, ext := range opts.TargetFileExts {
		targetFileExts[ext] = true
	}

	return &FileProducer{
		SourceFile:     opts.SourceFile,
		SourceDir:      opts.SourceDir,
		TargetDir:      opts.TargetDir,
		SourceFileExts: sourceFileExts,
		TargetFileExts: targetFileExts,
	}, nil
}

func (p *FileProducer) Produce(ctx context.Context, workQueue parallel.Queue) error {
	logger.Infow("Producing files", "source_dir", p.SourceDir, "target_dir", p.TargetDir, "source_exts", p.SourceFileExts, "target_exts", p.TargetFileExts)

	absSrcDir, err := filepath.Abs(p.SourceDir)
	if err != nil {
		return err
	}

	absTargetDir, err := filepath.Abs(p.TargetDir)
	if err != nil {
		return err
	}

	// send message for source file and return
	if p.SourceFile != "" {
		logger.Infow("Producing file", "path", p.SourceFile)
		for ext := range p.SourceFileExts {
			destFileName := fmt.Sprintf("%s.%s", p.SourceFile, ext)
			workQueue.Enqueue(ctx, messages.FileMessage{
				SourceFile: filepath.Join(absSrcDir, p.SourceFile),
				DestFile:   filepath.Join(absTargetDir, destFileName),
			})
		}
		return nil
	}

	// send message for source dir
	err = filepath.Walk(absSrcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filename := filepath.Base(path)

		if !info.IsDir() && p.SourceFileExts[strings.ToLower(filepath.Ext(filename))] {
			logger.Infow("Producing file", "path", path)
			for ext := range p.TargetFileExts {
				destFilePath := filepath.Join(absTargetDir, fmt.Sprintf("%s.%s", filename, ext))
				workQueue.Enqueue(ctx, messages.FileMessage{
					SourceFile: path,
					DestFile:   destFilePath,
				})
			}
		}

		return nil
	})

	if err != nil {
		logger.Errorw("Failed to walk source dir", "error", err)
		return err
	}

	logger.Infow("Produced files", "num_files", workQueue.Size())

	return nil
}
