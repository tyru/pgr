package pgr

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
)

type Bar struct {
	// current must be first member of struct
	// (https://code.google.com/p/go/issues/detail?id=5278)
	current  int64
	total    int64
	finished bool

	mu sync.RWMutex

	tmpl         *template.Template
	format       FormatFunc
	finishTmpl   *template.Template
	finishFormat FormatFunc
}

type FormatFunc func(*Bar) string

func NewBar(total int64, tmpl *template.Template) *Bar {
	return &Bar{current: 0, total: total, tmpl: tmpl, finishTmpl: tmpl}
}

func NewBarFunc(total int64, format FormatFunc) *Bar {
	return &Bar{current: 0, total: total, format: format, finishFormat: format}
}

// template MUST NOT print newline.
func (bar *Bar) SetTemplate(tmpl *template.Template) *Bar {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.tmpl = tmpl
	return bar
}

func (bar *Bar) SetCurrent(current int64) *Bar {
	atomic.StoreInt64(&bar.current, current)
	return bar
}

func (bar *Bar) Current() int64 {
	v := atomic.LoadInt64(&bar.current)
	total := bar.Total()
	if v > total {
		return total
	}
	return v
}

// XXX: https://code.google.com/p/go/issues/detail?id=5278
func (bar *Bar) SetTotal(total int64) *Bar {
	atomic.StoreInt64(&bar.total, total)
	return bar
}

func (bar *Bar) Total() int64 {
	return atomic.LoadInt64(&bar.total)
}

func (bar *Bar) Inc() {
	bar.Add(1)
}

func (bar *Bar) Add(n int64) {
	atomic.AddInt64(&bar.current, n)
}

func (bar *Bar) OnFinish(tmpl *template.Template) *Bar {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.finishTmpl = tmpl
	return bar
}

func (bar *Bar) OnFinishFunc(format FormatFunc) *Bar {
	bar.mu.Lock()
	defer bar.mu.Unlock()
	bar.finishFormat = format
	return bar
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
