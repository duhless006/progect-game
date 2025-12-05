package company

import "errors"

var ErrUnknokwnMinerType = errors.New("неизвестный тип майнера")
var ErrUnknokwnEquipmentType = errors.New("неизвестный тип оборудования")
var ErrInsifficientFunds = errors.New("недостаточно средств")
var ErrNotAllEquipmentPurchased = errors.New("не все оборудование закуплено")
var ErrEquipmentAlreadyPurchased = errors.New("оборудование уже закуплено")
var ErrCompanyStopped = errors.New("компания остановлена")
