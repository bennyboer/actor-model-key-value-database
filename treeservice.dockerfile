FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o bin/tree-service ./treeservice

FROM iron/go
COPY --from=builder /app/bin/tree-service /app/tree-service
EXPOSE 8090
ENTRYPOINT ["/app/tree-service"]
