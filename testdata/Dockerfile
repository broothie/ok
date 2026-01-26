FROM alpine:latest

# Install dependencies and shells
RUN apk add --no-cache \
    curl \
    bash \
    zsh \
    make \
    ruby \
    nodejs \
    npm \
    git \
    go

# Install Rake
RUN gem install rake

# Install Yarn
RUN npm install -g yarn

# Install Task (go-task)
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Install Just
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

# Verify installations
RUN echo "Verifying tools installation..." && \
    just --version && \
    make --version && \
    npm --version && \
    yarn --version && \
    rake --version && \
    task --version && \
    go version && \
    bash --version && \
    zsh --version

WORKDIR /workdir
ENV PATH="/root/go/bin:${PATH}"
COPY . .
RUN go install ./cmd/ok
WORKDIR /workdir/testdata/all-tools

CMD ["sleep", "infinity"]
