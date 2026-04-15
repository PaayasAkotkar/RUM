package rum

import (
	"fmt"
	"sync"
)

// Budget tracks spend against a limit for a model service
type Budget struct {
	mu    sync.Mutex
	Limit float64
	Spent float64
	Cost  float64
	perc  map[int]float64 // call count -> percent
}

// NewBudget creates a new budget with a limit and cost per call
func NewBudget(limit, cost float64) *Budget {
	return &Budget{
		Limit: limit,
		Cost:  cost,
		perc:  make(map[int]float64),
	}
}

// get funcs

func (b *Budget) GetCost() float64 {
	return b.Cost
}
func (b *Budget) GetSpend() float64 {
	return b.Spent
}
func (b *Budget) GetLimit() float64 {
	return b.Limit
}

// percent progress of the budget
func (b *Budget) percent() float64 {
	if b.Limit <= 0 {
		return 0
	}
	return (b.Spent / b.Limit) * 100
}

// Exhausted returns true when spend >= limit
func (b *Budget) Exhausted() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.percent() >= 100.0
}

// Left returns the remaining budget
// do use the percent func cause i am still figuring this func out
func (b *Budget) Left() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Limit - b.Spent
}

// Report returns report of the budget
func (b *Budget) Report() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return fmt.Sprintf("$%.4f / $%.4f (%.1f%%)", b.Spent, b.Limit, b.percent())
}

// end

// set funcs

// Spend adds one cost to the spend
func (b *Budget) Spend() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Spent += b.Cost

}

// Percent calcs the percent of the budget and stores it
func (b *Budget) Percent() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	calc := b.percent()
	count := len(b.perc) + 1
	if _, ok := b.perc[count]; !ok {
		b.perc[count] = calc
	}
	return calc
}

// end
