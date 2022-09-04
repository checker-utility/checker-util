package checkerutil

import (
	"sync"
	"time"
)

// CPMStruct holds all data for calculating the CPM
type CPMStruct struct {
	CPS struct {
		CPSMutex sync.Mutex
		CPS      int
	}
	CPMIn    chan int
	CPMArray []int
	CPMMutex sync.Mutex
}

// CalculateCPM calculates the CPM, making sure to lock the CPM arrays locker so you don't run into a race condition
func (c *Checker) CalculateCPM() int {
	c.CPM.CPMMutex.Lock()
	end := 0
	for _, a := range c.CPM.CPMArray {
		end += a
	}
	c.CPM.CPMMutex.Unlock()
	return end
}

// handleCPM sets the CPM.CPMArray using CPS, CPS is calculated using CPM.CPMIn and CPM.CPMOut
func (c *Checker) handleCPM() {
	go func(c *Checker) {
		for {
			time.Sleep(1 * time.Second)
			c.CPM.CPS.CPSMutex.Lock()
			c.CPM.CPMMutex.Lock()
			if len(c.CPM.CPMArray) != 60 {
				c.CPM.CPMArray = append(c.CPM.CPMArray, c.CPM.CPS.CPS)
				c.CPM.CPS.CPS = 0
				c.CPM.CPS.CPSMutex.Unlock()
				c.CPM.CPMMutex.Unlock()
				continue
			}
			cp := []int{}
			for i := 1; i < len(c.CPM.CPMArray); i++ {
				cp = append(cp, c.CPM.CPMArray[i])
			}
			cp = append(cp, c.CPM.CPS.CPS)
			c.CPM.CPMArray = cp
			c.CPM.CPS.CPSMutex.Unlock()
			c.CPM.CPMMutex.Unlock()
		}
	}(c)
	for range c.CPM.CPMIn {
		c.CPM.CPS.CPSMutex.Lock()
		c.CPM.CPS.CPS++
		c.CPM.CPS.CPSMutex.Unlock()
	}
}
