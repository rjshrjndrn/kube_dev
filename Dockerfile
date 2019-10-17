# Which base image to use to build
FROM golang as builder
# Changing to workdirectory
WORKDIR /work
# Copy the sourcecode
COPY webscrapper.go /work
# Run the build command
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix webscrapper_signals.go -o app .

# Base secure image for production app
FROM alpine
# Copy the artifact
COPY --from=builder  /work/app .
# Run application
cmd ./app
