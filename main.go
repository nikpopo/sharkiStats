package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

func req(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(body))

	return string(body)
}

func main() {

	var loanProfit float64
	var totalLend float64
	var loancnt int
	var secondcnt int
	var fulfilledD int
	var d14 int
	var d7 int
	var d16 int
	var maxtime int
	var mintime int
	var maxtimeindex int64
	var mintimeindex int64

	maxtime = -1
	mintime = 1e9

	for i := 0; i < 500; i += 50 {
		temp := req("https://sharky.fi/api/loan/my-loans?lender=wallet_name&network=mainnet&deployEnvironment=production&offset=" + strconv.Itoa(i))

		for j := 0; j < 50; j++ {
			principalLamports := gjson.Get(temp, "historyLoans." + strconv.Itoa(j+1) + ".principalLamports")
			amountRepaidLamports := gjson.Get(temp, "historyLoans." + strconv.Itoa(j+1) + ".amountRepaidLamports")

			loantime := gjson.Get(temp, "historyLoans." + strconv.Itoa(j+1) + ".durationSeconds")

			if loantime.Int() ==  604800 {
				d7++
			}
			if loantime.Int() == 1209600 {
				d14++
			}
			if loantime.Int() == 1382400 {
				d16++
			}

			totalLend += principalLamports.Float()

			if principalLamports.Int() != 0 {
				loancnt++

				dateTaken := gjson.Get(temp, "historyLoans." + strconv.Itoa(j+1) + ".dateTaken")
				dateRepaid := gjson.Get(temp, "historyLoans." + strconv.Itoa(j+1) + ".dateRepaid")

				dateTakenS := dateTaken.String()
				dateRepaidS := dateRepaid.String()

				var dayT int
				var dayR int
				var hourT int
				var hourR int
				var minT int
				var minR int
				var secT int
				var secR int
				var secTotal int

				if len(dateTakenS) != 0 && len(dateRepaidS) != 0 {
					//days
					dayT, _ = strconv.Atoi(string(dateTakenS[8]))
					dayT = dayT * 10
					tempo, _ := strconv.Atoi(string(dateTakenS[9]))
					dayT = dayT + tempo

					dayR, _ = strconv.Atoi(string(dateRepaidS[8]))
					dayR = dayR * 10
					tempo, _ = strconv.Atoi(string(dateRepaidS[9]))
					dayR = dayR + tempo

					//hours
					hourT, _ = strconv.Atoi(string(dateTakenS[11]))
					hourT = hourT * 10
					tempo, _ = strconv.Atoi(string(dateTakenS[12]))
					hourT = hourT + tempo

					hourR, _ = strconv.Atoi(string(dateRepaidS[11]))
					hourR = hourR * 10
					tempo, _ = strconv.Atoi(string(dateRepaidS[12]))
					hourR = hourR + tempo

					//mins
					minT, _ = strconv.Atoi(string(dateTakenS[14]))
					minT = minT * 10
					tempo, _ = strconv.Atoi(string(dateTakenS[15]))
					minT = minT + tempo

					minR, _ = strconv.Atoi(string(dateRepaidS[14]))
					minR = minR * 10
					tempo, _ = strconv.Atoi(string(dateRepaidS[15]))
					minR = minR + tempo

					//secs
					secT, _ = strconv.Atoi(string(dateTakenS[17]))
					secT = secT * 10
					tempo, _ = strconv.Atoi(string(dateTakenS[18]))
					secT = secT + tempo

					secR, _ = strconv.Atoi(string(dateRepaidS[17]))
					secR = secR * 10
					tempo, _ = strconv.Atoi(string(dateRepaidS[18]))
					secR = secR + tempo

					if dateTakenS[5] == dateRepaidS[5] && dateTakenS[6] == dateRepaidS[6] {
						//работаем в рамках одного месяца

						var tTotal int
						var rTotal int

						tTotal = 24 * 60 * 60 * dayT + 60 * 60 * hourT + 60 * minT + secT
						rTotal = 24 * 60 * 60 * dayR + 60 * 60 * hourR + 60 * minR + secR

						secTotal = rTotal - tTotal

					} else {
						var monthT int
						var days int

						monthT, _ = strconv.Atoi(string(dateTakenS[5]))
						monthT = monthT * 10
						tempo, _ = strconv.Atoi(string(dateTakenS[6]))
						monthT = monthT + tempo

						if monthT == 1 || monthT == 3 || monthT == 5 || monthT == 7 || monthT == 8 || monthT == 10 || monthT == 12 {
							days = 31 - dayT + dayR

							var tTotal int
							var rTotal int

							tTotal = 60 * 60 * hourT + 60 * minT + secT
							rTotal = 60 * 60 * hourR + 60 * minR + secR

							secTotal = rTotal - tTotal + 24 * 60 * 60 * days
						}
						var year int

						year, _ = strconv.Atoi(string(dateTakenS[0]))
						year = year * 1000
						tempo, _ = strconv.Atoi(string(dateTakenS[1]))
						year = year + 100 * tempo
						tempo, _ = strconv.Atoi(string(dateTakenS[2]))
						year = year + 10 * tempo
						tempo, _ = strconv.Atoi(string(dateTakenS[3]))
						year = year + tempo

						if monthT == 2 {
							if year % 4 == 0 {
								days = 29 - dayT + dayR

								var tTotal int
								var rTotal int

								tTotal = 60 * 60 * hourT + 60 * minT + secT
								rTotal = 60 * 60 * hourT + 60 * minR + secR

								secTotal = rTotal - tTotal + 24 * 60 * 60 * days
							} else {
								days = 28 - dayT + dayR

								var tTotal int
								var rTotal int

								tTotal = 60 * 60 * hourT + 60 * minT + secT
								rTotal = 60 * 60 * hourR + 60 * minR + secR

								secTotal = rTotal - tTotal + 24 * 60 * 60 * days
							}
						}

						if monthT == 4 || monthT == 6 || monthT == 9 || monthT == 11 {
							days = 30 - dayT + dayR

							var tTotal int
							var rTotal int

							tTotal = 60 * 60 * hourT + 60 * minT + secT
							rTotal = 60 * 60 * hourR + 60 * minR + secR

							secTotal = rTotal - tTotal + 24 * 60 * 60 * days
						}
					}

					fulfilledD++
					secondcnt += secTotal

					if secTotal >= maxtime {
						maxtime = secTotal
						maxtimeindex = loantime.Int()
					}
					if secTotal <= mintime {
						mintime = secTotal
						mintimeindex = loantime.Int()
					}
				}
			}

			if amountRepaidLamports.Int() - principalLamports.Int() > 0 {
				loanProfit += amountRepaidLamports.Float() - principalLamports.Float()
			} else {

			}
		}
	}

	var secondavg int
	secondavg = (secondcnt + ((loancnt - fulfilledD) * 7 * 24 * 60 * 60)) / loancnt

	var daycnt int
	daycnt = secondavg / (24*60*60)

	var hourcnt int
	hourcnt = (secondavg - (daycnt * 24 * 60 * 60)) / (60 * 60)

	var mincnt int
	mincnt = (secondavg - (daycnt * 24 * 60 * 60) - (hourcnt * 60 * 60)) / 60

	var seccnt int
	seccnt = secondavg - (daycnt * 24 * 60 * 60) - (hourcnt * 60 * 60) - mincnt * 60

	var secondavg1 int
	secondavg1 = secondcnt / loancnt

	var daycnt1 int
	daycnt1 = secondavg1 / (24*60*60)

	var hourcnt1 int
	hourcnt1 = (secondavg1 - (daycnt1 * 24 * 60 * 60)) / (60 * 60)

	var mincnt1 int
	mincnt1 = (secondavg1 - (daycnt1 * 24 * 60 * 60) - (hourcnt1 * 60 * 60)) / 60

	var seccnt1 int
	seccnt1 = secondavg1 - (daycnt1 * 24 * 60 * 60) - (hourcnt1 * 60 * 60) - mincnt1 * 60

	fmt.Println("Total proceed loans counter:\n", d7, "for 7 days\n", d14, "for 14 days\n", d16, "for 16 days\n", loancnt, "total")
	fmt.Print("Defaulted loans: ", loancnt - fulfilledD, " / ", math.Round(float64(loancnt - fulfilledD) / float64(loancnt) * 100 * 1000) / 1000)
	fmt.Println("%")
	fmt.Println("Avg loan time(defaulted loans included):", daycnt, "Days", hourcnt, "Hours", mincnt, "Minutes", seccnt, "Seconds")
	fmt.Println("Avg loan time(defaulted loans excluded):", daycnt1, "Days", hourcnt1, "Hours", mincnt1, "Minutes", seccnt1, "Seconds")

	fmt.Print("Total profit from loans, excluding defaulted ones: ")
	fmt.Printf("%.2f", loanProfit / 1e9)
	fmt.Println(" SOL")

	fmt.Print("Total accepted & repaid offers: ")
	fmt.Printf("%.2f", totalLend / 1e9)
	fmt.Println(" SOL")

	fmt.Print("Avg percent: ")
	fmt.Printf("%.2f", loanProfit / totalLend * 100)
	fmt.Println(" %")

	fmt.Println("Avg profit per loan:", math.Round((loanProfit / 1e9) / float64(loancnt) * 100) / 100, "SOL")

	fmt.Print("Max time loan: ", maxtime, " of ", maxtimeindex, " seconds, which is ", math.Round(float64(maxtime) / float64(maxtimeindex) * 100 * 1000) / 1000)
	fmt.Println("%")
	fmt.Print("Min time loan: ", mintime, " of ", mintimeindex, " seconds, which is ", math.Round(float64(mintime) / float64(mintimeindex) * 100 * 1000) / 1000)
	fmt.Println("%")
}
