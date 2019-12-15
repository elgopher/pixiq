FROM utensils/opengl:19.0.8

# Install git
RUN apk add --update git openssh-client

# Install GLFW dependencies
RUN apk add --update build-base libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev xorg-server-dev

# Install sound
RUN apk add --update alsa-lib-dev

# Install fresh version of Go
RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl go && \
    wget -O go.tgz https://dl.google.com/go/go1.13.5.src.tar.gz && \
    tar -C /usr/local -xzf go.tgz && \
    rm go.tgz && \
    cd /usr/local/go/src/ && \
    ./make.bash && \
    apk del .build-deps
ENV GOPATH=/opt/go/ PATH=$PATH:/usr/local/go/bin:/opt/go/bin

# Install golint
RUN go get -u golang.org/x/lint/golint
