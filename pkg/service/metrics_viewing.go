package service

import (
	"context"
	"fmt"

	"github.com/KanathipP/fl-kube-reader-backend/pkg/config"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/dtos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type MetricsViewingService struct {
	ClientSet        kubernetes.Interface
	MetricsClientSet metrics.Interface
	KubeConfig       config.KubeConfig
}

func NewMetricsViewingService(cs kubernetes.Interface, mcs metrics.Interface, kubeConfig config.KubeConfig) *MetricsViewingService {
	return &MetricsViewingService{cs, mcs, kubeConfig}
}

func NewMetricsViewingServiceFromConfig(cfg *rest.Config, kubeConfig config.KubeConfig) (*MetricsViewingService, error) {
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create clientset : %w", err)
	}

	mcs, err := metrics.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create metric clientset : %w", err)
	}
	return &MetricsViewingService{cs, mcs, kubeConfig}, nil
}

func (s *MetricsViewingService) GetPodUsage(ctx context.Context, pod, container string) (dtos.GetPodUsageDto, error) {
	pm, err := s.MetricsClientSet.MetricsV1beta1().PodMetricses(s.KubeConfig.Namespace).Get(ctx, pod, metav1.GetOptions{})
	if err != nil {
		return dtos.GetPodUsageDto{}, fmt.Errorf("get pod metrics: %w", err)
	}

	var cpuTotalMilli int64
	var memTotal int64

	for _, c := range pm.Containers {
		if container != "" && c.Name != container {
			continue
		}
		cpuTotalMilli += c.Usage.Cpu().MilliValue()
		memTotal += c.Usage.Memory().Value()
	}

	return dtos.GetPodUsageDto{
		Namespace:   s.KubeConfig.Namespace,
		Pod:         pod,
		Container:   container,
		CPU_m:       cpuTotalMilli,
		MemoryBytes: memTotal,
	}, nil
}

func (s *MetricsViewingService) GetNodesUsage(ctx context.Context) ([]dtos.GetNodeUsageDto, error) {
	// metrics
	nm, err := s.MetricsClientSet.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list node metrics error: %w", err)
	}
	// allocatable
	nodes, err := s.ClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list nodes error: %w", err)
	}

	alloc := make(map[string]struct {
		cpuMilli int64
		mem      int64
	})
	for _, n := range nodes.Items {
		alloc[n.Name] = struct {
			cpuMilli int64
			mem      int64
		}{
			cpuMilli: n.Status.Allocatable.Cpu().MilliValue(),
			mem:      n.Status.Allocatable.Memory().Value(),
		}
	}

	var out []dtos.GetNodeUsageDto

	for _, m := range nm.Items {
		a := alloc[m.Name]
		out = append(out, dtos.GetNodeUsageDto{
			Node:          m.Name,
			CPU_m:         m.Usage.Cpu().MilliValue(),
			MemoryBytes:   m.Usage.Memory().Value(),
			AllocCPU_m:    a.cpuMilli,
			AllocMemBytes: a.mem,
		})
	}
	return out, nil
}
