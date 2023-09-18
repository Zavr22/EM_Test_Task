FROM ubuntu:latest
LABEL authors="misha"

ENTRYPOINT ["top", "-b"]