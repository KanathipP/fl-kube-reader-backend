package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/KanathipP/fl-kube-reader-backend/pkg/config"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/dtos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PodViewingService struct {
	ClientSet  kubernetes.Interface
	kubeConfig config.KubeConfig
}

func NewPodViewingService(cs kubernetes.Interface, kubeConfig config.KubeConfig) *PodViewingService {
	return &PodViewingService{cs, kubeConfig}
}

func NewPodViewingServiceFromConfig(cfg *rest.Config, kubeConfig config.KubeConfig) (*PodViewingService, error) {
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create clientset : %w", err)
	}

	return &PodViewingService{cs, kubeConfig}, nil
}

func (s *PodViewingService) GetPods(ctx context.Context, component string, instance string) ([]dtos.GetPodDto, error) {
	var parts []string

	parts = append(parts, s.kubeConfig.LabelSelector)

	if component != "" {
		parts = append(parts, fmt.Sprintf("app.kubernetes.io/component=%s", component))
	}

	if instance != "" {
		parts = append(parts, fmt.Sprintf("app.kubernetes.io/instance=%s", instance))
	}

	selector := strings.Join(parts, ",")

	pods, err := s.ClientSet.CoreV1().Pods(s.kubeConfig.Namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, fmt.Errorf("get pods error %w", err)
	}

	var out []dtos.GetPodDto

	for _, p := range pods.Items {
		lbls := p.Labels
		if lbls == nil {
			lbls = map[string]string{}
		}
		var containers []string
		for _, c := range p.Spec.Containers {
			containers = append(containers, c.Name)
		}

		out = append(out, dtos.GetPodDto{
			Name:       p.Name,
			Namespace:  p.Namespace,
			NodeName:   p.Spec.NodeName,
			Phase:      string(p.Status.Phase),
			Labels:     lbls,
			Component:  lbls["app.kubernetes.io/component"],
			Instance:   lbls["app.kubernetes.io/instance"],
			Containers: containers,
		})
	}
	return out, nil
}
