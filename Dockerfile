FROM golang AS build
COPY . .
RUN CGO_ENABLED=0 go build .

FROM alpine:3.24
COPY --from=build /go/youtube-rss-server .
EXPOSE 8080
ENTRYPOINT ["./youtube-rss-server"]
