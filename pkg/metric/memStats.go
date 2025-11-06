package metric

import (
	"runtime"

	"github.com/POSIdev-community/aictl/pkg/logger"
)

func PrintMemStat(log *logger.Logger) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.StdErr("Memory Usage:")
	log.StdErrF("  Allocated: %v MB", bToMb(m.Alloc))
	log.StdErrF("  TotalAlloc: %v MB", bToMb(m.TotalAlloc))
	log.StdErrF("  Sys: %v MB", bToMb(m.Sys))
	log.StdErrF("  NumGC: %v", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
