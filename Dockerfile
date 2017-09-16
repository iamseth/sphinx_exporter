FROM scratch
MAINTAINER Seth Miller <seth@sethmiller.me>
COPY ./dist/sphinx_exporter.linux-amd64 /bin/sphinx_exporter
EXPOSE     9247
ENTRYPOINT [ "/bin/sphinx_exporter" ]
