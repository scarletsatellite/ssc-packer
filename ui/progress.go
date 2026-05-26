package ui

import (
	"fmt"
	"io"
	"strings"
)

type ProgressWriter struct {
	W           io.Writer
	Total       int64
	Written     int64
	FileName    string
	lastPercent int
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.W.Write(p)
	if n > 0 {
		pw.Written += int64(n)
		pw.draw()
	}
	return n, err
}

func (pw *ProgressWriter) draw() {
	if pw.Total <= 0 {
		return
	}

	percent := int((pw.Written * 100) / pw.Total)
	if percent > 100 {
		percent = 100
	}

	if percent == pw.lastPercent && percent < 100 {
		return
	}
	pw.lastPercent = percent

	barLength := 30
	filledLength := int((percent * barLength) / 100)

	var bar strings.Builder
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			bar.WriteString("█")
		} else {
			bar.WriteString("-")
		}
	}

	task := pw.FileName
	if len(task) > 20 {
		task = task[:17] + "..."
	}

	fmt.Printf("\r%-20s [%s] %3d%%  ", task, bar.String(), percent)

	if percent == 100 {
		fmt.Println()
	}
}
