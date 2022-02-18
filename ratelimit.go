package goratelimit

import (
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

var (
	m         = cmap.New()
	allow     = true
	imRunning = 0
)

// RateLimit: Limits Token to rate, returns true/false
func RateLimit(Token string, rate, rtimer int) bool {
	count := 1
	if imRunning == 0 {
		imRunning = 1
		go func() {
			fmt.Println("Running")
			for {
				time.Sleep(time.Duration(rtimer) * time.Second)
				start := time.Now()
				var currentBlockCount int
				for mlist, cont := range m.Items() {
					if cont.(int) > 0 {
						m.Set(mlist, cont.(int)-1)
					}
					if cont.(int) == 0 {
						m.Remove(mlist)
					}
					if cont.(int) >= rate {
						currentBlockCount++
					}
				}
				totalList := len(m.Items())
				done := time.Since(start)
				fmt.Printf("Session: %v Blocked: %v ProcessTime: %v\n", totalList, currentBlockCount, done)
			}
		}()
	}

	if chkAccess, ok := m.Get(Token); ok {
		if chkAccess.(int) >= rate {
			return false
		}
	}
	countNew, ok := m.Get(Token)
	if ok {
		count = countNew.(int) + 1
	}
	m.Set(Token, count)
	return allow
}
