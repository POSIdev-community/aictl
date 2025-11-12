package metric

import (
	"runtime"

	"github.com/POSIdev-community/aictl/pkg/logger"
)

func PrintMemStat(log *logger.Logger) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.StdErr("Memory Usage:")
	log.StdErrf("  Allocated: %v MB", bToMb(m.Alloc))
	log.StdErrf("  TotalAlloc: %v MB", bToMb(m.TotalAlloc))
	log.StdErrf("  Sys: %v MB", bToMb(m.Sys))
	log.StdErrf("  NumGC: %v", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
