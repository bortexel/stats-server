FROM golang:latest
WORKDIR /home/container
ADD stats .
ENV PORT 8080
HEALTHCHECK --interval=1m --timeout=3s \
  CMD curl -f http://127.0.0.1:8080/ || exit 1
EXPOSE 8080
CMD [ "./stats" ]
