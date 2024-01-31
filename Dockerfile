FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /cl-inter

FROM build AS run-test
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build /cl-inter /cl-inter

EXPOSE 8081

CMD ["/cl-inter"]



