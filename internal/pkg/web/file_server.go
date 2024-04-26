package web

import (
	"fmt"
	"log"
	"net/http"
)

// Serve files from the output csv directory on a given port.
//
// Note: This is a blocking call.
func Serve(dir string, port int) error {
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	addr := fmt.Sprintf(":%d", port)

	log.Printf("Serving files on %s", addr)
	return http.ListenAndServe(addr, nil)
}
