package mine

import (
	"context"
	"sync"
	"time"
)

type Equipment struct {
	Name   string
	Price  int
	Bought bool
}

type CompanyInfo struct {
	Balance      int
	ActiveMiners int
	TotalHired   map[string]int
	Equipment    []EquipmentInfo
}

type EquipmentInfo struct {
	Name   string
	Price  int
	Bought bool
}

type GameResult struct {
	Info     CompanyInfo
	Duration time.Duration
}

type Company struct {
	Balance    int
	Miner      []Miner
	TotalHired map[string]int
	mu         sync.Mutex
	Equipment  []Equipment
	Сtx        context.Context
	cancel     context.CancelFunc
	StartTime  time.Time
	EndTime    time.Time
}

func NewCompany() *Company {
	ctx, cancel := context.WithCancel(context.Background())

	return &Company{
		Balance:    0,
		Miner:      []Miner{},
		TotalHired: make(map[string]int),
		Equipment: []Equipment{
			{Name: "Pickaxe", Price: 3000, Bought: false},
			{Name: "Ventilation", Price: 15000, Bought: false},
			{Name: "Cart", Price: 50000, Bought: false},
		},
		Сtx:       ctx,
		cancel:    cancel,
		StartTime: time.Now(),
	}
}

func PassiveIncome(a *Company, ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.mu.Lock()
				a.Balance++
				a.mu.Unlock()
			}
		}
	}()
}

func (c *Company) HireMiner(miner Miner) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Balance < miner.Info().Price {
		return false
	}

	c.Balance -= miner.Info().Price

	// добавляем шахтёра в список
	c.Miner = append(c.Miner, miner)
	c.TotalHired[miner.Info().MinerClass]++

	// запускаем горутину
	go func(m Miner) {
		for coal := range miner.Run(c.Сtx) {
			c.mu.Lock()
			c.Balance += coal.Amount
			c.mu.Unlock()
		}

		// удаляем шахтёра после завершения работы
		c.mu.Lock()
		newMiners := []Miner{}
		for _, active := range c.Miner {
			if active != m {
				newMiners = append(newMiners, active)
			}
		}
		c.Miner = newMiners
		c.mu.Unlock()
	}(miner)

	return true
}

func (c *Company) Info() CompanyInfo {
	c.mu.Lock()
	defer c.mu.Unlock()

	eqInfo := make([]EquipmentInfo, len(c.Equipment))
	for i, eq := range c.Equipment {
		eqInfo[i] = EquipmentInfo{
			Name:   eq.Name,
			Bought: eq.Bought,
			Price:  eq.Price,
		}
	}

	return CompanyInfo{
		Balance:      c.Balance,
		ActiveMiners: len(c.Miner),
		TotalHired:   c.TotalHired,
		Equipment:    eqInfo,
	}
}

func (c *Company) BuyEquipment(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range c.Equipment {
		if c.Equipment[i].Name == name {
			if c.Equipment[i].Bought {
				return false
			}
			if c.Balance < c.Equipment[i].Price {
				return false
			}

			c.Balance -= c.Equipment[i].Price
			c.Equipment[i].Bought = true
			return true
		}
	}

	return false
}

func (c *Company) IsAllEquipmentBought() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, eq := range c.Equipment {
		if !eq.Bought {
			return false
		}
	}

	return true
}

func (c *Company) Finish() (GameResult, error) {
	if !c.IsAllEquipmentBought() {
		return GameResult{}, ErrorGameNotComplete
	}

	c.mu.Lock()

	if !c.EndTime.IsZero() {
		duration := c.EndTime.Sub(c.StartTime)
		info := c.Info()
		c.mu.Unlock()
		return GameResult{Info: info, Duration: duration}, nil
	}

	c.EndTime = time.Now()

	if c.cancel != nil {
		c.cancel()
	}

	duration := c.EndTime.Sub(c.StartTime)
	c.mu.Unlock()

	info := c.Info()

	return GameResult{Info: info, Duration: duration}, nil
}
