FROM netdata/netdata:stable

COPY go.d.conf /etc/netdata/go.d.conf

RUN apt-get update && apt-get install -y netcat-traditional

COPY entrypoint.sh /usr/local/bin/custom-entrypoint.sh
RUN chmod +x /usr/local/bin/custom-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/custom-entrypoint.sh"]