FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o bin/tree-cli ./treecli

FROM iron/go
COPY --from=builder /app/bin/tree-cli /app/tree-cli
EXPOSE 8091
ENTRYPOINT [ "/app/tree-cli" ]
