package serve

import (
	"net/http"
	"strings"

	"github.com/liuqiqi-Y/qiqiCache/concurrency"
)

const (
	DefaultBasePath = "/_qiqiCache/"
)

type HTTPServer struct{}

var Server HTTPServer

func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, DefaultBasePath) {
		http.Error(w, "wrong request path", http.StatusBadRequest)
		return
	}

	parts := strings.SplitAfterN(r.URL.Path[len(DefaultBasePath):], "/", 2)
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
