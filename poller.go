package pgr

import (
	"context"
	"errors"
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

var ErrCanceled = errors.New("canceled by context")

func (p *Poller) Show(ctx context.Context) error {
	termSave()
	tick := time.NewTicker(p.duration).C
	for {
		select {
		case <-ctx.Done():
			return ErrCanceled
		case <-tick:
			termRestore()
			if err := p.poll(); err == nil {
				return nil
			} else if err != errUnfinished {
				return err
			}
		}
	}
}

var errUnfinished = errors.New("not finished yet")

// poll() returns nil error if all bars are finished.
func (p *Poller) poll() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := range p.bars {
		if err := p.bars[i].template.Execute(os.Stdout, p.bars[i]); err != nil {
			return err
		}
		if p.bars[i].Current() < p.bars[i].Total() {
			err = errUnfinished
		}
	}
	return err
}
