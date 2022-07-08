# specify the base image for the go app.
FROM golang:1.18

#specify that we now need to execute any commands in this directory
#WORKDIR /go/src/github.com/postgres-go
WORKDIR /go/src/github.com/postgres-go

#copy everything from this project into the filesystem of the container
COPY . .

#obtain the package needed to run code. alternatively use go modules
RUN go get -u gorm.io/gorm
RUN go get gorm.io/driver/postgres

#compile the binary exe for our app
RUN go build -o main .

#start the application
CMD ["./main"]