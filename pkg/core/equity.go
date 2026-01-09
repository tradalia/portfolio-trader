//=============================================================================
/*
Copyright © 2023 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package core

import (
	"time"

	"github.com/tradalia/portfolio-trader/pkg/db"
)

//=============================================================================

func BuildGrossProfits(trades *[]db.Trade, tradeType string) (*[]time.Time, *[]float64){
	timeSlice := []time.Time{}
	equSlice  := []float64{}

	for _, tr := range *trades {
		if tradeType == db.TradeTypeAll || tr.TradeType == tradeType {
			etime := tr.ExitDate
			gross := tr.GrossProfit

			timeSlice = append(timeSlice, *etime)
			equSlice  = append(equSlice ,  gross)
		}
	}

	return &timeSlice, &equSlice
}

//=============================================================================

func BuildNetProfits(grossProfits *[]float64, costPerOper float64) *[]float64 {
	netSlice := []float64{}

	for _, gross := range *grossProfits {
		net := gross - 2 * costPerOper
		netSlice = append(netSlice, net)
	}

	return &netSlice
}

//=============================================================================

func BuildEquity(profits *[]float64) *[]float64 {
	equity := []float64{}
	value  := 0.0

	for _, profit := range *profits {
		value += profit

		equity = append(equity, value)
	}

	return &equity
}

//=============================================================================

func BuildDrawDown(equity *[]float64) (*[]float64, float64) {
	maxProfit    := 0.0
	currDrawDown := 0.0
	maxDrawDown  := 0.0
	drawDown     := []float64{}

	for _, currProfit := range *equity {
		if currProfit >= maxProfit {
			maxProfit = currProfit
			currDrawDown = 0
		} else {
			currDrawDown = currProfit - maxProfit
		}

		drawDown = append(drawDown, currDrawDown)

		if currDrawDown < maxDrawDown {
			maxDrawDown = currDrawDown
		}
	}

	return &drawDown, maxDrawDown
}

//=============================================================================

func CalcWinningPercentage(profits []float64, filter []int8) float64 {
	tot := 0
	pos := 0

	for i, profit := range profits {
		if profit != 0 {
			if filter == nil || filter[i] == 1 {
				tot++
				if profit > 0 {
					pos++
				}
			}
		}
	}

	if tot == 0 {
		return 0
	}

	return float64(pos * 10000 / tot) / 100
}

//=============================================================================

func CalcAverageTrade(profits []float64, filter []int8) float64 {
	sum := 0.0
	num := 0.0

	for i, profit := range profits {
		if profit != 0 {
			if filter == nil || filter[i] == 1 {
				sum += profit
				num++
			}
		}
	}

	return float64(int(sum * 100 / num)) / 100
}

//=============================================================================

func CalcMin(data []float64) float64 {
	minv := data[0]
	for _, value := range data {
		if value < minv {
			minv = value
		}
	}

	return minv
}

//=============================================================================
