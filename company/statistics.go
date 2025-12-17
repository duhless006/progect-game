package company

import (
	"progect-game/company/miners"
	"sync/atomic"
	"time"
)

type CompanyStatistics struct {
	balance     *atomic.Int64
	totalEarned *atomic.Int64

	totalMinersStatistics map[miners.MinerClass]int

	createdTime   time.Time
	completedTime *time.Time
}

func NewCompanyStatistics() *CompanyStatistics {
	return &CompanyStatistics{
		balance:               &atomic.Int64{},
		totalEarned:           &atomic.Int64{},
		totalMinersStatistics: make(map[miners.MinerClass]int),
		createdTime:           time.Now(),
	}
}

func (c *CompanyStatistics) Balance() int64 {
	if c == nil || c.balance == nil {
		return 0
	}
	return c.balance.Load()
}

func (c *CompanyStatistics) TotalEarned() int64 {
	return c.totalEarned.Load()
}

func (c *CompanyStatistics) TimeToComple() string {
	if c.completedTime == nil {
		return ""
	}

	return c.completedTime.Sub(c.createdTime).String()
}

func (c *CompanyStatistics) TotalMinersStatistics() map[miners.MinerClass]int {

	tmp := make(map[miners.MinerClass]int)
	for k, v := range c.totalMinersStatistics {
		tmp[k] = v
	}

	return tmp
}
