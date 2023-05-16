FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl

ARG DOCKER_VERSION="23.0.6"
RUN curl -fsSL https://download.docker.com/linux/static/stable/$(uname -m)/docker-$DOCKER_VERSION.tgz | tar zxvf - --strip 1 -C /usr/local/bin docker/docker


ARG HOST_DOCKER_GID=995
RUN groupadd -g ${HOST_DOCKER_GID} docker
RUN useradd -ms /bin/bash -g docker runner

ARG HOME="/home/runner"
WORKDIR $HOME

ARG DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}
RUN mkdir -p $DOCKER_CONFIG/cli-plugins \
    && curl -SL https://github.com/docker/compose/releases/download/$(curl -s https://api.github.com/repos/docker/compose/releases/latest | \
    grep tag_name | cut -d '"' -f 4)/docker-compose-$(echo $(uname -s)-$(uname -m) | tr '[A-Z]' '[a-z]') -o $DOCKER_CONFIG/cli-plugins/docker-compose \
    && chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose

RUN RUNNER_VERSION=$(curl -L https://api.github.com/repos/actions/runner/releases/latest | grep tag_name | cut -d '"' -f 4 | sed 's/v//g') \
    && RUNNER_ARCH=$(if [ "$(uname -m)" = "aarch64" ]; then echo "arm64"; elif [ "$(uname -m)" = "x86_64" ]; then echo "x64"; else echo "Unsupported architecture $(uname -m)" >&2; fi) \
    && curl -L https://github.com/actions/runner/releases/download/v$RUNNER_VERSION/actions-runner-linux-$RUNNER_ARCH-$RUNNER_VERSION.tar.gz | tar xz 

RUN ./bin/installdependencies.sh

USER runner

CMD echo "\n\n\n\n" | ./config.sh --url $URL --token $TOKEN && ./run.sh