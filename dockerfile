FROM ubuntu
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y
RUN apt-get install -y apt-utils
RUN apt-get install -y software-properties-common
RUN apt-get update -y
RUN apt-get install -y python3 python3-pip curl
RUN curl -OL https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
ENV PATH="$PATH:/usr/local/go/bin"
WORKDIR /usr/src/app
COPY . .
RUN go mod download && go mod verify
WORKDIR /usr/src/app/cmd/cs
RUN go build .
WORKDIR /usr/src/app/cmd/plugins/analyzer/pii
RUN go build .
WORKDIR /usr/src/app/cmd/plugins/connectors/workspace
RUN go build .
WORKDIR /usr/src/app/cmd/plugins/transformer/openxml
RUN go build .


WORKDIR /usr/src/app/cmd/plugins/transformer/pdf

RUN pip3 install -r requirements.txt
WORKDIR /usr/src/app

#CMD ["/bin/bash"]
WORKDIR /proj
CMD [ "/usr/src/app/cmd/cs/cs" ]