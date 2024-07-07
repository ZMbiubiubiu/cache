// 提供单机版的HTTP服务，为后续的通过HTTP通信各节点做准备

package cache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"cache/consistenthash"
)

const (
	_defaultBasePath = "/_cache/"
	_defaultReplicas = 50
)

type HTTPPool struct {
	selfAddr    string
	basePath    string
	mu          sync.Mutex // 保护下面字段
	peers       *consistenthash.ConsistentHash
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		selfAddr: self, // 主机名(ip):port
		basePath: _defaultBasePath,
	}
}

func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", h.selfAddr, fmt.Sprintf(format, v...))
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log("path=%s", r.URL.Path)
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// path:/<basepath>/<groupName>/<key>
	parts := strings.SplitN(r.URL.Path[len(_defaultBasePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName, key := parts[0], parts[1]

	group := GetGroup(groupName)
	if group == nil {
		h.Log("group(%s) not found", groupName)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	v, err := group.Get(key)
	if err != nil {
		h.Log("get failed, err=%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(v.ByteSlice())
}

func (h *HTTPPool) Set(peerAddrs ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistenthash.NewConsistentHash(_defaultReplicas, nil)
	h.peers.AddNodes(peerAddrs...)
	h.httpGetters = make(map[string]*httpGetter, len(peerAddrs))
	for _, addr := range peerAddrs {
		h.httpGetters[addr] = &httpGetter{baseURL: addr + h.basePath}
	}
}

func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if addr := h.peers.GetNode(key); addr != "" && addr != h.selfAddr {
		h.Log("Pick peer %s", addr)
		return h.httpGetters[addr], true
	}
	return nil, false
}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) PeerGet(group string, key string) ([]byte, error) {
	uri := fmt.Sprintf("%s%s/%s", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))

	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}
