package checkerutil

import (
	"os"
	"strings"
)

// Combo relates to loading combo(s) and using *Combo.GetCombo() to get a combo from the list in a thread safe way with no risk of a race condition
type Combo struct {
	ComboChan chan string
}

// LoadCombo returns a *Combo of which contains *Combo.ComboChan, you use *Combo.GetCombo() to interact with this channel
func LoadCombo(FileName, Deliminator string) (*Combo, error) {
	f, err := os.ReadFile(FileName)
	if err != nil {
		return nil, err
	}
	a := strings.Split(string(f), Deliminator)
	c := &Combo{ComboChan: make(chan string, len(a))}
	go func(a []string) {
		for _, p := range a {
			c.ComboChan <- p
		}
	}(a)
	return c, nil
}

// LoadCombosFromDir is the same as LoadCombo except can load multiple combos from one Dir
func LoadCombosFromDir(Dir, Deliminator string) (*Combo, error) {
	c := &Combo{ComboChan: make(chan string)}
	if !strings.HasSuffix(Dir, "/") {
		Dir += "/"
	}
	files, err := os.ReadDir(Dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		ff, err := os.ReadFile(Dir + f.Name())
		if err != nil {
			return nil, err
		}
		a := strings.Split(string(ff), Deliminator)
		for _, p := range a {
			c.ComboChan <- p
		}
	}
	return c, nil
}

// GetCombo returns a single entry into the *Combo.ComboChan channel
// If the combo channel is empty, it will return an empty string
func (c *Combo) GetCombo() string {
	if len(c.ComboChan) == 0 {
		return ""
	}
	return <-c.ComboChan
}
