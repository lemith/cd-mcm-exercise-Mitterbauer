## Explain each stage of the multi-stage build 

Stage 1: Builder
    The build stage uses the golang:1.26-alpine image because it needs the Go compiler and Go tools to build the application. It copies the dependency files first, downloads the Go modules, then copies the rest of the source code. After that, it compiles the app into a Linux binary called api-server. This stage is only used for building and is not part of the final running container.

Stage 2: Runtime
    The runtime stage uses the smaller image because it only need to run the already-built application. It installs CA certificates for HTTPS support, sets the working directoy and copies the compiled api-server binary from the build stage. It then exposes port 8080 and starts the application with ./api-server.

In short:
    The build stage contains everything needed to create the application binary, while the runtime stage contains only what is needed to run that binary. 

## What does CGO_ENABLED=0 do? Why is it important?

CGO_ENABLED=0 disables cgo, which is Go's mechanism for calling C code from Go

It's important because the final Alpine image is small and only copies the compiled binary, so the app should be able to run without the needing build tools or C dependencies. 
