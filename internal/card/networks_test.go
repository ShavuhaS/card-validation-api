package card

import (
	"reflect"
	"testing"
)

func TestNewNetworkStorage(t *testing.T) {
	storage, err := NewNetworkStorage()
	if err != nil {
		t.Fatalf("NewNetworkStorage() returned error: %v", err)
	}
	if storage == nil {
		t.Fatalf("NewNetworkStorage() returned nil storage")
	}
}

func TestNetworkStorage_FindIssuer(t *testing.T) {
	storage, err := NewNetworkStorage()
	if err != nil {
		t.Fatalf("Failed to initialize NetworkStorage: %v", err)
	}

	tests := []struct {
		name           string
		iin            string
		expectedIssuer CardIssuer
		expectedFound  bool
	}{
		{"Visa lower bound", "40000000", Visa, true},
		{"Visa middle", "41111111", Visa, true},
		{"Visa upper bound", "49999999", Visa, true},
		{"Mastercard 2221 lower bound", "22210000", Mastercard, true},
		{"Mastercard 2221 middle", "25000000", Mastercard, true},
		{"Mastercard 2720 upper bound", "27209999", Mastercard, true},
		{"Mastercard 51 lower bound", "51000000", Mastercard, true},
		{"Mastercard 55 upper bound", "55999999", Mastercard, true},
		{"Amex 34 lower bound", "34000000", AmericanExpress, true},
		{"Amex 34 upper bound", "34999999", AmericanExpress, true},
		{"Amex 37 lower bound", "37000000", AmericanExpress, true},
		{"Amex 37 upper bound", "37999999", AmericanExpress, true},
		{"Discover 6011 lower bound", "60110000", Discover, true},
		{"Discover 622126 lower bound", "62212600", Discover, true},
		{"Discover 644 lower bound", "64400000", Discover, true},
		{"Discover 65 lower bound", "65000000", Discover, true},
		{"Diners 30 lower", "30000000", DinersClubInternational, true},
		{"Diners 36 lower", "36000000", DinersClubInternational, true},
		{"Diners 38 lower", "38000000", DinersClubInternational, true},
		{"Diners 39 lower", "39000000", DinersClubInternational, true},
		{"JCB 3528 lower", "35280000", JCB, true},
		{"JCB 3589 upper", "35899999", JCB, true},
		{"UnionPay lower", "62000000", UnionPay, true},
		{"UnionPay 621", "62100000", UnionPay, true},
		{"UnionPay upper", "62999999", UnionPay, true},
		{"Unknown IIN 1", "10000000", "", false},
		{"Unknown IIN 2", "99999999", "", false},
		{"Invalid format", "4111abcd", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issuer, found := storage.FindIssuer(tt.iin)
			if found != tt.expectedFound {
				t.Errorf("FindIssuer(%q) found = %v, expected %v", tt.iin, found, tt.expectedFound)
			}
			if found && issuer != tt.expectedIssuer {
				t.Errorf("FindIssuer(%q) issuer = %v, expected %v", tt.iin, issuer, tt.expectedIssuer)
			}
		})
	}
}

func TestNetworkStorage_GetNumberLengths(t *testing.T) {
	storage, err := NewNetworkStorage()
	if err != nil {
		t.Fatalf("Failed to initialize NetworkStorage: %v", err)
	}

	tests := []struct {
		name            string
		issuer          CardIssuer
		expectedLengths []int
	}{
		{"Visa", Visa, []int{13, 16, 19}},
		{"Mastercard", Mastercard, []int{16}},
		{"AmericanExpress", AmericanExpress, []int{15}},
		{"Discover", Discover, []int{16, 17, 18, 19}},
		{"DinersClubInternational", DinersClubInternational, []int{14, 15, 16, 17, 18, 19}},
		{"JCB", JCB, []int{16, 17, 18, 19}},
		{"UnionPay", UnionPay, []int{16, 17, 18, 19}},
		{"Unknown Issuer", CardIssuer("Unknown"), []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lengths := storage.GetNumberLengths(tt.issuer)
			if len(lengths) == 0 && len(tt.expectedLengths) == 0 {
				return
			}
			if !reflect.DeepEqual(lengths, tt.expectedLengths) {
				t.Errorf("GetNumberLengths(%v) = %v, expected %v", tt.issuer, lengths, tt.expectedLengths)
			}
		})
	}
}
