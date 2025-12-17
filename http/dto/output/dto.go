package output

import (
	"progect-game/company"
	"progect-game/company/equipment"
	"progect-game/company/miners"

	"github.com/google/uuid"
)

type EquipmentDto struct {
	Pickaxe bool `json:"pickaxe"`

	Ventilation bool `json:"ventilation"`

	Trolleys bool `json:"trolleys"`
}

func NewEquipmentDTO(equipment equipment.Equipment) EquipmentDto {
	return EquipmentDto{
		Pickaxe:     equipment.PickaxesPurchased(),
		Ventilation: equipment.VentilationPurchased(),
		Trolleys:    equipment.TrolleysPurchased(),
	}
}

type MinersByTypeDTO map[uuid.UUID]miners.MinerInfo

func NewMinersByTypeDTO(m map[uuid.UUID]miners.Miner) MinersByTypeDTO {
	info := make(map[uuid.UUID]miners.MinerInfo)

	for id, miner := range m {
		info[id] = miner.Info()
	}
	return info
}

type AllMinerDTO map[miners.MinerClass]MinersByTypeDTO

func NewAllMinerDTO(m map[miners.MinerClass]map[uuid.UUID]miners.Miner) AllMinerDTO {
	info := make(map[miners.MinerClass]MinersByTypeDTO)

	for minersClass, minerMap := range m {
		info[minersClass] = NewMinersByTypeDTO(minerMap)
	}
	return info
}

type CompanyStatisticDTO struct {
	Balance             int64 `json:"balance"`
	TotalEarned         int64
	TotalMinerStatistic map[miners.MinerClass]int
	TotalTime           string
}

func NewCompanyStatisticDTO(companyStatistic company.CompanyStatistics) CompanyStatisticDTO {
	return CompanyStatisticDTO{
		Balance:             companyStatistic.Balance(),
		TotalEarned:         companyStatistic.TotalEarned(),
		TotalMinerStatistic: companyStatistic.TotalMinersStatistics(),
		TotalTime:           companyStatistic.TimeToComple(),
	}
}

type MinersSalariesDTO struct {
	LittleSalary   int64 `json:"littlePrice"`
	NormalSalary   int64 `json:"normalPrice"`
	PowerfulSalary int64 `json:"powerfulPrice"`
}

func NewminersSalariesDTO(
	littleMinerSalary int64,
	normalMinerSalary int64,
	powerfulMinerSalary int64,
) MinersSalariesDTO {
	return MinersSalariesDTO{
		LittleSalary:   littleMinerSalary,
		NormalSalary:   normalMinerSalary,
		PowerfulSalary: powerfulMinerSalary,
	}
}

type EquipmntPriceDTO struct {
	PickaxePrice     int64 `json:"pickaxePrice"`
	VentilationPrice int64 `json:"ventilationPrice"`
	TroleysPrice     int64 `json:"troleysPrice"`
}

func NewEquipmntPriceDTO(
	pickaxePrice int64,
	ventilationPrice int64,
	troleysPrice int64,
) EquipmntPriceDTO {
	return EquipmntPriceDTO{
		PickaxePrice:     pickaxePrice,
		VentilationPrice: ventilationPrice,
		TroleysPrice:     troleysPrice,
	}
}
