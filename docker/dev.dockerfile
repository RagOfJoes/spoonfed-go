FROM golang:alpine

# Alpine images doesn't include base tools
# so install necessary tools
RUN apk update && apk upgrade \
  && apk add --no-cache git make

WORKDIR /spoonfed-go

# Copies source code into container working dir
COPY . .

# Download project deps
# Add air for hot reload
RUN go mod download \
  && go get -u -v github.com/cosmtrek/air@master 

EXPOSE 8080

CMD [ "air", "-c", ".air.toml" ]
