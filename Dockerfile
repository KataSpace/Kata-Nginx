FROM vikings/alpine
ADD  bin/kn /kn
ENTRYPOINT ["/kn"]