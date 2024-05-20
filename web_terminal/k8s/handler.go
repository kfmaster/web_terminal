package k8s

import (
	"encoding/json"
	"net/http"
)

func RegisterHandlers(client *Client) {
	http.HandleFunc("/api/namespaces", handleNamespaces(client))
	http.HandleFunc("/api/pods", handlePods(client))
	http.HandleFunc("/api/containers", handleContainers(client))
}

func handleNamespaces(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespaces, err := c.ListNamespaces()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"namespaces": namespaces})
	}
}

func handlePods(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		if namespace == "" {
			http.Error(w, "Namespace is required", http.StatusBadRequest)
			return
		}

		pods, err := c.ListPods(namespace)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"pods": pods})
	}
}

func handleContainers(c *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		pod := r.URL.Query().Get("pod")
		if namespace == "" || pod == "" {
			http.Error(w, "Namespace and pod are required", http.StatusBadRequest)
			return
		}

		containers, err := c.ListContainers(namespace, pod)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{"containers": containers})
	}
}
