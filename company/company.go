package company

import (
	"context"
	"fmt"
	"progect-game/company/equipment"
	"progect-game/company/miners"
	"progect-game/models"

	"sync"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	incomeCh  chan miners.Coal
	miners    map[miners.MinerClass]map[uuid.UUID]miners.Miner
	equipment equipment.Equipment
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc

	//------
	statisctics *CompanyStatistics
}

func NewCompany(ctx context.Context) *Company {

	ctx, cancel := context.WithCancel(ctx)

	c := &Company{
		incomeCh:  make(chan miners.Coal),
		miners:    make(map[miners.MinerClass]map[uuid.UUID]miners.Miner),
		equipment: equipment.NewEquipment(),
		ctx:       ctx,
		cancel:    cancel,

		statisctics: NewCompanyStatistics(),
	}

	go c.baseIncome()
	go c.collectIncome()

	return c
}

// HireMiner нанимаем наши чебуреков
func (c *Company) HireMiner(minerClass miners.MinerClass) (miners.Miner, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var miner miners.Miner
	switch minerClass {
	case miners.LittleMinerClass:
		if c.statisctics.balance.Load() >= miners.LittleMinerSalary {
			miner = miners.NewLittleMiner()
			c.statisctics.balance.Add(-miners.LittleMinerSalary)
		} else {
			return nil, ErrInsifficientFunds
		}
	case miners.NormalMinerClass:
		if c.statisctics.balance.Load() >= miners.BasicMinerSalary {
			miner = miners.NewNormalMiner()
			c.statisctics.balance.Add(-miners.BasicMinerSalary)
		} else {
			return nil, ErrInsifficientFunds
		}
	case miners.PowerfulMinerClass:
		if c.statisctics.balance.Load() >= miners.PowerfulMinerSalary {
			miner = miners.NewPowerfulMiner()
			c.statisctics.balance.Add(-miners.PowerfulMinerSalary)
		} else {
			return nil, ErrInsifficientFunds
		}
	default:
		return nil, ErrUnknokwnMinerType
	}

	info := miner.Info()
	if c.miners[info.MinerClass] == nil {
		c.miners[info.MinerClass] = make(map[uuid.UUID]miners.Miner)
	}
	c.miners[minerClass][info.ID] = miner

	coalChan := miner.Run(c.ctx)

	go func() {
		select {
		case <-c.ctx.Done():
			return
		default:
			for v := range coalChan {
				select {
				case c.incomeCh <- v:
				case <-c.ctx.Done():
					return
				}
			}
		}
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.miners[minerClass], info.ID)

		if len(c.miners[minerClass]) == 0 {
			delete(c.miners, minerClass)
		}

	}()

	c.statisctics.totalMinersStatistics[minerClass]++
	return miner, nil
}

// GetAllMiners посмотреть всех майнеров
func (c *Company) GetAllMiners() map[miners.MinerClass]map[uuid.UUID]miners.Miner {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tmp := make(map[miners.MinerClass]map[uuid.UUID]miners.Miner)
	for minerClass, minerMap := range c.miners {
		tmp[minerClass] = make(map[uuid.UUID]miners.Miner)

		for k, v := range minerMap {
			tmp[minerClass][k] = v
		}
	}
	return tmp
}

// GetMinersByType посмотреть тип майнера
func (c *Company) GetMinersByType(minerType miners.MinerClass) map[uuid.UUID]miners.Miner {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.miners == nil {
		return make(map[uuid.UUID]miners.Miner)
	}

	minersMap, ok := c.miners[minerType]
	if !ok || minersMap == nil {
		return make(map[uuid.UUID]miners.Miner)
	}

	tmp := make(map[uuid.UUID]miners.Miner)
	for k, v := range c.miners[minerType] {
		tmp[k] = v
	}
	return tmp
}

