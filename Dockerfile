FROM golang:1.18-alpine AS build
ADD ./ /app
WORKDIR /app
RUN apk add git
RUN go build -o ggce

FROM alpine:latest AS final
ARG USER=nonroot
RUN adduser -D $USER
USER $USER
WORKDIR /app
COPY --from=build --chown=$USER:$USER /app/ggce .
CMD ["./ggce"]