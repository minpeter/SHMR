FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    curl

RUN useradd -ms /bin/bash runner
WORKDIR /home/runner

ENV RUNNER_VERSION="2.304.0"

RUN mkdir actions-runner && cd actions-runner
RUN curl -o actions-runner-linux-x64-$RUNNER_VERSION.tar.gz \
    -L https://github.com/actions/runner/releases/download/v$RUNNER_VERSION/actions-runner-linux-x64-$RUNNER_VERSION.tar.gz
RUN tar xzf ./actions-runner-linux-x64-$RUNNER_VERSION.tar.gz
RUN ./bin/installdependencies.sh
RUN chown -R runner:runner /home/runner/actions-runner

USER runner

CMD echo "\n\n\n\n" | ./config.sh --url $URL --token $TOKEN && ./run.sh