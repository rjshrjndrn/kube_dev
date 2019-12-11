## Commands
1. build docker image
    ` docker build -f Dockerfile <hub_username>/<application_name>:<version> .`
2. Run a image
    `docker run <imagename>`
    
## Shutting down gracefully
1. Handle sigterm
    1. rjshrjndrn/demo:nosigterm
    2. rjshrjndrn/demo:sigterm
2. Liveness Probe
3. readines probe
4. grace kill time

## livenessProbe
```
readinessProbe:
      httpGet:
        path: /start/healthz
        port: 8080
        httpHeaders:
        - name: Custom-Header
          value: Awesome
      initialDelaySeconds: 3
      periodSeconds: 3
```

## livenessProbe
```
livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
        httpHeaders:
        - name: Custom-Header
          value: Awesome
      initialDelaySeconds: 3
      periodSeconds: 3
```

## How a pod get terminated

1. A SIGTERM signal is sent to the main process (PID 1) in each container, and a “grace period” countdown starts (defaults to 30 seconds - see below to change it).
2. Upon the receival of the SIGTERM, each container should start a graceful shutdown of the running application and exit.
3. If a container doesn’t terminate within the grace period, a SIGKILL signal will be sent and the container violently terminated.

## Termination grace period
```
terminationGracePeriodSeconds: 60
```
