// +build ignore

package main

import (
	"context"
	"text/template"
	"time"

	"github.com/tyru/pgr"
)

func main() {
	b1 := pgr.NewBar(100, parseTemplate(`(b1) {{ current . }}/{{ total . }}`))
	b2 := pgr.NewBar(200, parseTemplate(`(b2) {{ current . }}/{{ total . }}`))
	b3 := pgr.NewBar(100, parseTemplate(`(b3) {{ current . }}/{{ total . }}`))
	go incBy(b1, 30*time.Millisecond)
	go incBy(b2, 20*time.Millisecond)
	go incBy(b3, 40*time.Millisecond)

	poller := pgr.NewPoller(100 * time.Millisecond).Add(b1)

	go func() {
		// Add new progress bar (b2)
		time.Sleep(1 * time.Second)
		poller.Add(b2)

		// Add new progress bar (b3)
		time.Sleep(1 * time.Second)
		poller.Add(b3)
	}()

	poller.Show(context.Background())
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
