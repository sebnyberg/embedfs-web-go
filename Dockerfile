FROM cgr.dev/chainguard/go AS builder

COPY . /app
RUN cd /app && go build -ldflags "-s -w" -o fileserver .

FROM cgr.dev/chainguard/glibc-dynamic

COPY --from=builder /app/fileserver /usr/local/bin/fileserver

ENTRYPOINT ["/usr/local/bin/fileserver"]