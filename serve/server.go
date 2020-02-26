package serve

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/liuqiqi-Y/qiqiCache/concurrency"
)

const (
	defaultBasePath = "/_qiqiCache/"
)

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(s string) *HTTPPool {
	return &HTTPPool{
		self:     s,
		basePath: defaultBasePath,
	}
}
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server: %s] %s\n", p.self, fmt.Sprintf(format, v...))
}
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		http.Error(w, "wrong request path", http.StatusBadRequest)
		return
	}
	p.Log("%s %s", r.Method, r.URL.Path)

	parts := strings.SplitAfterN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	//log.Printf("===>%v\n", r.URL.Path[len(p.basePath):])
	//log.Printf("%v===%v\n", parts[0], parts[1])
	group := concurrency.GetGroup(parts[0][:len(parts[0])-1]) //去除字符中最后的'/'
	if group == nil {
		http.Error(w, "have no this group: "+parts[0], http.StatusNotFound)
		return
	}
	v, err := group.Get(parts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(v.ByteSlice())
}
