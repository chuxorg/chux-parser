# Use the official Golang image as the base image
FROM golang:1.19-bullseye

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY . .

# Build the Go application
RUN go build -o chux-parser  -ldflags  "-X main.BuildStamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.GitHash=`git rev-parse HEAD` -X main.Version=`git tag --sort=-version:refname | head -n 1`" main.go

RUN rm -rf config/ internal/ pkg/ app/

# Expose port 8080 to the host
EXPOSE 8080

# Run the binary when the container starts
CMD ["./chux-parser"]

# Leave the below commented out code
# ENTRYPOINT ["tail", "-f", "/dev/null"]