package miners

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	PowerfulMinerClass  MinerClass = "powerful"
	PowerfulMinerSalary            = 450
)

type PowerfulMiner struct {
	id       uuid.UUID
	energy   *atomic.Int64
	power    *atomic.Int64
	interval time.Duration
}

func NewPowerfulMiner() *PowerfulMiner {
	const (
		powerfulMinerEnegy    = 60
		powerfulMinerPower    = 10
		powerfulMinerInterval = 1 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(powerfulMinerEnegy)
	power.Add(powerfulMinerPower)

	return &PowerfulMiner{

		id:       uuid.New(),
		energy:   energy,
		power:    power,
		interval: powerfulMinerInterval,
	}
}

func (p *PowerfulMiner) Run(ctx context.Context) <-chan Coal {
	coalChan := make(chan Coal)

	go func() {
		defer close(coalChan)
		defer fmt.Println("Сильный шахтёр начал  работу")

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(p.interval):
			}
			select {
			case <-ctx.Done():
				return
			case coalChan <- Coal(p.power.Load()):
				p.power.Add(3)

				if new := p.energy.Add(-1); new <= 0 {
					return
				}
			}
		}
	}()
	return coalChan
}

func (p *PowerfulMiner) Info() MinerInfo {
	return MinerInfo{
		ID:         p.id,
		MinerClass: PowerfulMinerClass,
		Energy:     p.energy.Load(),
		Power:      p.power.Load(),
		Interval:   time.Duration(p.interval.Seconds()),
	}
}
