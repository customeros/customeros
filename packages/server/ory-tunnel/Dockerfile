FROM ubuntu:22.04
RUN echo 'APT::Install-Suggests "0";' >> /etc/apt/apt.conf.d/00-docker
RUN echo 'APT::Install-Recommends "0";' >> /etc/apt/apt.conf.d/00-docker
RUN DEBIAN_FRONTEND=noninteractive \
  apt-get update \
  && apt-get install -y curl bash ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /tmp/
RUN bash -c 'bash <(curl https://raw.githubusercontent.com/ory/meta/master/install.sh) -b . ory'


RUN  mv ./ory /usr/local/bin/
COPY tunnel.sh /usr/local/bin/

ENTRYPOINT ["tail", "-f", "/dev/null"]
