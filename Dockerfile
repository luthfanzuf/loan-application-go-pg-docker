### STEP 1 BUILD BINARY

# specify the base image for the go app.
FROM golang:1.18-alpine3.16 AS builder

#specify that we now need to execute any commands in this directory
#WORKDIR /go/src/github.com/postgres-go
WORKDIR /app

#copy everything from this project into the filesystem of the container
COPY . .

#obtain the package needed to run code. alternatively use go modules
RUN go get -u gorm.io/gorm && \
    go get gorm.io/driver/postgres && \
    #compile the binary exe for our app
    go build -o main .

# STEP 2 COPY BINARY TO NEW ALPINE 

FROM alpine:3.16
WORKDIR /app

COPY --from=builder /app/main .

#start the application
CMD ["./main"]