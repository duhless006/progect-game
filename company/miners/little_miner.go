package miners

import (
	"context"
	"fmt"
	"sync/atomic"

	"time"

	"github.com/google/uuid"
)

const (
	LittleMinerClass  MinerClass = "little"
	LittleMinerSalary            = 5
)

type LittleMiner struct {
	id       uuid.UUID
	energy   *atomic.Int64
	power    *atomic.Int64
	interval time.Duration
}

func NewLittleMiner() *LittleMiner {
	const (
		littleMinerEnergy   = 30
		littleMinerPower    = 1
		littleMinerInterval = 3 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(littleMinerEnergy)
	power.Add(littleMinerPower)

	return &LittleMiner{
		id:       uuid.New(),
		energy:   energy,
		power:    power,
		interval: littleMinerInterval,
	}
}

func (l *LittleMiner) Run(ctx context.Context) <-chan Coal {
	coalChan := make(chan Coal)

	go func() {
		defer close(coalChan)
		defer fmt.Println("Шахтёр начал работу")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Шахтёр остановлен")
				return
			case <-time.After(l.interval):
			}
			select {
			case <-ctx.Done():
				return
			case coalChan <- Coal(l.power.Load()):
				if new := l.energy.Add(-1); new <= 0 {
					return
				}
			}
		}
	}()
	return coalChan
}

func (l *LittleMiner) Info() MinerInfo {
	return MinerInfo{
		ID:         l.id,
		MinerClass: LittleMinerClass,
		Energy:     l.energy.Load(),
		Power:      l.power.Load(),
		Interval:   time.Duration(l.interval.Seconds()),
	}
}
