package main

import (
	"context"
	"time"

	"github.com/tyru/pgr"
)

// This is the simplest example
func main() {
	bar := pgr.NewBar(100, pgr.DefaultTemplate)
	go incBy(bar, 30*time.Millisecond)

	poller := pgr.NewPoller(100 * time.Millisecond).Add(bar)
	poller.Show(context.Background())
}

func incBy(p *pgr.Bar, d time.Duration) {
	for {
		p.Inc()
		time.Sleep(d)
	}
}
