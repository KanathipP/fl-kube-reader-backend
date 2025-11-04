package service

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PodViewingService struct {
	ClientSet kubernetes.Interface
}

func NewPodViewingService(cs kubernetes.Interface) *PodViewingService {
	return &PodViewingService{cs}
}

func NewPodViewingServiceFromConfig(cfg *rest.Config) (*PodViewingService, error) {
	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("create clientset : %w", err)
	}

	return &PodViewingService{cs}, nil
}

type PodSummary struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	NodeName   string            `json:"nodename"`
	Phase      string            `json:"phase"`
	Labels     map[string]string `json:"labels"`
	Component  string            `json:"component"`
	Instance   string            `json:"instance"`
	Containers []string          `json:"containers"`
}

func (s *PodViewingService) GetPods(ctx context.Context, namespace string, baseSelector string, component string, instance string) ([]PodSummary, error) {
	var parts []string

	if baseSelector != "" {
		parts = append(parts, baseSelector)
	}

	if component != "" {
		parts = append(parts, fmt.Sprintf("app.kubernetes.io/component=%s", component))
	}

	if instance != "" {
		parts = append(parts, fmt.Sprintf("app.kubernetes.io/instance=%s", instance))
	}

	selector := strings.Join(parts, ",")

	pods, err := s.ClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, fmt.Errorf("get pods error %w", err)
	}

	var out []PodSummary

	for _, p := range pods.Items {
		lbls := p.Labels
		if lbls == nil {
			lbls = map[string]string{}
		}
		var containers []string
		for _, c := range p.Spec.Containers {
			containers = append(containers, c.Name)
		}

		out = append(out, PodSummary{
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

func (s *PodViewingService) GetPodLog(ctx context.Context, namespace, pod, container string, tail int, previous bool) (string, error) {
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

	req := s.ClientSet.CoreV1().Pods(namespace).GetLogs(pod, opts)

	data, err := req.DoRaw(ctx)
	if err != nil {
		return "", fmt.Errorf("get logs for pod %s: %w", pod, err)
	}

	return string(data), nil
}
