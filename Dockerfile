FROM debian:latest

# RUN apk add --no-cache ca-certificates

COPY view/ /docer/view/
COPY docer /docer/docer
COPY public/ /docer/public/

EXPOSE 8089

WORKDIR /docer

ENTRYPOINT /docer/docer

