FROM scratch

COPY sctl /

WORKDIR /root
ENTRYPOINT ["/sctl"]
