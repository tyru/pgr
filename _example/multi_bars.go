// +build ignore

package main

import (
	"context"
	"time"

	"github.com/tyru/pgr"
)

func main() {
	tmpl, err := pgr.NewBarTemplate().Parse(`<<<{{ name . }}>>> {{ current . }}/{{ total . }}{{ println }}`)
	if err != nil {
		panic(err)
	}

	p1 := pgr.NewBar("p1").SetTemplate(tmpl)
	p2 := pgr.NewBar("p2").SetTotal(200)
	p3 := pgr.NewBar("p3")
	go incBy(p1, 30*time.Millisecond)
	go incBy(p2, 20*time.Millisecond)
	go incBy(p3, 40*time.Millisecond)

	poller := pgr.NewPoller(100*time.Millisecond, p1, p2, p3)
	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	// defer cancel()
	poller.Show(ctx)
}

func incBy(p *pgr.Bar, d time.Duration) {
	for {
		p.Inc()
		time.Sleep(d)
	}
}
