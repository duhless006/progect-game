package equipment

import "fmt"

type EquipmentType string

const (
	EquipmentTypePickaxe     EquipmentType = "pickaxe"
	EquipmentTypeVentilation EquipmentType = "ventilation"
	EquipmentTypeTrolleys    EquipmentType = "trolleys"
)
const (
	EquipmentPickaxeConst     int64 = 3000
	EquipmentVentilationConst int64 = 15000
	EquipmentTrolleysConst    int64 = 50000
)

type Equipment struct {
	//кирки
	pickaxe bool

	//вентиляция
	ventilation bool

	//вагонетки
	trolleys bool

	//costs map[EquipmentType]int64
}

func NewEquipment() Equipment {
	return Equipment{}
}

func (e *Equipment) BuyPickaxe() {
	e.pickaxe = true

	fmt.Println("pickaxe is purchased")
}

func (e *Equipment) BuyVentilation() {
	e.ventilation = true

	fmt.Println("ventilation is purchased")
}

func (e *Equipment) BuyTrolleys() {
	e.trolleys = true

	fmt.Println("trolleys is purchased")
}

func (e *Equipment) PickaxesPurchased() bool {
	return e.pickaxe
}

func (e *Equipment) VentilationPurchased() bool {
	return e.ventilation
}

func (e *Equipment) TrolleysPurchased() bool {
	return e.trolleys
}

func (e *Equipment) GetStatus() map[EquipmentType]bool {
	return map[EquipmentType]bool{
		EquipmentTypePickaxe:     e.pickaxe,
		EquipmentTypeVentilation: e.ventilation,
		EquipmentTypeTrolleys:    e.trolleys,
	}
}

func (e *Equipment) AllBought() bool {
	return e.pickaxe && e.ventilation && e.trolleys
}
