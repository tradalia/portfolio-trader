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
	"log/slog"
	"sync"
	"time"

	"github.com/tradalia/portfolio-trader/pkg/core"
	"github.com/tradalia/portfolio-trader/pkg/db"
)

//=============================================================================

var jobs = struct {
	sync.RWMutex
	m map[uint]*Process
}{m: make(map[uint]*Process)}

//-----------------------------------------------------------------------------

var workers = core.WorkerPool{}

//=============================================================================
//===
//=== Init
//===
//=============================================================================

func init() {
	num := 4
	workers.Init(num, 100)
	go periodicCleanup()
}

//=============================================================================
//===
//=== API methods
//===
//=============================================================================

func Start(req *Request, ts *db.TradingSystem, trades *[]db.Trade, risk float64) {
	jobs.Lock()
	defer jobs.Unlock()

	sp, ok := jobs.m[ts.Id]
	if ok {
		slog.Error("Stopping a previous simulation process", "tsId", ts.Id)
		sp.Stop()
		delete(jobs.m, ts.Id)
	}

	sp = NewProcess(ts, trades, req, risk)
	workers.Submit(sp.Start)

	jobs.m[ts.Id] = sp
}

//=============================================================================

func Stop(tsId uint) bool {
	jobs.Lock()
	defer jobs.Unlock()

	sp, ok := jobs.m[tsId]
	if ok {
		sp.Stop()
		delete(jobs.m, tsId)
	}

	return ok
}

//=============================================================================

func GetResult(tsId uint) *Result {
	jobs.Lock()
	defer jobs.Unlock()

	sp, ok := jobs.m[tsId]
	if !ok {
		return &Result{
			Status: SimStatusIdle,
		}
	}

	return sp.GetResult()
}

//=============================================================================
//===
//=== Cleanup process
//===
//=============================================================================

func periodicCleanup() {
	for {
		time.Sleep(time.Minute * 5)
		purge()
	}
}

//=============================================================================

func purge() {
	jobs.Lock()
	defer jobs.Unlock()

	for tsId, sp := range jobs.m {
		if sp.result.Status == SimStatusComplete {
			delta := time.Now().Sub(sp.result.EndTime)
			if delta.Minutes() >= 30 {
				slog.Info("purge: Purging simulation process entry for trading system", "tsId", tsId)
				delete(jobs.m, tsId)
			}
		}
	}
}

//=============================================================================
