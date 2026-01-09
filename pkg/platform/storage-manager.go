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
	"log/slog"
	"strconv"

	"github.com/tradalia/core/auth"
	"github.com/tradalia/core/req"
)

//=============================================================================

type EquityRequest struct {
	Username string            `json:"username"`
	Images   map[string][]byte `json:"images"`
}
//-----------------------------------------------------------------------------

func NewEquityRequest() *EquityRequest {
	return &EquityRequest{
		Images: map[string][]byte{},
	}
}

//=============================================================================
//===
//=== Public functions
//===
//=============================================================================

func SetEquityChart(id uint, er *EquityRequest) error {
	slog.Info("SetEquityChart: Sending equity chart to storage manager", "id", id, "username", er.Username)

	token,err := auth.Token()
	if err != nil {
		return err
	}

	client :=req.GetClient("bf")
	url    := platform.Storage +"/v1/trading-systems/"+ strconv.Itoa(int(id)) +"/equity-chart"

	err = req.DoPut(client, url, &er, "", token)
	if err != nil {
		slog.Error("SetEquityChart: Got an error when sending to storage-manager", "id", id, "error", err.Error())
		return req.NewServerError("Cannot communicate with storage-manager: %v", err.Error())
	}

	slog.Info("SetEquityChart: Equity chart saved", "id", id)
	return nil
}

//=============================================================================

func DeleteEquityChart(username string, id uint) error {
	slog.Info("DeleteEquityChart: Deleting equity chart from the storage manager", "id", id, "username", username)

	token,err := auth.Token()
	if err != nil {
		return err
	}

	client :=req.GetClient("bf")
	url    := platform.Storage +"/v1/trading-systems/"+ strconv.Itoa(int(id)) +"/equity-chart"
	er     := EquityRequest{
		Username: username,
	}

	err = req.DoDelete(client, url, &er, "", token)
	if err != nil {
		slog.Error("DeleteEquityChart: Got an error when sending to storage-manager", "id", id, "error", err.Error())
		return req.NewServerError("Cannot communicate with storage-manager: %v", err.Error())
	}

	slog.Info("DeleteEquityChart: Equity chart deleted", "id", id, "username", username)
	return nil
}

//=============================================================================
