package handlers

import (
	"log/slog"

	"github.com/KanathipP/fl-kube-reader-backend/pkg/response"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/service"
	"github.com/gofiber/fiber/v2"
)

type PodViewingHandler struct {
	podViewingService *service.PodViewingService
}

func NewPodViewingHandler(podViewingService *service.PodViewingService) *PodViewingHandler {
	return &PodViewingHandler{
		podViewingService: podViewingService,
	}
}

// GetPods godoc
// @Summary     List pods
// @Description Returns pods filtered by namespace, component, and instance
// @Tags        pods
// @Accept      json
// @Produce     json
// @Param       namespace  query     string  false  "Namespace (default: default)"
// @Param       component  query     string  false  "K8s component label: app.kubernetes.io/component"
// @Param       instance   query     string  false  "K8s instance label: app.kubernetes.io/instance"
// @Success     200        {array}   map[string]interface{}  "List of pods"
// @Failure     500        {object}  response.ErrorResponse  "Failed to retrieve pods"
// @Router      /pod-viewing [get]
func (h *PodViewingHandler) GetPods(ctx *fiber.Ctx) error {
	component := ctx.Query("component", "")
	instance := ctx.Query("instance", "")

	slog.Info("GetPods",
		"component", component,
		"instance", instance,
	)

	pods, err := h.podViewingService.GetPods(
		ctx.Context(),
		component,
		instance,
	)
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get pods: "+err.Error())
	}

	return response.OK(ctx, pods)
}
