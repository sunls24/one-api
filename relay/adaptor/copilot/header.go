package copilot

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/songquanpeng/one-api/common/cache"
)

var salt, _ = os.LookupEnv("COPILOT_SALT")

type vscodeHeader struct {
	machineID string
	sessionID string

	lastSessionTime   int64
	updateSessionTime int64
}

func createHeader(key string) *vscodeHeader {
	sum := sha256.Sum256([]byte(key + salt))
	return &vscodeHeader{
		machineID:         hex.EncodeToString(sum[:]),
		updateSessionTime: genUpdateSessionTime(),
	}
}

func genUpdateSessionTime() int64 {
	// 5min-60min
	return int64(rand.Intn(3600) + 300)
}

var cacheHeader = cache.NewCache[*vscodeHeader](createHeader)

func GetSessionID(key string) string {
	cached, _ := cacheHeader.Get(key)
	now := time.Now().Unix()
	if now-cached.lastSessionTime > cached.updateSessionTime {
		cached.lastSessionTime = now
		cached.updateSessionTime = genUpdateSessionTime()
		cached.sessionID = uuid.New().String() + strconv.FormatInt(now*1000, 10)
	}
	return cached.sessionID
}

func GetMachineID(key string) string {
	cached, _ := cacheHeader.Get(key)
	return cached.machineID
}
