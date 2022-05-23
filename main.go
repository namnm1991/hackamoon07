package main

import (
	"fmt"

	"github.com/montanaflynn/stats"
)

func main() {

	// ======================================================
	// send an email
	// emails := []string{"nam@krystal.app"}
	// subject := "Welcome to Krystal SmartAlert"
	// content := fmt.Sprintf("S.O.S %d", rand.Intn(100))
	// sendEmail(emails, subject, content)

	// ======================================================
	// draw a spark line
	// sparkline()

	// ======================================================
	// simple calculation
	// [mean - k * sigma..mean + k * sigma] range
	// (sigma stands for the standard deviation),
	// where k is typically 2 (95%), 3 (99.76%),

	nums := []float64{3, 5, 9, 1, 8, 6, 58, 9, 4, 10}
	m, _ := stats.Mean(nums)
	sd, _ := stats.StandardDeviation(nums)

	fmt.Printf("mean [%.3f], standard deviation: [%.3f]\n", m, sd)

	o, _ := stats.QuartileOutliers([]float64{-1000, 1, 3, 4, 4, 6, 6, 6, 6, 7, 8, 15, 18, 100})
	fmt.Printf("%+v\n", o)
}
