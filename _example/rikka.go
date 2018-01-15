package main

import (
	"context"
	"time"

	"github.com/tyru/pgr"
)

func main() {
	bar := pgr.NewBarFunc(100, rikka())
	poller := pgr.NewPoller(100 * time.Millisecond).Add(bar)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	poller.Show(ctx)
}

func rikka() pgr.FormatFunc {
	i := 0
	return func(*pgr.Bar) string {
		i++
		switch i % 4 {
		case 0:
			return "(b 回ω・)b"
		case 1:
			return "(σ回ω・)σ"
		case 2:
			return "(q 回ω・)q"
		default:
			return "(ρ回ω・)ρ"
		}
	}
}
