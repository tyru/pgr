package pgr

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
)

type Poller struct {
	duration time.Duration
	bars     []*Bar

	mu      sync.RWMutex
	running bool
	out     io.Writer
}

func NewPoller(d time.Duration) *Poller {
	return &Poller{duration: d, out: colorable.NewColorableStdout()}
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

var errCantChangeOut = errors.New("cannot change out while running")

func (p *Poller) SetOut(out io.Writer) error {
	p.mu.RLock()
	if p.running {
		return errCantChangeOut
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.out = out
	return nil
}

var ErrCanceled = errors.New("canceled by context")

func (p *Poller) Show(ctx context.Context) error {
	p.mu.Lock()
	p.running = true
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		p.running = false
		p.mu.Unlock()
	}()

	lines := 0
	for {
		select {
		case <-ctx.Done():
			return ErrCanceled
		case <-time.NewTimer(p.duration).C:
			if lines > 0 {
				termPrevLine(p.out, lines)
			}
			if l, err := p.poll(); err == nil {
				return nil
			} else if err != errUnfinished {
				return err
			} else {
				lines = l
			}
		}
	}
}

var errUnfinished = errors.New("not finished yet")

// poll() returns nil error if all bars are finished.
func (p *Poller) poll() (lines int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, bar := range p.bars {
		termClearLine(p.out)
		if bar.tmpl != nil {
			if err := bar.tmpl.Execute(p.out, bar); err != nil {
				return 0, err
			}
		} else /* if bar.format != nil */ {
			if s := bar.format(bar); s != "" {
				if _, err := p.out.Write([]byte(s)); err != nil {
					return 0, err
				}
			}
		}
		if _, err := p.out.Write([]byte{byte('\n')}); err != nil {
			return 0, err
		}
		if bar.Current() < bar.Total() {
			err = errUnfinished
		}
	}
	return len(p.bars), err
}

func (p *Poller) SetDuration(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.duration = d
}
