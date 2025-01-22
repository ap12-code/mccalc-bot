FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /mccalc-bot

# ------------

FROM alpine:3.12

RUN apk update
RUN apk add tzdata bash curl

COPY --from=builder /mccalc-bot /bin/

# user
RUN adduser -D bot && chown -R bot /bin
USER bot


CMD [ "/bin/mccalc-bot" ]