FROM golang:1.17 as build

ADD ./ /build
WORKDIR /build
RUN go build -o ggce

FROM gcr.io/distroless/base-debian11
WORKDIR /ggce
COPY --from=build /build/ggce .
CMD ["./ggce"]