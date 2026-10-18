package handler

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// MetricsResponse — ответ с метриками приложения.
type MetricsResponse struct {
	Uptime       string         `json:"uptime"`
	GoVersion    string         `json:"go_version"`
	NumGoroutine int            `json:"num_goroutine"`
	NumCPU       int            `json:"num_cpu"`
	GOMAXPROCS   int            `json:"gomaxprocs"`
	Memory       MemoryMetrics  `json:"memory"`
	GC           GCMetrics      `json:"gc"`
	Process      ProcessMetrics `json:"process"`
}

// MemoryMetrics — метрики потребления памяти.
type MemoryMetrics struct {
	AllocMB      float64 `json:"alloc_mb"`       // Текущее выделение (heap)
	TotalAllocMB float64 `json:"total_alloc_mb"` // Суммарно выделено за всё время
	SysMB        float64 `json:"sys_mb"`         // Общее потребление от ОС
	HeapAllocMB  float64 `json:"heap_alloc_mb"`
	HeapSysMB    float64 `json:"heap_sys_mb"`
	HeapObjects  uint64  `json:"heap_objects"`    // Количество живых объектов в heap
	StackInUseMB float64 `json:"stack_in_use_mb"`
}

// GCMetrics — метрики сборщика мусора.
type GCMetrics struct {
	NumGC        uint32  `json:"num_gc"`         // Количество запусков GC
	PauseTotalMs float64 `json:"pause_total_ms"` // Суммарные паузы GC
	LastPauseMs  float64 `json:"last_pause_ms"`  // Последняя пауза GC
}

// ProcessMetrics — метрики процесса.
type ProcessMetrics struct {
	PID int `json:"pid"`
}

// GetMetrics godoc
// @Summary      Метрики приложения
// @Description  Возвращает потребление ресурсов: RAM, горутины, GC, uptime
// @Tags         Monitoring
// @Produce      json
// @Success      200  {object}  MetricsResponse
// @Router       /metrics [get]
func GetMetrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	bytesToMB := func(b uint64) float64 {
		return float64(b) / 1024 / 1024
	}

	resp := MetricsResponse{
		Uptime:       time.Since(startTime).Round(time.Second).String(),
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		GOMAXPROCS:   runtime.GOMAXPROCS(0),
		Memory: MemoryMetrics{
			AllocMB:      bytesToMB(m.Alloc),
			TotalAllocMB: bytesToMB(m.TotalAlloc),
			SysMB:        bytesToMB(m.Sys),
			HeapAllocMB:  bytesToMB(m.HeapAlloc),
			HeapSysMB:    bytesToMB(m.HeapSys),
			HeapObjects:  m.HeapObjects,
			StackInUseMB: bytesToMB(m.StackInuse),
		},
		GC: GCMetrics{
			NumGC:        m.NumGC,
			PauseTotalMs: float64(m.PauseTotalNs) / 1e6,
			LastPauseMs:  float64(m.PauseNs[(m.NumGC+255)%256]) / 1e6,
		},
		Process: ProcessMetrics{
			PID: os.Getpid(),
		},
	}

	c.JSON(http.StatusOK, resp)
}
