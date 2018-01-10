package pgr

import (
	"context"
	"os"
	"time"
)

type Poller struct {
	duration time.Duration
	bars     []*Bar
}

func NewPoller(d time.Duration, bars ...*Bar) *Poller {
	return &Poller{duration: d, bars: bars}
}

func (p *Poller) Show(ctx context.Context) {
	termSave()
	tick := time.NewTicker(p.duration).C
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			termRestore()
			if p.poll() {
				return
			}
		}
	}
}

func (p *Poller) poll() bool {
	doneAll := true
	for i := range p.bars {
		p.bars[i].template.Execute(os.Stdout, p.bars[i])
		if p.bars[i].Current() < p.bars[i].Total() {
			doneAll = false
		}
	}
	return doneAll
}
