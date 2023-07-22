FROM golang:1.20

WORKDIR /app

RUN apt-get update
RUN apt install -y protobuf-compiler

#RUN useradd -u 1000 -m -s /bin/bash user \
#    && echo "[url \"git@gitlab.com:\"]" | tee -a /root/.gitconfig /home/user/.gitconfig \
#    && echo "    insteadOf = https://gitlab.com/" | tee -a /root/.gitconfig /home/user/.gitconfig \
#    && mkdir -p /root/.ssh /home/user/.ssh /app \
#    && chown -R 1000:1000 /home/user /app

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
