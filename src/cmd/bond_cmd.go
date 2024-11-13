package cmd

import (
	"errors"
	"fmt"
	fi "go-pricer/pkg/fixed-income"

	"github.com/spf13/cobra"
)

var (
	n             int // N number of periods until maturity
	couponRate    float64
	couponPayment float64
	yield         float64
	faceValue     float64
	idate         string
	mdate         string
	frequency     int
	sdate         string
	dcc           string // day-count convention
	dp            bool   // dirty-price
	// verbose       bool
)

// bondCmd represents the ping command
var bondCmd = &cobra.Command{
	Use:   "bond",
	Short: "This computes the price of a bond",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ValidateBondFlags(); err != nil {
			fmt.Printf("error invalid arguments: %s", err.Error())
			return
		}

		price, err := fi.ComputeBondPrice(n, couponRate, couponPayment, yield, faceValue, idate, mdate, sdate, frequency, dp, dcc)
		if err != nil {
			fmt.Printf("error when computing bond price: %s", err.Error())
			return
		}

		fmt.Printf("%.2f$\n", price)
	},
}

// ValidateBondFlags checks for flag logic
func ValidateBondFlags() error {
	if yield == 0.0 {
		return errors.New("a yield value is required to compute the bond price")
	}
	if faceValue == 0.0 {
		return errors.New("face value is required to compute the bond price")
	}
	if couponRate != 0.0 && couponPayment != 0.0 {
		return errors.New("please specify only one of --coupon-rate or --coupon-payment")
	}
	if couponRate == 0.0 && couponPayment == 0.0 {
		return errors.New("a coupon of sorts, rate or payment, needs to be specified")
	}
	if dp && sdate == "" {
		return errors.New("settlement date is required when computing dirty price")
	}

	if n == 0 {
		if idate == "" && mdate == "" {
			return errors.New("either the number of periods (--periods) or both --issue-date and --maturity-date are required")
		}
		if idate == "" || mdate == "" {
			return errors.New("both --issue-date and --maturity-date must be provided together if --periods is not specified")
		}
	}

	return nil
}

func init() {
	bondCmd.Flags().IntVarP(&n, "periods", "n", 0, "Number of periods until maturity (computed if dates are given)")
	bondCmd.Flags().Float64VarP(&couponRate, "coupon-rate", "c", 0.0, "Coupon rate as a decimal (e.g., 0.05 for 5%)")
	bondCmd.Flags().Float64VarP(&couponPayment, "coupon-payment", "C", 0.0, "Coupon as a payment)")
	bondCmd.Flags().Float64VarP(&yield, "yield", "y", 0.0, "Yield to maturity as a decimal (e.g., 0.03 for 3%)")
	bondCmd.Flags().Float64VarP(&faceValue, "face-value", "F", 0.0, "Face value of the bond, typically 1000")
	bondCmd.Flags().StringVarP(&idate, "issue-date", "i", "", "Bond's original issue date (format: YYYY-MM-DD), \"now\" or \"today\" are acceptable formats")
	bondCmd.Flags().StringVarP(&mdate, "maturity-date", "m", "", "Bond's maturity date (format: YYYY-MM-DD)")
	bondCmd.Flags().IntVarP(&frequency, "frequency", "f", 2, "Payment frequency: 1 for annual, 2 for semiannual, 4 for quarterly, 12 for monthly (default is semiannual)")
	bondCmd.Flags().StringVarP(&sdate, "settlement-date", "s", "", "Date for pricing calculation (defaults to today if omitted)")
	bondCmd.Flags().StringVarP(&dcc, "day-count-convention", "d", "30/360", "Day-count convention (default: 30/360)")
	bondCmd.Flags().BoolVarP(&dp, "dirty-price", "D", false, "Compute bond's \"dirty price\" (bond price + accrued interest)")
	// bondCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	if err := bondCmd.MarkFlagRequired("yield"); err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	if err := bondCmd.MarkFlagRequired("face-value"); err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	addBondCommand(rootCmd)
}

// addBondCommand adds the command to the root command
func addBondCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(bondCmd)
}
