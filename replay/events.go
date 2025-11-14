package replay

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// ParseReplaySpeed parses replay speed strings like "1x", "10x", "0.5x"
func ParseReplaySpeed(speed string) (int, error) {
	speed = strings.TrimSpace(strings.ToLower(speed))
	speed = strings.TrimSuffix(speed, "x")

	multiplier, err := strconv.ParseFloat(speed, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid speed format: %w", err)
	}
	if multiplier <= 0 {
		return 0, fmt.Errorf("speed multiplier must be positive")
	}
	return int(multiplier), nil
}

func StreamLiveEvents(kubeconfig string, speedMultiplier int) error {
	// Resolve kubeconfig path with proper fallbacks
	if kubeconfig == "" {
		if envPath := os.Getenv("KUBECONFIG"); envPath != "" {
			kubeconfig = envPath
		} else if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	watch, err := clientset.CoreV1().Events("").Watch(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watch.Stop()

	for ev := range watch.ResultChan() {
		e, ok := ev.Object.(*v1.Event)
		if !ok {
			continue
		}

		// Build kubectl command
		kind := strings.ToLower(e.InvolvedObject.Kind)
		name := e.InvolvedObject.Name
		namespace := e.InvolvedObject.Namespace

		var cmd string
		if namespace == "" {
			// Cluster-scoped resource (e.g., Node, PersistentVolume)
			cmd = fmt.Sprintf("kubectl describe %s %s", kind, name)
		} else {
			// Namespaced resource
			cmd = fmt.Sprintf("kubectl describe %s %s -n %s", kind, name, namespace)
		}

		// Add event metadata as comment
		cmd = fmt.Sprintf("%s # reason=%s message=%s", cmd, e.Reason, e.Message)

		fmt.Println(cmd)
		time.Sleep(time.Second / time.Duration(speedMultiplier))
	}
	return nil
}
