package timewheel

import (
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	channel := make(chan time.Time)
	beginTime := time.Now()

	Delay(time.Second, "", func() { channel <- time.Now() })
	execAt := <-channel
	delayDuration := execAt.Sub(beginTime)

	if delayDuration < time.Second || delayDuration > 3*time.Second {
		t.Error("wrong execute time")
	}
}
