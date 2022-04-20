package checkerutil

import (
	"os"
	"strings"
	"time"
)

// Checker holds a map of all the different possible writes as well as all the configuration for the checker utility and the CPM calculations
type Checker struct {
	Outputs map[string]*Output
	CPM     *CPMStruct
	Console *ConsoleTitle
	Dir     string
}

// MakeChecker takes in nothing, it returns a (mostly empty) *Checker, you should use the methods attached to configure the checker.
func MakeChecker() *Checker {
	return &Checker{
		Outputs: make(map[string]*Output),
		CPM: &CPMStruct{
			CPMIn: make(chan int),
		},
		Console: &ConsoleTitle{},
	}
}

// SetDir will set the default directory that will be applied to both automatically and manually added outputs.
func (c *Checker) SetDir(Dir string) {
	c.Dir = Dir
}

// ConfigureConsoleTitle configures the *Checker.ConsoleTitle struct
func (c *Checker) ConfigureConsoleTitle(Activated bool, Delay time.Duration, Format string) {
	c.Console = &ConsoleTitle{
		Activated: Activated,
		Delay:     Delay,
		Format:    Format,
	}
	if Activated {
		c.handleCPM()
	}
}

//AddOutput explicitly adds an output to the outputs list
/*
   ID referrs to how you later use the outputter when calling *Checker.Input
   FileName is optional, if one is supplied, it will use it, but if not, it will just use ID and append .txt to the end (if the file is not already created, it will create it)
   Delay referrs to how often it will write
*/
func (c *Checker) AddOutput(ID, FileName string, Delay time.Duration) error {
	if FileName == "" {
		FileName = ID + ".txt"
	}
	if c.Dir != "" && !strings.HasPrefix(FileName, c.Dir) {
		FileName = c.Dir + FileName
	}
	f, err := os.OpenFile(FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// If this happens, it basically makes the entire utility not work. The only logical way this could happen is if the program doesn't have permissions or some other obscure thing.
		panic(err)
	}
	c.Outputs[ID] = &Output{
		FileName: FileName,
		ID:       ID,
		Delay:    Delay,
		File:     f,
		Input:    make(chan string),
	}
	go func(ID string) {
		o := c.Outputs[ID]
		toWrite := ""
		go func(o *Output) {
			for i := range o.Input {
				toWrite += i
			}
		}(o)
		for {
			time.Sleep(o.Delay)
			if toWrite != "" {
				_, err := o.File.WriteString(toWrite)
				if err == nil {
					toWrite = ""
				}
			}
		}
	}(ID)
	return nil
}

// Input will actually write to files and use the *Output struct
// If you input an invalid ID, it will add that ID to the outputs list automatically, the default delay will be 1 second and the rules from AddOutput do apply.
func (c *Checker) Input(ID, Input string) {
	if _, ok := c.Outputs[ID]; !ok {
		c.AddOutput(ID, "", 1*time.Second)
	}
	if !strings.HasSuffix(Input, "\n") {
		Input += "\n"
	}
	o := c.Outputs[ID]
	o.Input <- Input
}
