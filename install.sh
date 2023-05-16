#!/bin/sh

# 도커 설치 확인
if ! [ -x "$(command -v docker)" ]; then
  echo 'Error: docker is not installed.' >&2
  exit 1
fi

# curl 설치 확인
if ! [ -x "$(command -v curl)" ]; then
  echo 'Error: curl is not installed.' >&2
  exit 1
fi

if [ -d ~/.SHMR ]; then
  rm -rf ~/.SHMR
fi

mkdir -p ~/.SHMR/bin

cd ~/.SHMR

curl -O https://raw.githubusercontent.com/minpeter/SHMR/main/Dockerfile

SHMR_VERSION=$(curl -L https://api.github.com/repos/minpeter/SHMR/releases/latest | grep tag_name | cut -d '"' -f 4 | sed 's/v//g') \
SHMR_ARCH=$(if [ "$(uname -m)" = "aarch64" ]; then echo $(uname -s)_"arm64" | tr '[A-Z]' '[a-z]'; elif [ "$(uname -m)" = "x86_64" ]; then echo $(uname -s)_"amd64" | tr '[A-Z]' '[a-z]'; else echo "Unsupported architecture $(uname -m)" >&2; fi) \
curl -L https://github.com/minpeter/SHMR/releases/download/$SHMR_VERSION/SHMR_$SHMR_VERSION"_"$SHMR_ARCH.tar.gz | tar xz

docker build  --build-arg HOST_DOCKER_GID=$(stat -c "%g" /var/run/docker.sock) \
              --build-arg DOCKER_VERSION=$(docker version --format '{{.Client.Version}}') \
              -t gh-action-runner:latest .

mv SHMR ~/.SHMR/bin/shmr

if [ -f ~/.bashrc ]; then
  echo "export PATH=\$PATH:~/.SHMR/bin" >> ~/.bashrc
elif [ -f ~/.zshrc ]; then
  echo "export PATH=\$PATH:~/.SHMR/bin" >> ~/.zshrc
fi