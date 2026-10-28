package tipDataCache

import (
	"github.com/Jagtech-Services/core-go-lib/milieu"
	"github.com/Jagtech-Services/go-tari-grpc-lib/v3/tari_generated"
	"github.com/Jagtech-Services/go-tari-lib/v3/nodeGRPC"
	"sync"
)

type tipDataStruct struct {
	tipResponse *tari_generated.TipInfoResponse
	mutex       sync.RWMutex
}

var running = false

var tipData *tipDataStruct = nil

func UpdateTipData(core *milieu.Milieu) {
	running = true
	defer func() {
		running = false
	}()
	tipResponse, err := nodeGRPC.GetTipInfo()
	if err != nil {
		core.Debug(err.Error())
		core.CaptureException(err)
		return
	}
	if tipData == nil {
		core.Debug("Initializing the tip cache")
		tipData = &tipDataStruct{}
	}
	tipData.mutex.Lock()
	tipData.tipResponse = tipResponse
	tipData.mutex.Unlock()
}

func GetTipData() *tari_generated.TipInfoResponse {
	if tipData == nil {
		return nil
	}
	tipData.mutex.RLock()
	defer tipData.mutex.RUnlock()
	return tipData.tipResponse
}
