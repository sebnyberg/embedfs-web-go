package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/klauspost/compress/gzhttp"
	"github.com/sebnyberg/flagtags"
	"github.com/urfave/cli/v2"
)

// Note: updates require a corresponding change to the go:embed directive below
const path = "static"

//go:embed static
var embeddedFS embed.FS

type Config struct {
	Addr             string `value:"localhost:0"`
	GzipMinSizeBytes int    `value:"1024"`
	LogText          bool   `value:"false"`
}

func main() {
	var cfg Config
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app := &cli.App{
		Usage: "Run the static file server",
		Flags: flagtags.MustParseFlags(&cfg),
		Action: func(c *cli.Context) error {
			return Run(cfg, log)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func Run(cfg Config, log *slog.Logger) error {
	if cfg.LogText {
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	// Load 'path' from the embedded filesystem stored in embeddedFs
	staticFs, err := fs.Sub(embeddedFS, path)
	if err != nil {
		return fmt.Errorf(
			"mount subpath at %s in embedded fs failed, "+
				"please ensure that both //go:embed and const path have the same value, "+
				"err: %w",
			path, err,
		)
	}

	// Open a TCP socket on --addr / $ADDR
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("listen on %s failed: %w", cfg.Addr, err)
	}

	// Create a Gzip withGzip for the file server
	withGzip, err := gzhttp.NewWrapper(gzhttp.MinSize(cfg.GzipMinSizeBytes))
	if err != nil {
		return fmt.Errorf("create gzip wrapper failed: %w", err)
	}

	// Create a File Server with Gzip compression
	srv := withGzip(http.FileServer(http.FS(staticFs)))

	// Run the server
	log.Info("running file server", "addr", lis.Addr())
	return http.Serve(lis, srv)
}
