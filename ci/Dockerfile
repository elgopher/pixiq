FROM utensils/opengl:19.0.8

RUN apk add --update go git openssh-client

# Install GLFW dependencies
RUN apk add --update build-base libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev xorg-server-dev

# Install sound
RUN apk add --update alsa-lib-dev
