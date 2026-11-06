package cache

import (
	"context"

	"github.com/cxd-scu/bbot/biz/model"
	"github.com/dgraph-io/ristretto"
)

var ristrettoCache *ristretto.Cache

func InitRistretto(ctx context.Context) {
	var err error
	ristrettoCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 100,
		MaxCost:     100,
		BufferItems: 16, // 缓冲池大小，并发高时推荐 64
		Metrics:     true,
	})
	if err != nil {
		panic(err)
	}
}

func Set(k string, v any) {
	ristrettoCache.Set(k, v, 1)
	ristrettoCache.Wait()
}

func Get[T any](k string) (T, bool) {
	val, ok := ristrettoCache.Get(k)
	if !ok {
		var zeroT T
		return zeroT, false
	}

	return val.(T), true
}

func GetOrCreateSession(sessionID string) *model.Session {
	s, ok := Get[*model.Session](sessionID)
	if !ok {
		s = &model.Session{
			SessionID:    sessionID,
			CheckPointID: "",
		}
		Set(sessionID, s)
	}
	return s
}

func SetSession(session *model.Session) {
	Set(session.SessionID, session)
}
