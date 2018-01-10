// +build ignore

package main

import (
	"bytes"
	"context"
	"math"
	"strings"
	"text/template"
	"time"

	"github.com/tyru/pgr"
)

func main() {
	b1 := pgr.NewBar(math.MaxInt64, hasiruTemplate())
	b2 := pgr.NewBar(math.MaxInt64, uwaaaaTemplate())
	b3 := pgr.NewBar(math.MaxInt64, parseTemplate(`(b3) {{ current . }}/{{ total . }}`))
	go incBy(b1, 30*time.Millisecond)
	go incBy(b2, 20*time.Millisecond)
	go incBy(b3, 40*time.Millisecond)

	poller := pgr.NewPoller(100 * time.Millisecond).Add(b1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Add new progress bar (b2)
		time.Sleep(2 * time.Second)
		poller.Add(b2)

		// Add new progress bar (b3)
		time.Sleep(2 * time.Second)
		poller.Add(b3)

		// Speed up 10x faster
		time.Sleep(2 * time.Second)
		poller.SetDuration(10 * time.Millisecond)

		// end
		time.Sleep(2 * time.Second)
		cancel()
	}()

	poller.Show(ctx)
}

func hasiruTemplate() *template.Template {
	t := pgr.NewBarTemplate().Funcs(template.FuncMap{
		"hasiru": func() func() string {
			forward := true
			dash := 0
			const maxDash = 10
			return func() string {
				var aa string
				if forward {
					if dash == 0 {
						aa = "┏( ^o^)┛"
					} else {
						aa = strings.Repeat("　", dash-1) + "三┏( ^o^)┛"
					}
				} else {
					if dash == 0 {
						aa = strings.Repeat("　", maxDash) + "┗(^o^ )┓"
					} else {
						aa = strings.Repeat("　", maxDash-dash) + "┗(^o^ )┓三"
					}
				}
				if dash >= maxDash {
					forward = !forward
					dash = 0
				} else {
					dash++
				}
				return aa
			}
		}(),
	})
	return template.Must(t.Parse(`{{ hasiru }}`))
}

func uwaaaaTemplate() *template.Template {
	t := pgr.NewBarTemplate().Funcs(template.FuncMap{
		"uwaaaa": func() func() string {
			start := 0
			const b1 = "▂"
			const b2 = "▅"
			const b3 = "▇"
			const b4 = "▇"
			const b5 = "▓"
			const b6 = "▒"
			const b7 = "░"
			wave := []string{b1, b2, b3, b4, b5, b6, b7, b6, b5, b4, b3, b2}
			const face = " ('ω')"
			return func() string {
				var buf bytes.Buffer

				// wave 1
				j := start
				for i := 0; i < len(wave); i++ {
					buf.WriteString(wave[j])
					j = (j + 1) % len(wave)
				}

				buf.WriteString(face)

				// wave 2: reverse of wave 1
				j = start
				for i := 0; i < len(wave); i++ {
					buf.WriteString(wave[len(wave)-1-j])
					j = (j + 1) % len(wave)
				}

				start = (start + 1) % len(wave)
				return buf.String()
			}
		}(),
	})
	return template.Must(t.Parse(`{{ uwaaaa }}`))
}

func parseTemplate(format string) *template.Template {
	return template.Must(pgr.NewBarTemplate().Parse(format))
}

func incBy(p *pgr.Bar, d time.Duration) {
	for {
		p.Inc()
		time.Sleep(d)
	}
}
