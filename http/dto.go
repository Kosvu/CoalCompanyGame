package http

import (
	"coalcompny/mine"
	"encoding/json"
	"errors"
	"time"
)

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type BuyEquipmentDTO struct {
	Name    string `json:"name"`
	Buy     bool   `json:"buy"`
	Balance int    `json:"balance"`
}

type EquipmentDTO struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Bought bool   `json:"bought"`
}

type FinishDTO struct {
	Success  bool             `json:"success"`
	Duration time.Duration    `json:"duration"`
	Info     mine.CompanyInfo `json:"info"`
}

type MinerDTO struct {
	Class           string `json:"class"`
	RemainingEnergy int    `json:"remainingenergy"`
}

type HireDTO struct {
	Class string `json:"class"`
}

type HireMinerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Balance int    `json:"balance,omitempty"`
}

func (h HireDTO) ValidateForCreate() error {
	if h.Class == "" {
		return errors.New("class is empty")
	}
	return nil
}

func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
