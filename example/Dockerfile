FROM ubuntu:latest

COPY confd-6.6.linux.x86_64.installer.bin /

RUN apt update && \
  apt install -y openssl openssh-client && \
  rm -rf /var/lib/apt/lists/* && \
  /confd-6.6.linux.x86_64.installer.bin /opt/confd && \
  rm /confd-6.6.linux.x86_64.installer.bin 

COPY root.fxs confd.conf /opt/confd/etc/confd/

EXPOSE 22 830 80 4565

CMD bash -c "source /opt/confd/confdrc && \
   confd --conf /opt/confd/etc/confd/confd.conf --foreground -v"
