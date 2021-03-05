FROM golang:1.15-alpine AS build
RUN apk add make ca-certificates
WORKDIR /usr/src
COPY . .
RUN make
FROM alpine
COPY --from=build /usr/src/bin/jqrp.linux-amd64 /usr/local/bin/jqrp
ENTRYPOINT ["jqrp"]
LABEL maintainer="bauerdominic@gmail.com"
