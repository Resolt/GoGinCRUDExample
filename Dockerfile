FROM golang:1.17-alpine AS build
ADD ./ /app
WORKDIR /app
RUN go build -o ggce

FROM alpine
WORKDIR /app
COPY --from=build /app/ggce .
CMD ["./ggce"]