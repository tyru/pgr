package pgr

import (
	"context"
	"os"
	"sync"
	"time"
)

type Poller struct {
	duration time.Duration
	bars     []*Bar

	mu sync.RWMutex
}

func NewPoller(d time.Duration) *Poller {
	return &Poller{duration: d}
}

func (p *Poller) Add(bars ...*Bar) *Poller {
	if len(bars) == 0 {
		return p
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bars = append(p.bars, bars...)
	return p
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
