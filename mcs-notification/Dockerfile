FROM golang:1.21.6

RUN apt-get update && apt-get install -y ca-certificates openssl

ARG cert_location=/usr/local/share/ca-certificates
# Get certificate from "github.com"
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
# Get certificate from "proxy.golang.org"
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/proxy.golang.crt
# Update certificates
RUN update-ca-certificates

ENV genmqpkg_incnls=1
ENV genmqpkg_incsdk=1
ENV genmqpkg_inctls=1

RUN mkdir -p /opt/mqm
RUN cd /opt/mqm
RUN curl -LO "https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/9.3.4.1-IBM-MQC-Redist-LinuxX64.tar.gz"
RUN tar -zxf ./*.tar.gz
RUN rm -f ./*.tar.gz
RUN bin/genmqpkg.sh -b /opt/mqm

RUN apt install libxml2-dev libxslt1-dev liblzma-dev zlib1g-dev -y

WORKDIR /app

ARG SERVICE_NAME=go-clean-architecture
RUN mkdir -p bin && mkdir -p ${SERVICE_NAME}

# Download lib
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Build file binary
WORKDIR /app/${SERVICE_NAME}
COPY . ./
WORKDIR /app/${SERVICE_NAME}/cmd
RUN go build -o /app/bin/main .
WORKDIR /app
RUN rm -rf ${SERVICE_NAME}/

EXPOSE 8088
EXPOSE 8089

ENV GOOS linux
ENV CGO_ENABLED 0

CMD ["/app/bin/main"] 