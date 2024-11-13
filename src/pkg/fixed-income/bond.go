package fixedincome

import (
	"errors"
	"math"
	"time"
)

// computeN returns the amount of periods to input into bond pricing formula
func computeN(issueDate string, maturityDate string, frequency int) (int, error) {
	var idate time.Time
	var err error

	if issueDate == "today" || issueDate == "now" {
		idate = time.Now()
	} else {
		idate, err = time.Parse("2006-01-02", issueDate)
		if err != nil {
			return 0, errors.New("invalid issue date format; expected YYYY-MM-DD")
		}
	}

	mdate, err := time.Parse("2006-01-02", maturityDate)
	if err != nil {
		return 0, errors.New("invalid maturity date format; expected YYYY-MM-DD")
	}

	if mdate.Before(idate) {
		return 0, errors.New("maturity date must be after issue date")
	}

	nbYears := mdate.Year() - idate.Year()
	if mdate.YearDay() < idate.YearDay() {
		nbYears--
	}

	var N int
	// annual payments
	if frequency == 1 {
		N = nbYears
	} else { // other payment frequencies
		N = nbYears * frequency
	}

	return N, nil
}

// ComputeBondPrice returns the clean price of a bond depending on either its defined periods or issue and maturity date
func ComputeBondPrice(n int, cRate float64, cPayment float64, yield float64, faceValue float64, iDate string, mDate string, sDate string, frequency int, dirtyPrice bool, dcc string) (float64, error) {
	var err error

	if n == 0 {
		n, err = computeN(iDate, mDate, frequency)
		if err != nil {
			return 0, err
		}
	}

	var couponPayment float64
	if cPayment == 0.0 {
		couponPayment = faceValue * cRate / float64(frequency)
	} else {
		couponPayment = cPayment
	}

	periodicYield := yield / float64(frequency)

	price := applyBondFormula(n, couponPayment, periodicYield, faceValue)

	if dirtyPrice {
		accruedInterest, err := computeAccruedInterest(sDate, iDate, frequency, couponPayment, dcc)
		if err != nil {
			return 0, err
		}
		price += accruedInterest
	}

	return price, nil

}

// computeAccruedInterest calculates accrued interest based on day-count convention
func computeAccruedInterest(settlementDateStr, issueDateStr string, frequency int, couponPayment float64, dayCountConvention string) (float64, error) {

	settlementDate, err := time.Parse("2006-01-02", settlementDateStr)
	if err != nil {
		return 0, errors.New("invalid settlement date format; expected YYYY-MM-DD")
	}

	var issueDate time.Time
	if issueDateStr == "now" || issueDateStr == "today" {
		issueDate = time.Now()
	} else {
		issueDate, err = time.Parse("2006-01-02", issueDateStr)
		if err != nil {
			return 0, errors.New("invalid issue date format; expected YYYY-MM-DD")
		}
	}

	// Calculate the last coupon date
	lastCouponDate := calculateLastCouponDate(issueDate, settlementDate, frequency)
	if settlementDate.Before(lastCouponDate) {
		return 0, errors.New("settlement date is before the last coupon date")
	}

	daysBetween := settlementDate.Sub(lastCouponDate).Hours() / 24

	var daysInPeriod float64
	switch dayCountConvention {
	case "30/360":
		daysInPeriod = 360 / float64(frequency)
	case "Actual/Actual":
		nextCouponDate := lastCouponDate.AddDate(0, 12/frequency, 0)
		daysInPeriod = nextCouponDate.Sub(lastCouponDate).Hours() / 24
	case "Actual/360":
		daysInPeriod = 360 / float64(frequency)
	default:
		return 0, errors.New("unsupported day-count convention")
	}

	accruedInterest := (daysBetween / daysInPeriod) * couponPayment

	return accruedInterest, nil
}

// calculateLastCouponDate calculates the last coupon date before the settlement date
func calculateLastCouponDate(issueDate, settlementDate time.Time, frequency int) time.Time {
	periodMonths := 12 / frequency
	lastCouponDate := issueDate

	for lastCouponDate.AddDate(0, periodMonths, 0).Before(settlementDate) {
		lastCouponDate = lastCouponDate.AddDate(0, periodMonths, 0)
	}

	return lastCouponDate
}

// applyBondFormula returns the value of a bond, takes as input
// n number of periods until maturity
// c coupon rate
// r discount rate
// f face value of the bond
func applyBondFormula(n int, c float64, yield float64, f float64) float64 {
	price := 0.0

	// PV of coupon payments
	for i := 1; i <= n; i++ {
		price += c / (math.Pow((1 + yield), float64(i)))
	}

	// PV of face value
	price += f / (math.Pow((1 + yield), float64(n)))

	return price
}
