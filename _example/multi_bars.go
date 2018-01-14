// +build ignore

package main

import (
	"context"
	"text/template"
	"time"

	"github.com/tyru/pgr"
)

func main() {
	b1 := pgr.NewBar(100, parseTemplate(`(b1) {{ current . }}/{{ total . }} {{ bar . "[" "=" ">" " " "]" 70 }}`)).
		OnFinish(parseTemplate(`(b1) {{ current . }}/{{ total . }} Finished!`))
	b2 := pgr.NewBar(200, parseTemplate(`(b2) {{ current . }}/{{ total . }} {{ bar . "[" "=" ">" " " "]" 70 }}`)).
		OnFinish(parseTemplate(`(b2) {{ current . }}/{{ total . }} Finished!`))
	b3 := pgr.NewBar(300, parseTemplate(`(b3) {{ current . }}/{{ total . }} {{ bar . "[" "=" ">" " " "]" 70 }}`)).
		OnFinish(parseTemplate(`(b3) {{ current . }}/{{ total . }} Finished!`))

	poller := pgr.NewPoller(100 * time.Millisecond).Add(b1)

	go func() {
		go incBy(b1, 30*time.Millisecond)

		// Add new progress bar (b2)
		time.Sleep(1 * time.Second)
		go incBy(b2, 20*time.Millisecond)
		poller.Add(b2)

		// Add new progress bar (b3)
		time.Sleep(1 * time.Second)
		go incBy(b3, 10*time.Millisecond)
		poller.Add(b3)
	}()

	poller.Show(context.Background())
}

func parseTemplate(format string) *template.Template {
	return template.Must(pgr.NewTemplate().Parse(format))
}

func incBy(p *pgr.Bar, d time.Duration) {
	for {
		p.Inc()
		time.Sleep(d)
	}
}
