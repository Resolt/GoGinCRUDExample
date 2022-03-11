FROM golang:1.17-alpine AS build
ADD ./ /build
WORKDIR /build
RUN go build -o ggce

FROM alpine
WORKDIR /ggce
COPY --from=build /build/ggce .
CMD ["./ggce"]