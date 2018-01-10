package pgr

import (
	"strconv"
	"sync/atomic"
	"text/template"
)

type Bar struct {
	// current must be first member of struct
	// (https://code.google.com/p/go/issues/detail?id=5278)
	current int64
	total   int64
	tmpl    *template.Template
}

func NewBar(total int64, tmpl *template.Template) *Bar {
	return &Bar{current: 0, total: total, tmpl: tmpl}
}

// template MUST NOT print newline.
func (p *Bar) SetTemplate(tmpl *template.Template) *Bar {
	p.tmpl = tmpl
	return p
}

func (p *Bar) SetCurrent(current int64) *Bar {
	atomic.StoreInt64(&p.current, current)
	return p
}

func (p *Bar) Current() int64 {
	v := atomic.LoadInt64(&p.current)
	total := p.Total()
	if v > total {
		return total
	}
	return v
}

// XXX: https://code.google.com/p/go/issues/detail?id=5278
func (p *Bar) SetTotal(total int64) *Bar {
	atomic.StoreInt64(&p.total, total)
	return p
}

func (p *Bar) Total() int64 {
	return atomic.LoadInt64(&p.total)
}

func (p *Bar) Inc() {
	p.Add(1)
}

func (p *Bar) Add(n int64) {
	atomic.AddInt64(&p.current, n)
}

// template MUST NOT print newline.
func NewBarTemplate() *template.Template {
	return template.New("pgr.Poller").Funcs(funcMaps)
}

var funcMaps = template.FuncMap{
	"current": func(value interface{}, args ...string) string {
		if bar, ok := value.(*Bar); ok {
			return strconv.FormatInt(bar.Current(), 10)
		}
		return ""
	},
	"total": func(value interface{}, args ...string) string {
		if bar, ok := value.(*Bar); ok {
			return strconv.FormatInt(bar.Total(), 10)
		}
		return ""
	},
}
