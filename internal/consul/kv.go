package consul

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var _kvClient = &http.Client{Timeout: 3 * time.Second}

const _kvPrefix = "notebookum/controller"

// KVGet fetches a single key from Consul KV under the controller prefix.
// Returns defaultVal when the key is absent or Consul is unreachable.
func KVGet(consulURL, key, defaultVal string) string {
	url := fmt.Sprintf("%s/v1/kv/%s/%s?raw", consulURL, _kvPrefix, key)
	resp, err := _kvClient.Get(url)
	if err != nil {
		log.Printf("consul_kv: GET %s failed: %v (using default)", key, err)
		return defaultVal
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode != http.StatusNotFound {
			log.Printf("consul_kv: GET %s returned HTTP %d (using default)", key, resp.StatusCode)
		}
		return defaultVal
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return defaultVal
	}
	return strings.TrimSpace(string(body))
}

// KVGetList fetches a JSON-encoded string array from Consul KV.
// Returns defaultVal when the key is absent, parse fails, or Consul is unreachable.
func KVGetList(consulURL, key string, defaultVal []string) []string {
	url := fmt.Sprintf("%s/v1/kv/%s/%s?raw", consulURL, _kvPrefix, key)
	resp, err := _kvClient.Get(url)
	if err != nil {
		return defaultVal
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return defaultVal
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return defaultVal
	}
	var result []string
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("consul_kv: key %s is not a valid JSON array (using default)", key)
		return defaultVal
	}
	return result
}
