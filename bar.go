package pgr

import (
	"strconv"
	"strings"
	"sync/atomic"
	"text/template"
)

type Bar struct {
	// current must be first member of struct
	// (https://code.google.com/p/go/issues/detail?id=5278)
	current int64
	total   int64
	tmpl    *template.Template
	format  FormatFunc
}

type FormatFunc func(*Bar) string

func NewBar(total int64, tmpl *template.Template) *Bar {
	return &Bar{current: 0, total: total, tmpl: tmpl}
}

func NewBarFunc(total int64, format FormatFunc) *Bar {
	return &Bar{current: 0, total: total, format: format}
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
func NewTemplate() *template.Template {
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
	"percent": func(value interface{}) string {
		if bar, ok := value.(*Bar); ok {
			percent := int(float64(bar.Current()) / float64(bar.Total()) * 100)
			return strconv.Itoa(percent) + "%"
		}
		return ""
	},
	"bar": func(value interface{}, prefix, complete, current, incomplete, suffix string, opt ...int) string {
		if bar, ok := value.(*Bar); ok {
			col := 80
			if len(opt) > 0 {
				col = opt[0]
			}
			ccWidth := col - len(prefix+current+suffix)
			if ccWidth <= 0 {
				return "" // no space
			}

			// len(complete) = n, len(uncomplete) = m
			//
			//   n : n + m = bar.Current() : bar.Total()
			//   (n + m) * bar.Current() = n * bar.Total()
			//
			// col - len(prefix+current+suffix) = ccWidth
			// ccWidth = n + m
			//
			//   ccWidth * bar.Current() = n * bar.Total()
			//   n = ccWidth * bar.Current() / bar.Total()

			completeCount := int(float64(ccWidth) * float64(bar.Current()) / float64(bar.Total()))
			incompleteCount := ccWidth - completeCount
			return prefix +
				strings.Repeat(complete, completeCount) +
				current +
				strings.Repeat(incomplete, incompleteCount) +
				suffix
		}
		return ""
	},
}

// `{{ percent . }} {{ bar . "[" "=" ">" " " "]" }}`
var DefaultTemplate = template.Must(NewTemplate().Parse(`{{ percent . }} {{ bar . "[" "=" ">" " " "]" }}`))
