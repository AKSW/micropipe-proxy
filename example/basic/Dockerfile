FROM busybox

COPY . /app

HEALTHCHECK --interval=5s --timeout=1s \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./proxy"]
