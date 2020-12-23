FROM alpine:3
COPY ecr-scan /bin/ecr-scan
ENTRYPOINT [ "ecr-scan" ]
