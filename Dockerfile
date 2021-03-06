FROM heroku/heroku:18-build as build

COPY . /app
WORKDIR /app

# Setup buildpack
RUN mkdir -p /tmp/buildpack/heroku/go /tmp/build_cache /tmp/env
RUN curl https://codon-buildpacks.s3.amazonaws.com/buildpacks/heroku/go.tgz | tar xz -C /tmp/buildpack/heroku/go

#Execute Buildpack
RUN STACK=heroku-18 /tmp/buildpack/heroku/go/bin/compile /app/app/heroku/* /tmp/build_cache /tmp/env

# Prepare final, minimal image
FROM heroku/heroku:18

COPY --from=build /app /app
ENV HOME /app
WORKDIR /app

ADD ./.profile.d /app/.profile.d
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

RUN useradd -m heroku
USER heroku
CMD bash heroku-exec.sh && /app/bin/heroku

# FROM golang:1.12-alpine

# RUN mkdir /app && \
#     apk add git

# ADD . /app
# WORKDIR /app

# RUN mkdir -p /app/bin && \
#     go build -o bin/server src/server.go

# ENTRYPOINT /app/bin/server

# EXPOSE 2010
# EXPOSE 7080
