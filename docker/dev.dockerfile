FROM golang:alpine

# Alpine images doesn't include base tools
# so install necessary tools
RUN apk update && apk upgrade \
  && apk add --no-cache git make

WORKDIR /spoonfed-go

# Copies source code into container working dir
COPY . .
# Download project deps
# Add realize for hot reaload
# See: https://github.com/cosmtrek/air/issues/114
RUN go mod download \
  && go get -v github.com/cosmtrek/air@b538c70423fb3590435c003dda15bf6a2f61187c \
  && echo "Finished downloading dependencies"
# Verify modules
RUN go mod verify

EXPOSE 8080

CMD [ "air", "-c", ".air.toml" ]
