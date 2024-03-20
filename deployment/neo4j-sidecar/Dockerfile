FROM debian:bullseye

RUN apt-get update && \
      apt-get install -y wget && \
      wget -O /tmp/cypher-shell_5.18.1_all.deb https://dist.neo4j.org/cypher-shell/cypher-shell_5.18.1_all.deb && \
      apt-get install -y openjdk-17-jdk && \
      dpkg -i /tmp/cypher-shell_5.18.1_all.deb
