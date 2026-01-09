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

package platform

import (
	"fmt"

	"github.com/tradalia/core/auth"
	"github.com/tradalia/core/datatype"
	"github.com/tradalia/core/req"
)

//=============================================================================
//===
//=== Response
//===
//=============================================================================

const (
	DirectionStrongBear = -2
	DirectionBear       = -1
	DirectionNeutral    =  0
	DirectionBull       =  1
	DirectionStrongBull =  2
)

const (
	VolatilityQuiet        = 0
	VolatilityNormal       = 1
	VolatilityVolatile     = 2
	VolatilityVeryVolatile = 3
)

//=============================================================================

type DataProductAnalysisResponse struct {
	Id              uint             `json:"id"`
	Symbol          string           `json:"symbol"`
	From            datatype.IntDate `json:"from"`
	To              datatype.IntDate `json:"to"`
	Days            int              `json:"days"`
	DailyResults    []*DailyResult   `json:"dailyResults"`
}

//=============================================================================

type DailyResult struct {
	Date            datatype.IntDate `json:"date"`
	Price           float64          `json:"price"`
	PercDailyChange float64          `json:"percDailyChange"`
	Sqn100          float64          `json:"sqn100"`
	TrueRange       float64          `json:"trueRange"`
	PercAtr20       float64          `json:"percAtr20"`
	Direction       int              `json:"direction"`
	Volatility      int              `json:"volatility"`
}

//=============================================================================
//===
//=== Public functions
//===
//=============================================================================

func AnalyzeDataProduct(c *auth.Context, id uint, backDays int) (*DataProductAnalysisResponse, error) {
	c.Log.Info("AnalyzeDataProduct: Asking data product analysis to data collector", "id", id, "backDays", backDays)

	token  := c.Token
	client := req.GetClient("bf")
	url    := fmt.Sprintf("%s/v1/data-products/%d/analysis?backDays=%d", platform.Data, id, backDays)

	var res DataProductAnalysisResponse
	err := req.DoGet(client, url, &res, token)
	if err != nil {
		c.Log.Error("AnalyzeDataProduct: Got an error when accessing the data-collector", "id", id, "error", err.Error())
		return nil,req.NewServerError("Cannot communicate with data-manager: %v", err.Error())
	}

	c.Log.Info("AnalyzeDataProduct: Analysis received", "id", id)
	return &res,nil
}

//=============================================================================
