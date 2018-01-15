package pgr

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
)

type Poller struct {
	duration time.Duration
	lessFunc LessFunc
	bars     []*Bar

	mu      sync.RWMutex
	running bool

	out io.Writer
}

type LessFunc func(*Bar, *Bar) bool

func NewPoller(d time.Duration) *Poller {
	return &Poller{duration: d, out: colorable.NewColorableStdout()}
}

func (p *Poller) SetSorter(less LessFunc) *Poller {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lessFunc = less
	return p
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

func (p *Poller) SetDuration(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.duration = d
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
			if lines > 0 {
				termPrevLine(p.out, lines)
			}
			p.Poll()
			return ErrCanceled
		case <-time.NewTimer(p.duration).C:
			if lines > 0 {
				termPrevLine(p.out, lines)
			}
			if l, err := p.Poll(); err == nil {
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
func (p *Poller) Poll() (lines int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.lessFunc != nil {
		sort.SliceStable(p.bars, func(i, j int) bool {
			return p.lessFunc(p.bars[i], p.bars[j])
		})
	}

	for i := range p.bars {
		termClearLine(p.out)
		if err := p.drawLine(p.bars[i]); err != nil {
			return 0, err
		}
		if _, err := p.out.Write([]byte{byte('\n')}); err != nil {
			return 0, err
		}
		if p.bars[i].Current() < p.bars[i].Total() {
			err = errUnfinished
		}
	}
	return len(p.bars), err
}

func (p *Poller) drawLine(bar *Bar) error {
	tmpl := bar.tmpl
	format := bar.format
	if bar.finished != 0 || bar.Current() >= bar.Total() {
		bar.finished = 1
		tmpl = bar.finishTmpl
		format = bar.finishFormat
	}
	if tmpl != nil {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, bar); err != nil {
			return err
		}
		s := bytes.Replace(buf.Bytes(), []byte{byte('\n')}, []byte{}, -1)
		if _, err := p.out.Write([]byte(s)); err != nil {
			return err
		}
	} else if format != nil {
		if s := format(bar); s != "" {
			buf := bytes.Replace([]byte(s), []byte{byte('\n')}, []byte{}, -1)
			if _, err := p.out.Write(buf); err != nil {
				return err
			}
		}
	}
	return nil
}
