FROM golang:1.10-alpine AS builder
RUN apk add git vim curl

# dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/noahklein/vimgolf
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch
COPY --from=builder /app ./
COPY --from=builder /usr/bin/vim /usr/bin/vim

ENV GIN_MODE="release"
ENTRYPOINT ["./app"]
