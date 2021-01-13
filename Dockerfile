FROM golang:latest as build
WORKDIR /app
COPY . .
RUN make server

FROM alpine:latest
COPY --from=build /app/bin/fantasy /fantasy
RUN chmod +x /fantasy
CMD ["/fantasy"]