// BuyEquipment- покупка оборудования
func (c *Company) BuyEquipment(equipmentType equipment.EquipmentType) (equipment.Equipment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if equipmentType == equipment.EquipmentTypePickaxe {
		if c.statisctics.balance.Load() >= equipment.EquipmentPickaxeConst {
			c.equipment.BuyPickaxe()
			c.statisctics.balance.Add(-equipment.EquipmentPickaxeConst)
		} else {
			return equipment.Equipment{}, ErrInsifficientFunds
		}
	} else if equipmentType == equipment.EquipmentTypeTrolleys {
		if c.statisctics.balance.Load() >= equipment.EquipmentTrolleysConst {
			c.equipment.BuyTrolleys()
			c.statisctics.balance.Add(-equipment.EquipmentTrolleysConst)
		} else {
			return equipment.Equipment{}, ErrInsifficientFunds
		}
	} else if equipmentType == equipment.EquipmentTypeVentilation {
		if c.statisctics.balance.Load() >= equipment.EquipmentVentilationConst {
			c.equipment.BuyVentilation()
			c.statisctics.balance.Add(-equipment.EquipmentVentilationConst)
		} else {
			return equipment.Equipment{}, ErrInsifficientFunds
		}
	}
	return c.equipment, nil
}

// Получить оборудование компании
func (c *Company) GetEquipment() equipment.Equipment {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.equipment //возможно нужен будет указатель
}

// Завершить игру и получить финальную статистику
func (c *Company) FinishGame() (*CompanyStatistics, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.equipment.AllBought() {
		return &CompanyStatistics{}, ErrNotAllEquipmentPurchased
	}

	if c.statisctics.completedTime != nil {
		return &CompanyStatistics{}, ErrCompanyStopped
	}

	c.cancel()
	now := time.Now()
	c.statisctics.completedTime = &now

	return c.statisctics, nil
}

func (c *Company) GetStatistics() CompanyStatistics {
	return *c.statisctics
}

func (c *Company) collectIncome() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case income, ok := <-c.incomeCh:
			if !ok {
				return
			}
			c.statisctics.balance.Add(int64(income))
			fmt.Println("balance:", c.statisctics.balance.Load())
			c.statisctics.totalEarned.Add(int64(income))
		}
	}
}

// Постоянно генерирует базовый доход
func (c *Company) baseIncome() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			select {
			case c.incomeCh <- 1:
			case <-c.ctx.Done():
				return
			default:
			}
		}
	}

}

// GetState возвращает текущее состояние компании
func (c *Company) GetState() models.GameState {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Считаем общее количество шахтёров
	totalMiners := 0
	for _, minerMap := range c.miners {
		totalMiners += len(minerMap)
	}

	return models.GameState{
		CompanyID:   "default_company",
		Money:       c.statisctics.balance.Load(),
		TotalEarned: c.statisctics.totalEarned.Load(),
		MinersCount: totalMiners,
		Pickaxe:     c.equipment.PickaxesPurchased(),
		Ventilation: c.equipment.VentilationPurchased(),
		Trolleys:    c.equipment.TrolleysPurchased(),
		CreatedAt:   time.Now(),
	}
}

// SetState загружает состояние в компанию
func (c *Company) SetState(state models.GameState) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.statisctics.balance.Store(state.Money)
	c.statisctics.totalEarned.Store(state.TotalEarned)

	if state.Pickaxe && !c.equipment.PickaxesPurchased() {
		c.equipment.BuyPickaxe()
	}
	if state.Ventilation && !c.equipment.VentilationPurchased() {
		c.equipment.BuyVentilation()
	}
	if state.Trolleys && !c.equipment.TrolleysPurchased() {
		c.equipment.BuyTrolleys()
	}

	return nil
}

// GetMoney возвращает текущий баланс
func (c *Company) GetMoney() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.statisctics.balance.Load()
}

// GetTotalEarned возвращает общий заработок
func (c *Company) GetTotalEarned() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.statisctics.totalEarned.Load()
}

// GetMinerCount возвращает количество шахтёров
func (c *Company) GetMinerCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := 0
	for _, minerMap := range c.miners {
		total += len(minerMap)
	}
	return total
}
