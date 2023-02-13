FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
ENV PATH=/bin

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=certs /etc/passwd /etc/passwd

COPY snr /usr/local/bin/snr

WORKDIR /data

VOLUME /data

EXPOSE 7447/tcp

ENV SQLITE_DATABASE=/data/snr.db
ENTRYPOINT [ "/usr/local/bin/snr" ]
