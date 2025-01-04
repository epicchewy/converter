package executors

import (
	"context"
	"os/exec"
	"path/filepath"
	"sync/atomic"

	"github.com/epicchewy/converter/cli/pkg/logger"
	"github.com/epicchewy/converter/cli/pkg/messages"
)

const (
	ffmpegCodec = "libmp3lame"
)

type FileExecutor struct{}

func NewFileExecutor() *FileExecutor {
	return &FileExecutor{}
}

func (e *FileExecutor) Do(ctx context.Context, item interface{}) (interface{}, error) {
	fileItem := item.(messages.FileMessage)

	srcPath, destPath := fileItem.SourceFile, fileItem.DestFile

	ext := filepath.Ext(srcPath)

	var cmd *exec.Cmd
	if ext == ".flac" {
		cmd = e.getFlacToMp3Cmd(srcPath, destPath)
	} else {
		cmd = e.getDefaultCmd(srcPath, destPath)
	}

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	logger.Infow("Converted file", "src", srcPath, "dest", destPath)

	return nil, nil
}

func (e *FileExecutor) getFlacToMp3Cmd(srcPath, destPath string) *exec.Cmd {
	return exec.Command(
		"ffmpeg",
		"-i", srcPath,
		"-codec:a", ffmpegCodec,
		"-q:a", "2",
		"-map_metadata", "0",
		"-id3v2_version", "3",
		"-r", "44100",
		"-y",
		destPath,
	)
}

func (e *FileExecutor) getDefaultCmd(srcPath, destPath string) *exec.Cmd {
	return exec.Command(
		"ffmpeg",
		"-i", srcPath,
		"-codec:a", "flac",
		"-y",
		destPath,
	)
}

type FileResultHandler struct {
	successCount atomic.Int32
	failureCount atomic.Int32
}

func (h *FileResultHandler) Handle(ctx context.Context, i interface{}, err error) {
	if err != nil {
		logger.Errorw("Failed to convert file", "error", err)
		h.failureCount.Add(1)
	} else if i != nil {
		h.successCount.Add(1)
	}
}
