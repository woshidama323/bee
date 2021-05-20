FROM golang:1.15 AS build

WORKDIR /src
# enable modules caching in separate layer
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN make binary

FROM debian:10.2-slim

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*; \
    groupadd -r bee --gid 999; \
    useradd -r -g bee --uid 999 --no-log-init -m bee;

# make sure mounted volumes have correct permissions
RUN mkdir -p /home/bee/.bee && chown 999:999 /home/bee/.bee &&  \
       wget https://github.com/ethersphere/bee-clef/releases/download/v0.4.9/bee-clef_0.4.9_arm64.rpm &&  \
       rpm -i bee-clef_0.4.9_arm64.rpm
       

COPY --from=build /src/dist/bee /usr/local/bin/bee

EXPOSE 1633 1634 1635
USER root
WORKDIR /root
VOLUME /root/.bee

ENTRYPOINT ["bee","systemctl start bee-clef"]
