package dtos

type GetPodUsageDto struct {
	Namespace     string  `json:"namespace"`
	Pod           string  `json:"pod"`
	Container     string  `json:"container,omitempty"`
	CPU_m         int64   `json:"cpu_m"`                    // millicores
	MemoryBytes   int64   `json:"memory_bytes"`             // bytes
	CPUPercent    float64 `json:"cpu_percent,omitempty"`    // optional (ถ้ามี limit/allocatable)
	MemoryPercent float64 `json:"memory_percent,omitempty"` // optional
}
