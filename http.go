// 提供单机版的HTTP服务，为后续的通过HTTP通信各节点做准备

package cache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	_defaultBasePath = "/_cache/"
)

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self, // 主机名(ip):port
		basePath: _defaultBasePath,
	}
}

func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", h.self, fmt.Sprintf(format, v...))
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
