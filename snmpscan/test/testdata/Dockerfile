# FROM python:3.7-slim
# WORKDIR /app
# RUN pip install snmpsim
# COPY public.snmprec .
# RUN adduser --system snmpsim
# EXPOSE 161/udp

# CMD ["snmpsimd.py","--data-dir=/app","--agent-udpv4-endpoint=0.0.0.0:161","--process-user=snmpsim","--process-group=nogroup" ]


FROM python:3.7-slim
WORKDIR /app
RUN pip install snmpsim
# add user as root
RUN useradd -m docker && echo "docker:docker" | chpasswd && adduser docker sudo
# switch user
USER docker
COPY ./public.snmprec /app

EXPOSE 161/udp
# execute snmpsimd as user=docekr (root)
CMD snmpsimd.py --agent-udpv4-endpoint=0.0.0.0:161 --process-user=docker --process-group=docker --data-dir=/app --$EXTRA_FLAGS