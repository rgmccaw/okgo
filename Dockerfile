FROM scratch
MAINTAINER Ryan McCaw "rgmccaw@gmail.com"
EXPOSE 8080

WORKDIR /

# copy binary into image
COPY okgo /

ENTRYPOINT ["./okgo"]
