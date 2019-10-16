FROM golang as builder
COPY webscrapper_signals.go /work
WORKDIR /work
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix webscrapper_signals.go -o app .

FROM alpine
COPY --from=builder  /work/app .
cmd /work/app
