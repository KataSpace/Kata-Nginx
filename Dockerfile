FROM debian:10
ADD  bin/kn /kn
ENTRYPOINT ["/kn"]