package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/klauspost/compress/gzhttp"
)

//go:embed static
var embeddedFs embed.FS

var (
	readTimeout = flag.Int("read-timeout", 5, "Read timeout in seconds")
	idleTimeout = flag.Int("idle-timeout", 120, "Idle timeout in seconds")
	path        = flag.String("path", "", "Path to static directory")
	addr        = flag.String("addr", "localhost:0", "Address, e.g. localhost:8080")
)

func main() {
	flag.Parse()
	staticFs, err := fs.Sub(embeddedFs, *path)
	errexit(err)
	lis, err := net.Listen("tcp", *addr)
	errexit(err)
	wrapper, err := gzhttp.NewWrapper(
		gzhttp.MinSize(1),
		gzhttp.ContentTypeFilter(gzhttp.CompressAllContentTypeFilter),
	)
	errexit(err)
	srv := &http.Server{
		Handler:     wrapper(http.FileServer(http.FS(staticFs))),
		ReadTimeout: time.Duration(*readTimeout) * time.Second,
		IdleTimeout: time.Duration(*idleTimeout) * time.Second,
	}
	fmt.Printf("Serving on http://%s\n", lis.Addr())
	errexit(srv.Serve(lis))
}

func errexit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
