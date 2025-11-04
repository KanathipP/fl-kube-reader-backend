package dtos

type GetNodeUsageDto struct {
	Node          string `json:"node"`
	CPU_m         int64  `json:"cpu_m"`        // used (sum from metrics)
	MemoryBytes   int64  `json:"memory_bytes"` // used (sum from metrics)
	AllocCPU_m    int64  `json:"alloc_cpu_m"`
	AllocMemBytes int64  `json:"alloc_mem_bytes"`
}
