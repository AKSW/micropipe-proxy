FROM node

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY package.json /usr/src/app/
RUN npm install
COPY . /usr/src/app

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=1s \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./micropipe-proxy-linux"]
