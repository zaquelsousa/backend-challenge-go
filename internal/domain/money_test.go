package domain

import "testing"

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
