package domain

import (
	"math"
	"testing"
)

func TestParseMoney(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    int64
        wantErr bool
    }{
        {"whole amount", "25", 2500, false},
        {"two decimals", "25.90", 2590, false},
        {"one decimal", "25.9", 2590, false},
        {"zero", "0.00", 0, false},
        {"invalid input", "abc", 0, true},
        {"too many decimals", "25.999", 0, true},
        {"negative amount", "-25.00", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseMoney(tt.input, "BRL")

            if (err != nil) != tt.wantErr {
                t.Fatalf("ParseMoney() error = %v, wantErr %v", err, tt.wantErr)
            }

            if err == nil && got.Amount != tt.want {
                t.Errorf("Amount = %d, want %d", got.Amount, tt.want)
            }
        })
    }
}

func TestMoneyAdd(t *testing.T) {
    a, _ := NewMoney(2500, "BRL")
    b, _ := NewMoney(750, "BRL")

    got, err := a.Add(b)
    if err != nil {
        t.Fatal(err)
    }

    if got.Amount != 3250 {
        t.Errorf("Amount = %d, want 3250", got.Amount)
    }
}

func TestMoneyAddCurrencyMismatch(t *testing.T) {
    brl, _ := NewMoney(1000, "BRL")
    usd, _ := NewMoney(1000, "USD")

    if _, err := brl.Add(usd); err == nil {
        t.Fatal("expected currency mismatch error")
    }
}

func TestMoneyIsNegative(t *testing.T){
	tests := []struct {
		name   string
		amount int64
		want   bool
	}{
		{"negative amount", -100, true},
		{"zero amount", 0, false},
		{"positive amount", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.amount, "BRL")
			if err != nil {
				t.Fatal(err)
			}
			
			if got := m.IsNegative(); got != tt.want {
				t.Errorf("IsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneySubtract(t *testing.T){
	tests := []struct {
		name      string
		a         int64
		b         int64
		want      int64
		wantErr   bool
	}{
		{"normal subtraction", 2500, 750, 1750, false},
		{"subtract to zero", 1000, 1000, 0, false},
		{"insufficient amount", 500, 750, 0, true},
		{"subtract zero", 1000, 0, 1000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewMoney(tt.a, "BRL")
			if err != nil {
				t.Fatal(err)
			}

			b, err := NewMoney(tt.b, "BRL")
			if err != nil {
				t.Fatal(err)
			}

			got, err := a.Subtract(b)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Subtract() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && got.Amount != tt.want {
				t.Errorf("Amount = %d, want %d", got.Amount, tt.want)
			}

			if err == nil && got.Currency != "BRL" {
				t.Errorf("Currency = %s, want BRL", got.Currency)
			}
		})
	}
}

func TestMoneySubtractCurrencyMismatch(t *testing.T) {
	brl, err := NewMoney(2500, "BRL")
	if err != nil {
		t.Fatal(err)
	}

	usd, err := NewMoney(750, "USD")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := brl.Subtract(usd); err == nil {
		t.Fatal("expected currency mismatch error")
	}
}

func TestMoneySubtractOverflow(t *testing.T) {
	a, err := NewMoney(math.MinInt64, "BRL")
	if err != nil {
		t.Fatal(err)
	}

	b, err := NewMoney(1, "BRL")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := a.Subtract(b); err == nil {
		t.Fatal("expected subtraction overflow error")
	}
}
