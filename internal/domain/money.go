package domain

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Money struct {
	Amount int64
	Currency string
}

func validateCurrency(code string) error {
	if len(code) != 3 {
		return fmt.Errorf("invalid currency")
	}
	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf("invalid currency")
		}
	}
	return nil
}

func NewMoney(amountMinor int64, currency string) (Money, error) {
	if err := validateCurrency(currency); err != nil {
		return Money{}, err
	}
	return Money{Amount: amountMinor, Currency: currency}, nil
}

func ParseMoney(value string, currency string) (Money, error){
	if err := validateCurrency(currency); err != nil {
		return Money{}, err
	}

	parts := strings.Split(value, ".")
	if len(parts) > 2 || len(parts[0]) == 0{
		return Money{}, fmt.Errorf("invalid money format")
	}

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 0 {
		return Money{}, fmt.Errorf("invalid money amount")
	}

	fraction := "00"
	if len(parts) == 2 {
		fraction = parts[1]
	}

	if len(fraction) > 2  || len(fraction) == 0 {
		return Money{}, fmt.Errorf("invalid fractional amount")
	}

	if len(fraction) == 1 {
		fraction  += "0"
	}

	cents, err := strconv.ParseInt(fraction, 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("invalid fractional amount")
	}

	if whole > (math.MaxInt64-cents)/100 {
		return Money{}, fmt.Errorf("money amount overflow")
	}

	return NewMoney(whole*100+cents, currency)
}

func (m Money) Add(other Money) (Money, error){
	if m.Currency != other.Currency{
		return Money{}, fmt.Errorf("currency mismatch: %s and %s",m.Currency, other.Currency)
	}

	if (other.Amount > 0 && m.Amount > math.MaxInt64-other.Amount)||
	(other.Amount < 0 && m.Amount < math.MinInt64-other.Amount){
		return Money{}, fmt.Errorf("money addition overflow")
	}

	return NewMoney(m.Amount+other.Amount, m.Currency)
}
