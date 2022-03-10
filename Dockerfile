FROM golang:1.17 AS build
ADD ./ /build
WORKDIR /build
RUN go build -o ggce

FROM gcr.io/distroless/base
WORKDIR /ggce
COPY --from=build /build/ggce .
CMD ["./ggce"]