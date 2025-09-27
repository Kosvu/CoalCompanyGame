package mine

import (
	"context"
	"time"
)

type Coal struct {
	Amount int
}

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}

type SmallMiner struct {
	Price   int
	Energy  int
	GetCoal int
	Pause   time.Duration
}

func (s *SmallMiner) Run(ctx context.Context) <-chan Coal {
	ch := make(chan Coal)

	go func() {
		defer close(ch)
		for s.Energy > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.Pause):
				ch <- Coal{Amount: s.GetCoal}
				s.Energy--
			}
		}
	}()

	return ch
}

func (s *SmallMiner) Info() MinerInfo {
	return MinerInfo{
		MinerClass:      "Small",
		RemainingEnergy: s.Energy,
		Price:           5,
	}
}

func NewSmallMiner() *SmallMiner {
	return &SmallMiner{Price: 5, Energy: 30, GetCoal: 1, Pause: 3 * time.Second}
}

type NormalMiner struct {
	Price   int
	Energy  int
	GetCoal int
	Pause   time.Duration
}

func (n *NormalMiner) Run(ctx context.Context) <-chan Coal {
	ch := make(chan Coal)

	go func() {
		defer close(ch)
		for n.Energy > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(n.Pause):
				ch <- Coal{Amount: n.GetCoal}
				n.Energy--
			}
		}
	}()

	return ch
}

func (n *NormalMiner) Info() MinerInfo {
	return MinerInfo{
		MinerClass:      "Normal",
		RemainingEnergy: n.Energy,
		Price:           50,
	}
}

func NewNormalMiner() *NormalMiner {
	return &NormalMiner{Price: 50, Energy: 45, GetCoal: 3, Pause: 2 * time.Second}
}

type StrongMiner struct {
	Price   int
	Energy  int
	GetCoal int
	Pause   time.Duration
}

func (st *StrongMiner) Run(ctx context.Context) <-chan Coal {
	ch := make(chan Coal)

	go func() {
		defer close(ch)
		for st.Energy > 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(st.Pause):
				ch <- Coal{Amount: st.GetCoal}
				st.Energy--
			}
		}
	}()

	return ch
}

func (st *StrongMiner) Info() MinerInfo {
	return MinerInfo{
		MinerClass:      "Strong",
		RemainingEnergy: st.Energy,
		Price:           450,
	}
}

func NewStrongMiner() *StrongMiner {
	return &StrongMiner{Price: 450, Energy: 60, GetCoal: 10, Pause: 1 * time.Second}
}

type MinerInfo struct {
	MinerClass      string
	RemainingEnergy int
	Price           int
}
