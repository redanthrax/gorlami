FROM golang:1.20.5

ENV NODE_VERSION=19.9.0
RUN apt install -y curl
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
ENV NVM_DIR=/root/.nvm

RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}
ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"
RUN npm i -D postcss postcss-cli

COPY . /go/src/redanthrax/gorlami-server
WORKDIR /go/src/redanthrax/gorlami-server/frontend
RUN npm install

WORKDIR /go/src/redanthrax/gorlami-server

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY . .

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -polling -log-prefix=false\
    -build="go build ." -command="./server"\
    -directory="./" -recursive="true"