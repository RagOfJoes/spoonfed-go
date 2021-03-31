# =========== Build executable binary =========== #
FROM golang:alpine AS builder

# Alpine images doesn't include base tools
# so install necessary tools
RUN apk update && apk add --no-cache git make 

WORKDIR /build

# COPY GO MODULE FILES TO ALLOW FOR CACHING OF MODULE FETCHING
COPY . /build/

RUN go mod download
RUN go mod verify

RUN make build

# =========== Build small Docker image =========== #
FROM alpine AS prod

# Add maintainer label
LABEL maintainer "Victor Ragojos <vhsvragojos@gmail.com>"

RUN apk --no-cache add ca-certificates

# Copy built binary to directory
COPY --from=builder build/app ./spoonfed-go/

CMD [ "./app" ]
