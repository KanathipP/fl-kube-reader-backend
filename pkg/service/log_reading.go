package service

import (
	"context"
	"fmt"

	"github.com/KanathipP/fl-kube-reader-backend/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type LogReadingService struct {
	ClientSet  kubernetes.Interface
	kubeConfig config.KubeConfig
}

func NewLogReadingService(cs kubernetes.Interface, kubeConfig config.KubeConfig) *LogReadingService {
	return &LogReadingService{cs, kubeConfig}
}

func NewLogReadingServiceFromConfig(cfg *rest.Config, kubeConfig config.KubeConfig) (*LogReadingService, error) {
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create clientset : %w", err)
	}

	return &LogReadingService{cs, kubeConfig}, nil
}

func (s *LogReadingService) ReadPodLog(ctx context.Context, pod, container string, tail int, previous bool) (string, error) {
	opts := &corev1.PodLogOptions{
		Previous: previous,
	}

	if container != "" {
		opts.Container = container
	}

	if tail > 0 {
		t := int64(tail)
		opts.TailLines = &t
	}

	req := s.ClientSet.CoreV1().Pods(s.kubeConfig.Namespace).GetLogs(pod, opts)

	data, err := req.DoRaw(ctx)
	if err != nil {
		return "", fmt.Errorf("get logs for pod %s: %w", pod, err)
	}

	return string(data), nil
}
