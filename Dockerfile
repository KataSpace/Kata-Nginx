FROM debian:10
RUN apt-get update && \
    apt-get install -y procps
ADD  bin/kn /kn
ENTRYPOINT ["/kn"]