FROM iron/base
MAINTAINER Track "24245@163.com"
LABEL build_time="20190516172321"
RUN mkdir -p /home/blogserver/api_docs && mkdir -p /var/log/blogserver && mkdir -p /home/blogserver/config && mkdir -p /home/blogserver/images
COPY ./docs/swaggerui/ /home/blogserver/api_docs/
COPY ./config/config.toml /home/blogserver/config/
ENV APP_CONFIG_PATH /home/blogserver/config/config.toml
EXPOSE 8888
ADD  ./build/blog_server-linux-amd64-20190516172321 /
ENTRYPOINT ["./blog_server-linux-amd64-20190516172321"]
