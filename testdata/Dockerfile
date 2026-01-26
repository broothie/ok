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

# Install Nx
RUN npm install -g nx

# Install Task (go-task)
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Install Just
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

WORKDIR /workdir
ENV PATH="/root/go/bin:${PATH}"
COPY . .
RUN go install ./cmd/ok

WORKDIR /workdir/testdata/all-tools
RUN npm install

# Verify installations
RUN echo "Verifying tools installation..." && \
    ok --version && \
    just --version && \
    make --version && \
    npm --version && \
    yarn --version && \
    nx --version && \
    rake --version && \
    task --version && \
    bash --version && \
    zsh --version
