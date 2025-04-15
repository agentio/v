FROM golang:1.22.5 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o v .

FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:slim
COPY --from=builder /app/v /usr/local/bin/v
RUN apt-get install jq vim -y
