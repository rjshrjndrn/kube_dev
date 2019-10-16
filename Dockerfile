FROM golang as builder
WORKDIR /work
COPY webscrapper_signals.go /work
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix webscrapper_signals.go -o app .

FROM alpine
COPY --from=builder  /work/app .
cmd ./app
