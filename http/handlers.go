package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"progect_game/company"
	"progect_game/company/equipment"
	"progect_game/company/miners"
	input_dto "progect_game/http/dto/input"
	output_dto "progect_game/http/dto/output"
)

type HTTPHeandlers struct {
	company     *company.Company
	closeServer func() error
}

func NewHTTPHeandlers(company *company.Company) *HTTPHeandlers {
	return &HTTPHeandlers{
		company: company,
	}
}

func (heandler *HTTPHeandlers) SetCloseServerFunc(f func() error) {
	heandler.closeServer = f
}

// создание нового майнера
func (h *HTTPHeandlers) HandleCreateNewMiner(w http.ResponseWriter, r *http.Request) {
	var minerDTO input_dto.MinerDTO

	if err := json.NewDecoder(r.Body).Decode(&minerDTO); err != nil {
		http.Error(w, NewError(err).ToString(), http.StatusBadRequest)
		return
	}
	minerType := miners.MinerClass(minerDTO.MinerType)

	miner, err := h.company.HireMiner(minerType)
	if err != nil {
		errDTO := NewError(err)

		if errors.Is(err, company.ErrUnknokwnMinerType) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else if errors.Is(err, company.ErrInsifficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(miner.Info(), "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}

}

// вывод информации о нанятых майнерах
func (h *HTTPHeandlers) HandlerGetMiner(w http.ResponseWriter, r *http.Request) {
	minerType := r.URL.Query().Get("type")
	if minerType != "" {
		minersByType := h.company.GetMinersByType(miners.MinerClass(minerType))
		outputDTO := output_dto.NewMinersByTypeDTO(minersByType)

		b, err := json.MarshalIndent(outputDTO, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		if _, err := w.Write(b); err != nil {
			http.Error(w, NewError(err).ToString(), http.StatusInternalServerError)
		}
		return
	}
	miners := h.company.GetAllMiners()
	outputDTO := output_dto.NewAllMinerDTO(miners)

	b, err := json.MarshalIndent(outputDTO, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(b); err != nil {
		http.Error(w, NewError(err).ToString(), http.StatusInternalServerError)
	}
}

// покупка оборудования
func (h *HTTPHeandlers) HandleByeEquipment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var equipmentInputDTO input_dto.EquipmentDTO

	err := json.NewDecoder(r.Body).Decode(&equipmentInputDTO)
	if err != nil {
		http.Error(w, NewError(err).ToString(), http.StatusBadRequest)
		return
	}
	equipmentType := equipment.EquipmentType(equipmentInputDTO.EquipmentType)
	equipment, err := h.company.BuyEquipment(equipmentType)
	if err != nil {
		errDTO := NewError(err)

		if errors.Is(err, company.ErrUnknokwnEquipmentType) {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		} else if errors.Is(err, company.ErrInsifficientFunds) {
			http.Error(w, errDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		}
		return
	}
	equipmentOutputDTO := output_dto.NewEquipmentDTO(equipment)
	b, err := json.MarshalIndent(equipmentOutputDTO, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// проверка купленного оборудования
func (h *HTTPHeandlers) HandleCheckEquipment(w http.ResponseWriter, r *http.Request) {
	minerType := h.company.GetEquipment()
	minerDTO := output_dto.NewEquipmentDTO(minerType)

	b, err := json.MarshalIndent(minerDTO, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// статистика компании
func (h *HTTPHeandlers) HandleGetCompanyStatistics(w http.ResponseWriter, r *http.Request) {

	if h.company == nil {
		log.Printf("ERROR: company is nil in HTTPHandlers")
		http.Error(w, `{"error": "Company not initialized"}`, http.StatusInternalServerError)
		return
	}
	stats := h.company.GetStatistics()

	statsDTO := output_dto.NewCompanyStatisticDTO(stats)

	b, err := json.MarshalIndent(statsDTO, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// check о зарплатах майнеров
func (h *HTTPHeandlers) HandleGetMinerSalaries(w http.ResponseWriter, r *http.Request) {
	salariesOutputDto := output_dto.NewminersSalariesDTO(
		miners.LittleMinerSalary,
		miners.BasicMinerSalary,
		miners.PowerfulMinerSalary,
	)
	b, err := json.MarshalIndent(salariesOutputDto, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// инфо о стоимости оборудования
func (h *HTTPHeandlers) HandleGetEquipmentPrice(w http.ResponseWriter, r *http.Request) {
	priceOutputDto := output_dto.NewEquipmntPriceDTO(
		equipment.EquipmentPickaxeConst,
		equipment.EquipmentTrolleysConst,
		equipment.EquipmentVentilationConst,
	)
	b, err := json.MarshalIndent(priceOutputDto, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}
}

// хендлер для завершения игры если условия игры выполнены
func (h *HTTPHeandlers) HandleCompleateGame(w http.ResponseWriter, r *http.Request) {
	stats, err := h.company.FinishGame()
	if err != nil {
		errorDTO := NewError(err)
		if errors.Is(err, company.ErrEquipmentAlreadyPurchased) {
			http.Error(w, errorDTO.ToString(), http.StatusForbidden)
		} else {
			http.Error(w, errorDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	statisticOutputDTO := output_dto.NewCompanyStatisticDTO(*stats)

	b, err := json.MarshalIndent(statisticOutputDTO, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response:", err)
	}

	go func() {
		if err := h.closeServer(); err != nil {
			fmt.Println("failed to close HTTP server:", err)
		}
	}()
}
