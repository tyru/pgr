package pgr

import (
	"context"
	"os"
	"time"
)

type Poller time.Duration

func NewPoller(d time.Duration) *Poller {
	p := Poller(d)
	return &p
}

func (poll *Poller) Show(ctx context.Context, bars ...*Bar) {
	termSave()
	tick := time.NewTicker(time.Duration(*poll)).C
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			termRestore()
			if poll.poll(bars) {
				return
			}
		}
	}
}

func (poll *Poller) poll(bars []*Bar) bool {
	done := true
	for i := range bars {
		bars[i].template.Execute(os.Stdout, bars[i])
		if bars[i].Current() < bars[i].Total() {
			done = false
		}
	}
	return done
}
