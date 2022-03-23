FROM golang:1.18-alpine AS build
ADD ./ /app
WORKDIR /app
RUN apk add git
RUN go build -o ggce

FROM alpine
WORKDIR /app
COPY --from=build /app/ggce .
CMD ["./ggce"]