// @title           fl-kube-reader-backend API
// @version         1.0
// @description     Kubernetes Pod & Node metrics viewer backend
// @BasePath        /api/v1
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"reflect"

	"github.com/KanathipP/fl-kube-reader-backend/pkg/config"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/handlers"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/routes"
	"github.com/KanathipP/fl-kube-reader-backend/pkg/service"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {
	os.Setenv("TZ", "Asia/Bangkok")

	godotenvEnv := config.LoadDotEnvConfig()
	slog.Info("ENV has been loaded")

	kubeRestCfg, err := config.LoadKubeRestConfig()
	if err != nil {
		log.Fatalf("cannot create K8s config: %v", err)
	}

	cs, err := kubernetes.NewForConfig(kubeRestCfg)
	if err != nil {
		log.Fatalf("cannot create K8s ClientSet: %v", err)
	}
	mcs, err := metrics.NewForConfig(kubeRestCfg)
	if err != nil {
		log.Fatalf("cannot create K8s MetricsClientSet: %v", err)
	}

	podViewingService := service.NewPodViewingService(cs, godotenvEnv.KubeConfig)
	podViewingHandler := handlers.NewPodViewingHandler(podViewingService)
	metricsViewingService := service.NewMetricsViewingService(cs, mcs, godotenvEnv.KubeConfig)
	metricsViewingHandler := handlers.NewMetricsViewingHandler(metricsViewingService)
	logReadingService := service.NewLogReadingService(cs, godotenvEnv.KubeConfig)
	logReadingHandler := handlers.NewLogReadingHandler(logReadingService)

	app := fiber.New(fiber.Config{
		JSONDecoder: func(b []byte, v any) error {
			dec := json.NewDecoder(bytes.NewReader(b))
			dec.DisallowUnknownFields()
			if err := dec.Decode(v); err != nil {
				return fmt.Errorf("decode: %w", err)
			}
			if err := dec.Decode(new(struct{})); err != io.EOF {
				return fmt.Errorf("decode: trailing data")
			}

			rv := reflect.ValueOf(v)
			for rv.Kind() == reflect.Pointer {
				rv = rv.Elem()
			}
			validate := validator.New()
			if rv.Kind() == reflect.Struct {
				if err := validate.Struct(v); err != nil {
					return err
				}
			}
			return nil
		},
	})

	routes.SetupRoutes(app, podViewingHandler, metricsViewingHandler, logReadingHandler)
	port := godotenvEnv.ServerConfig.Port
	fmt.Println("Server is running on port " + port)
	app.Listen("localhost:" + port)
}
