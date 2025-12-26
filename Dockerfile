FROM golang:1.25.5-alpine3.23 AS builder

RUN apk add --update nodejs npm && node --version && npm --version && npm install -g pnpm@latest-10

WORKDIR /build
COPY . .

# build ui
RUN export CI=true &&  cd ui && pnpm i && pnpm build

# build backend
RUN go build -o gogame ./main/

FROM public.ecr.aws/docker/library/alpine:3.20

COPY --from=builder /build/gogame /gogame

CMD ["/gogame"]
