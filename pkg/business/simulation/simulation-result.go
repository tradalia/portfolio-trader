//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

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

package simulation

import (
	"time"

	"github.com/tradalia/core/datatype"
)

//=============================================================================

type Result struct {
	FirstTradeDate datatype.IntDate  `json:"firstTradeDate"`
	LastTradeDate  datatype.IntDate  `json:"lastTradeDate"`
	Runs           int               `json:"runs"`
	InitialCapital float64           `json:"initialCapital"`
	RuinPercentage int               `json:"ruinPercentage"`
	Risk           float64           `json:"risk"`

	Status         string            `json:"status"`
	StartTime      time.Time         `json:"startTime"`
	EndTime        time.Time         `json:"endTime"`
	Step           int               `json:"step"`

	GrossAll       *Details          `json:"grossAll"`
	GrossLong      *Details          `json:"grossLong"`
	GrossShort     *Details          `json:"grossShort"`
	NetAll         *Details          `json:"netAll"`
	NetLong        *Details          `json:"netLong"`
	NetShort       *Details          `json:"netShort"`
}

//=============================================================================

func NewResult(first, last datatype.IntDate, runs int, initialCapital float64, ruinPerc int, risk float64) *Result {
	return &Result{
		FirstTradeDate: first,
		LastTradeDate : last,
		Runs          : runs,
		InitialCapital: initialCapital,
		RuinPercentage: ruinPerc,
		Risk          : risk,
	}
}

//=============================================================================

type Details struct {
	Equities     string        `json:"equities"`
	MaxDrawdowns *Distribution `json:"maxDrawdowns"`
}

//=============================================================================

type Distribution struct {
	XAxis []string  `json:"xAxis"`
	YAxis []float64 `json:"yAxis"`
}

//=============================================================================
