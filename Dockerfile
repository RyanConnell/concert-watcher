FROM golang:1.20 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/concert-watch cmd/watcher/main.go

FROM alpine
COPY --from=build /bin/concert-watch /bin/concert-watch
CMD [ "/bin/concert-watch" ]
