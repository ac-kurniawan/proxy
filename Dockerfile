FROM golang:1.21 as build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod tidy
COPY . ./
RUN CGO_ENABLED=0 go build -o main

FROM gcr.io/distroless/static-debian11
WORKDIR /app
USER nonroot:nonroot
COPY --from=build /app/main /app/main
COPY --from=build /app/properties.yml /app/properties.yml
EXPOSE 8080
ENTRYPOINT ["./main"]