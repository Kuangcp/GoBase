# build stage

FROM node:12.16-alpine AS builder

WORKDIR /data

COPY package.json  /data
COPY package-lock.json /data

RUN npm install

COPY . /data

RUN npm run build

# run stage

FROM python:3.9.6-alpine

COPY --from=builder /data/dist /data

WORKDIR /data

CMD ["/usr/local/bin/python3", "-m", "http.server"]

