# embedfs-web-go

Just an EmbedFS example serving a static webpage.

## Changing the path

To change the path from `static`, update the path in both the constant
variable, and the `//go:embed` directive:

```bash
--- main.go     2024-05-29 19:06:09
+++ main2.go    2024-05-29 19:25:09
@@ -15,9 +15,9 @@
 // Note: updates require a corresponding change to the go:embed directive below
-const path = "static"
+const path = "path/to/dir"

-//go:embed static
+//go:embed path/to/dir
```

## Running locally

```bash
go run main.go --addr "localhost:8080"
```

## Running as a static binary

```bash
go build -o fileserver .
./fileserver --addr "localhost:8080"
```

## Running as a Docker image

Please note socket bind to the virtual ethernet device (`0.0.0.0`) rather than
loopback (`localhost`):

```bash
docker build . -t fileserver
docker run -it --rm -p 8080:8080 --user 65534:65534 \
  fileserver --addr 0.0.0.0:8080
```
