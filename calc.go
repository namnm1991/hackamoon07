package main

import (
	"errors"
	"math"
)

func mean(nums []float64) (float64, error) {
	if (len(nums)) == 0 {
		return 0, errors.New("empty array of nums")
	}

	var sum, mean float64
	for _, num := range nums {
		sum += num
	}
	mean = sum / float64(len(nums))

	return mean, nil
}

func sd(nums []float64) (float64, error) {
	if (len(nums)) == 0 {
		return 0, errors.New("empty array of nums")
	}

	mean, err := mean(nums)
	if err != nil {
		return 0, err
	}

	var sd float64
	for _, num := range nums {
		sd += math.Pow(num-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(nums)))

	return sd, nil
}

// var num[10] float64
// var sum,mean,sd float64
//   fmt.Println("******  Enter 10 elements  *******")
//   for i := 1; i <= 10; i++ {
// 	  fmt.Printf("Enter %d element : ",i)
// 	  fmt.Scan(&num[i-1])
// 	  sum += num[i-1]
//   }
//   mean = sum/10;

//   for j := 0; j < 10; j++ {
// 	   // The use of Pow math function func Pow(x, y float64) float64
// 	  sd += math.Pow(num[j] - mean, 2)
//   }
//   // The use of Sqrt math function func Sqrt(x float64) float64
//   sd = math.Sqrt(sd/10)

//   fmt.Println("The Standard Deviation is : ",sd)
