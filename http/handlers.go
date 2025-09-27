package http

import (
	"coalcompny/mine"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	company         *mine.Company
	closeServerFunc func() error
}

func NewHTTPHandlers(company *mine.Company) *HTTPHandlers {
	return &HTTPHandlers{
		company: company,
	}
}

func (h *HTTPHandlers) SetCloseServerFunc(f func() error) {
	h.closeServerFunc = f
}

/*
pattern: /company
method: GET
info: -

succeed:
	-status code: 200 OK
	-response body: JSON represented found tasks

failed:
	-status code: 400,500 ...
	-response body: JSON with error + time
*/

func (h *HTTPHandlers) HandleGetCompanyInfo(w http.ResponseWriter, r *http.Request) {
	companyinfo := h.company.Info()

	b, err := json.MarshalIndent(companyinfo, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http responce:", err)
		return
	}
}

/*
pattern: /company
method: POST
info: JSON in HTTP request body

suceed:
	- status code: 201 created
	- response body: JSON represent created task

failed:
	-status code: 400,409,500 ...
	-responce body: JSON with error + time
*/

func (h *HTTPHandlers) HandleFinishGame(w http.ResponseWriter, r *http.Request) {
	result, err := h.company.Finish()

	if err != nil {
		if errors.Is(err, mine.ErrorGameNotComplete) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorDTO{
				Message: "Not all equipment is bought",
				Time:    time.Now(),
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FinishDTO{
		Success:  true,
		Duration: result.Duration,
		Info:     result.Info,
	})

	if h.closeServerFunc != nil {
		go h.closeServerFunc()
	}
}

/*
pattern: /miners
method: GET
info: -

succeed:
	-status code: 200 OK
	-response body: JSON represented found tasks

failed:
	-status code: ---
	-response body: ---
*/

func (h *HTTPHandlers) HandleGetMiners(w http.ResponseWriter, r *http.Request) {
	infos := []MinerDTO{}
	for _, m := range h.company.Miner {
		infos = append(infos, MinerDTO{
			Class:           m.Info().MinerClass,
			RemainingEnergy: m.Info().RemainingEnergy,
		})
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(infos); err != nil {
		fmt.Println("failed to write miners:", err)
		return
	}
}

/*
pattern: /miners/{class}
method: GET
info: pattern

succed:
	-status code: 200 OK
	-respsonse body: JSON represented found task

failed:
	- status code: 400, 404, 500 ...
	-response body: JSON with error + time
*/

func (h *HTTPHandlers) HandleGetMinersByClass(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["class"]
	infos := []MinerDTO{}

	for _, m := range h.company.Miner {
		if m.Info().MinerClass == title {
			infos = append(infos, MinerDTO{
				Class:           m.Info().MinerClass,
				RemainingEnergy: m.Info().RemainingEnergy,
			})
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(infos); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}

}

/*
pattern: /miners
method: POST
info: JSON in HTTP request body

suceed:
  - status code: 201 created
  - response body: JSON represent created task

failed:

	-status code: 400,409,500 ...
	-responce body: JSON HireMinerResponse
*/
func (h *HTTPHandlers) HandleHireMiner(w http.ResponseWriter, r *http.Request) {
	var hireDTO HireDTO
	if err := json.NewDecoder(r.Body).Decode(&hireDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := hireDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	var miner mine.Miner

	switch hireDTO.Class {
	case "Small":
		miner = mine.NewSmallMiner()
	case "Normal":
		miner = mine.NewNormalMiner()
	case "Strong":
		miner = mine.NewStrongMiner()
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HireMinerResponse{
			Success: false,
			Message: "Hire error",
			Balance: h.company.Info().Balance,
		})

		return
	}
	ok := h.company.HireMiner(miner)
	if !ok {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(HireMinerResponse{
			Success: false,
			Message: "Not enough balance",
			Balance: h.company.Info().Balance,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(HireMinerResponse{
		Success: true,
		Message: "Successful",
		Balance: h.company.Info().Balance,
	}); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /equipment
method: GET
info: ---

succeed:
	status code: 200 OK
	response: JSON represent created task
*/

func (h *HTTPHandlers) HandleGetEquipment(w http.ResponseWriter, r *http.Request) {
	companyInfo := h.company.Info()
	arr := []EquipmentDTO{}

	for _, eq := range companyInfo.Equipment {
		arr = append(arr, EquipmentDTO{
			Name:   eq.Name,
			Price:  eq.Price,
			Bought: eq.Bought,
		})
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(arr); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /equipment/{name}
method: POST {бизнес действие}
info: ---

succeed:
	status code: 200 OK
	response: JSON with equipment name

failed:
	status code: 409 {conflict}
*/

func (h *HTTPHandlers) HandleBuyEquipment(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	buy := h.company.BuyEquipment(name)

	if !buy {

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorDTO{
			Message: "Not enough money to buy",
			Time:    time.Now(),
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BuyEquipmentDTO{
		Name:    name,
		Buy:     true,
		Balance: h.company.Info().Balance,
	})

}
