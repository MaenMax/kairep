FROM git.kaiostech.com:4567/cloud/docker-buildenv/generic

RUN yum install -y \
    rsyslog

ADD templates/hostname /bin/hostname
ADD configs/rsyslog.conf /etc/rsyslog.conf
RUN chmod +x /bin/hostname

RUN mkdir -p /data/autopush
ADD ./ /data/autopush

WORKDIR /data/autopush

ADD docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]
#CMD ["/docker-entrypoint.sh"]
