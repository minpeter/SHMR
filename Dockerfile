FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    curl

RUN useradd -ms /bin/bash runner
WORKDIR /home/runner

RUN mkdir actions-runner && cd actions-runner
RUN curl -o actions-runner-linux-x64-2.304.0.tar.gz \
    -L https://github.com/actions/runner/releases/download/v2.304.0/actions-runner-linux-x64-2.304.0.tar.gz
RUN tar xzf ./actions-runner-linux-x64-2.304.0.tar.gz
RUN ./bin/installdependencies.sh
RUN chown -R runner:runner /home/runner/actions-runner

USER runner

ENV URL=""
ENV TOKEN=""

CMD echo "\n\n\n\n" | ./config.sh --url $URL --token $TOKEN && ./run.sh