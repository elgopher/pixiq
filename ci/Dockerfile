FROM ubuntu:bionic

# Build Mesa
RUN apt-get update && \
    apt-get install -y libtool-bin autoconf python-pip libx11-dev libxext-dev x11proto-core-dev x11proto-gl-dev libglew-dev bison flex xvfb wget pkg-config zlib1g-dev llvm-dev && \
    wget https://mesa.freedesktop.org/archive/mesa-18.2.4.tar.xz && \
    tar xf mesa-18.2.4.tar.xz && \
    rm mesa-18.2.4.tar.xz && \
    mkdir mesa-18.2.4/build && \
    cd mesa-18.2.4/build && \
    ../configure --disable-dri \
               --disable-egl \
               --disable-gbm \
               --with-gallium-drivers=swrast,swr \
               --with-platforms=x11 \
               --prefix=/usr/local/ \
               --enable-gallium-osmesa \
               --disable-xvmc --disable-vdpau --disable-va \
               --with-swr-archs=avx && \
    make && \
    make install && \
    cd ../.. && \
    rm mesa-18.2.4 -Rf && \
    apt-get remove --auto-remove -y libtool-bin autoconf python-pip

# Setup our environment variables.
ENV XVFB_WHD="1920x1080x24"\
    DISPLAY=":99" \
    LIBGL_ALWAYS_SOFTWARE="1" \
    GALLIUM_DRIVER="softpipe" \
    MESA_DEBUG="incomplete_tex,incomplete_fbo"

# Install git
RUN apt-get install -y git

# Install GLFW dependencies
RUN apt-get install -y libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev

# Install go
RUN wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.14.2.linux-amd64.tar.gz && \
    rm go1.14.2.linux-amd64.tar.gz
ENV GOPATH=/opt/go/ PATH=$PATH:/usr/local/go/bin:/opt/go/bin

# Install make
RUN apt-get install -y make

# Install golint
RUN go get -u golang.org/x/lint/golint

# Install golangci-lint
RUN wget -O - https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.26.0
