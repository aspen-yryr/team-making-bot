FROM node:19.4

WORKDIR /var/works

ARG PROTOBUF_VERSION=3.19.6
ARG GEN_GRPC_VERSION=1.4.2

RUN apt-get -qq update && apt-get -qq install -y curl unzip

RUN curl -sSL \
    https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protoc-$PROTOBUF_VERSION-linux-x86_64.zip \
    -o protoc.zip

RUN unzip -qq protoc.zip && \
    mv ./bin/protoc /usr/local/bin/protoc

RUN npm i grpc-web

RUN curl -sSL \
    https://github.com/grpc/grpc-web/releases/download/$GEN_GRPC_VERSION/protoc-gen-grpc-web-$GEN_GRPC_VERSION-linux-x86_64 \
    -o protoc-gen-grpc-web

RUN mv protoc-gen-grpc-web /usr/local/bin && \
    chmod +x /usr/local/bin/protoc-gen-grpc-web

COPY scripts/gen_proto.sh /var/works/

ENTRYPOINT [ "/var/works/gen_proto.sh" ]
