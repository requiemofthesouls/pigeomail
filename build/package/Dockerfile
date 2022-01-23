##
## Build
## docker build -t pigeomail -f build/package/Dockerfile .
##
FROM golang:1.17-buster AS build

WORKDIR /app
ADD . .
RUN go mod vendor

RUN go build -o ./build/ci/pigeomail

##
## Deploy
##
FROM alpine:3.15

WORKDIR /
COPY --from=build /app/build/ci/pigeomail /pigeomail
EXPOSE 21025
USER nonroot:nonroot

ENTRYPOINT ["/pigeomail"]