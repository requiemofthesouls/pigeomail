##
## Build
## docker build -t pigeomail -f build/package/Dockerfile .
##
FROM golang:1.20-buster AS build

WORKDIR /app
ADD . .
RUN go mod vendor

RUN go build -o /app/build/ci/pigeomail

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /
COPY --from=build /app/build/ci/pigeomail /pigeomail
EXPOSE 21025
EXPOSE 8443
USER nonroot:nonroot

ENTRYPOINT ["./pigeomail"]