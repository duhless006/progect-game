package miners

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	NormalMinerClass MinerClass = "normal"
	BasicMinerSalary            = 50
)

type NormalMiner struct {
	id uuid.UUID

	energy   *atomic.Int64
	power    *atomic.Int64
	interval time.Duration
}

func NewNormalMiner() *NormalMiner {
	const (
		NormalMinerEnergy   = 45
		NormalMinerPower    = 3
		NormalMinerInterval = 2 * time.Second
	)

	energy := &atomic.Int64{}
	power := &atomic.Int64{}

	energy.Add(NormalMinerEnergy)
	power.Add(NormalMinerPower)

	return &NormalMiner{
		id:       uuid.New(),
		energy:   energy,
		power:    power,
		interval: NormalMinerInterval,
	}
}

func (n *NormalMiner) Run(ctx context.Context) <-chan Coal {
	coalChan := make(chan Coal)

	go func() {
		defer close(coalChan)
		fmt.Println("нормальный шахтёр начал работу")

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(n.interval):
			}

			select {
			case <-ctx.Done():
				return
			case coalChan <- Coal(n.power.Load()):
				if new := n.energy.Add(-1); new <= 0 {
					return
				}
			}

		}
	}()
	return coalChan
}

func (n *NormalMiner) Info() MinerInfo {
	return MinerInfo{
		ID:         n.id,
		MinerClass: NormalMinerClass,
		Energy:     n.energy.Load(),
		Power:      n.power.Load(),
		Interval:   time.Duration(n.interval.Seconds()),
	}
}
