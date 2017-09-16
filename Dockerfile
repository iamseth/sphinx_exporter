FROM busybox:1.27.2

MAINTAINER Seth Miller <seth@sethmiller.me>
COPY ./dist/sphinx_exporter.linux-amd64 /bin/sphinx_exporter
EXPOSE 9247
CMD ["/bin/sphinx_exporter"]
