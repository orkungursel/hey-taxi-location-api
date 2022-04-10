FROM golang:1.18-buster AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags build -ldflags="-s -w" -trimpath -o /bin/app cmd/main.go

FROM gcr.io/distroless/static
COPY --from=build /bin/app /bin/app

ENV ACTIVE_PROFILE=production
EXPOSE 8080

ENTRYPOINT ["/bin/app"]