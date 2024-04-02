package handlers
import (
	"fmt"
	"net/http"
)
func root (w http.ResponseWriter, r *http.Request) {
	route := r.PathValue("route")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(fmt.Sprintf("Hello, World! %v", route)))
	fmt.Fprintf(w, "Hello print, World!",route)
}

func RootHandler(rootRouter *http.ServeMux) {
	rootRouter.HandleFunc("/{route}", root)
	rootRouter.HandleFunc("/", root)


}