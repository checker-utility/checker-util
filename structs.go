package checkerutil

import (
	"os"
	"sync"
	"time"
)

// Output holds information like the file name, delay to write, etc
type Output struct {
	FileName string
	Delay    time.Duration
	ID       string
	File     *os.File
	Input    chan string
	InputNum struct {
		Num   int
		Mutex sync.Mutex
	}
}
