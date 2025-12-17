package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"progect-game/company"
	"progect-game/company/equipment"
	"progect-game/company/miners"
	"progect-game/database"
	input_dto "progect-game/http/dto/input"
	output_dto "progect-game/http/dto/output"
	"time"
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
	if database.Redis != nil {
		// Ключ для кеша
		cacheKey := "company_stats"

		// Пробуем получить из Redis
		cached, err := database.Redis.Get(context.Background(), cacheKey).Result()
		if err == nil {
			// Нашли в кеше - отдаём
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cached))
			return
		} else {
			fmt.Printf("Не нашли в кеше: %v\n", err)

		}
	}

	// 2. Вычисляем данные (как обычно)
	response := map[string]interface{}{
		"balance":      h.company.GetMoney(),
		"total_earned": h.company.GetTotalEarned(),
		"miners_count": h.company.GetMinerCount(),
	}

	jsonData, _ := json.Marshal(response)

	// 3. Сохраняем в Redis если он подключен
	if database.Redis != nil {
		fmt.Println("Сохраняем в Redis...")
		ctx := context.Background()
		err := database.Redis.Set(ctx, "company_stats", jsonData, 30*time.Second).Err()
		if err != nil {
			fmt.Printf("Ошибка сохранения в Redis: %v\n", err)
		} else {
			fmt.Println("Сохранено в Redis на 30 секунд")
		}

	}

	// 4. Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

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

func (h *HTTPHeandlers) SaveGame(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем ПОЛНОЕ состояние компании
	state := h.company.GetState()

	// 2. Сохраняем ВСЁ состояние
	err := database.SaveState(state)
	if err != nil {
		http.Error(w, "Не удалось сохранить: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Отправляем ответ
	response := map[string]interface{}{
		"status": "saved",
		"data": map[string]interface{}{
			"money":        state.Money,
			"total_earned": state.TotalEarned,
			"miners_count": state.MinersCount,
			"pickaxe":      state.Pickaxe,
			"ventilation":  state.Ventilation,
			"trolleys":     state.Trolleys,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LoadGame - GET /api/load
func (h *HTTPHeandlers) LoadGame(w http.ResponseWriter, r *http.Request) {
	// 1. Загружаем состояние из БД
	state, err := database.LoadState()
	if err != nil {
		http.Error(w, "Не удалось загрузить: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Загружаем состояние в компанию
	if err := h.company.SetState(state); err != nil {
		http.Error(w, "Не удалось применить сохранение: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Отправляем ответ
	response := map[string]interface{}{
		"status": "loaded",
		"data":   state,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
