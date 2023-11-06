package log

import (
	"strings"
	"sync"

	gr "github.com/awesome-fc/golang-runtime"
	"github.com/gin-gonic/gin"
)

var m = make(map[string]struct{})
var lock = new(sync.RWMutex)

// AddRid mark request with requestID is resolving
func AddRid(rid string) {
	lock.Lock()
	defer lock.Unlock()

	if rid != "" {
		m[rid] = struct{}{}
	}
}

// RemoveRid mark request with requestID is finished
func RemoveRid(rid string) {
	lock.Lock()
	defer lock.Unlock()

	delete(m, rid)
}

// GetRid get all requestID are resolving
func GetRid() string {
	lock.RLock()
	defer lock.RUnlock()

	keys := make([]string, len(m))
	for k := range m {
		if k != "" {
			keys = append(keys, k)
		}
	}

	return strings.Join(keys, ",")
}

// RequestIDMiddleware collects requestID from request header
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := gr.NewFromContext(c.Request)

		AddRid(ctx.RequestID)
		defer RemoveRid(ctx.RequestID)

		c.Next()
	}
}
