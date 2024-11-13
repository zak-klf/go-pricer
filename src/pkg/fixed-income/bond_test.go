package fixedincome

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestComputeN(t *testing.T) {
	tests := []struct {
		name         string
		issueDate    string
		maturityDate string
		frequency    int
		expectedN    int
		expectError  bool
	}{
		{"Valid Semiannual", "2020-01-01", "2030-01-01", 2, 20, false},
		{"Valid Annual", "2020-01-01", "2030-01-01", 1, 10, false},
		{"Invalid Date Format", "invalid-date", "2030-01-01", 2, 0, true},
		{"Maturity Before Issue", "2030-01-01", "2020-01-01", 2, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := computeN(tt.issueDate, tt.maturityDate, tt.frequency)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedN, n)
			}
		})
	}
}

func TestComputeBondPrice(t *testing.T) {
	tests := []struct {
		name          string
		n             int
		cRate         float64
		cPayment      float64
		yield         float64
		faceValue     float64
		iDate         string
		mDate         string
		sDate         string
		frequency     int
		dirtyPrice    bool
		dcc           string
		expectedPrice float64
		expectError   bool
	}{
		{"Clean Price - Semiannual", 0, 0.05, 0, 0.04, 1000, "2020-01-01", "2030-01-01", "", 2, false, "30/360", 1081.76, false},
		{"Dirty Price - Semiannual", 0, 0.05, 0, 0.04, 1000, "2020-01-01", "2030-01-01", "2024-12-01", 2, true, "30/360", 1102.87, false},
		{"Invalid Date Format", 0, 0.05, 0, 0.04, 1000, "invalid", "2030-01-01", "", 2, false, "30/360", 0, true},
		{"Missing Settlement for Dirty", 0, 0.05, 0, 0.04, 1000, "2020-01-01", "2030-01-01", "", 2, true, "30/360", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := ComputeBondPrice(tt.n, tt.cRate, tt.cPayment, tt.yield, tt.faceValue, tt.iDate, tt.mDate, tt.sDate, tt.frequency, tt.dirtyPrice, tt.dcc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InEpsilon(t, tt.expectedPrice, price, 0.01)
			}
		})
	}
}

func TestComputeAccruedInterest(t *testing.T) {
	tests := []struct {
		name            string
		settlementDate  string
		issueDate       string
		frequency       int
		couponPayment   float64
		dcc             string
		expectedAccrued float64
		expectError     bool
	}{
		{"30/360 Convention", "2024-12-01", "2020-01-01", 2, 25.0, "30/360", 21.25, false},
		{"Actual/Actual Convention", "2024-12-01", "2020-01-01", 2, 25.0, "Actual/Actual", 20.79, false},
		{"Unsupported DCC", "2024-12-01", "2020-01-01", 2, 25.0, "Unsupported", 0, true},
		{"Invalid Date Format", "invalid-date", "2020-01-01", 2, 25.0, "30/360", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accrued, err := computeAccruedInterest(tt.settlementDate, tt.issueDate, tt.frequency, tt.couponPayment, tt.dcc)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InEpsilon(t, tt.expectedAccrued, accrued, 0.01)
			}
		})
	}
}

func TestCalculateLastCouponDate(t *testing.T) {
	issueDate, _ := time.Parse("2006-01-02", "2020-01-01")
	settlementDate, _ := time.Parse("2006-01-02", "2024-12-01")
	expectedDate, _ := time.Parse("2006-01-02", "2024-07-01")

	t.Run("Calculate Last Coupon Date", func(t *testing.T) {
		lastCouponDate := calculateLastCouponDate(issueDate, settlementDate, 2)
		assert.Equal(t, expectedDate, lastCouponDate)
	})
}
