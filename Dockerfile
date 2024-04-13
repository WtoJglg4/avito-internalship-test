FROM golang:1.22-alpine AS builder

WORKDIR /usr/local/src/

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]

ENV GO111MODULE=on

RUN go mod download

COPY . ./
RUN go build -o ./bin/main cmd/main.go

FROM alpine AS runner

RUN apk update && apk add postgresql-client

COPY --from=builder /usr/local/src/bin/main /
# COPY docs ./docs/
COPY schema ./schema/
COPY .env /
COPY configs ./configs/

COPY wait-for-postgres.sh /
RUN chmod +x /wait-for-postgres.sh



CMD [ "/main" ]