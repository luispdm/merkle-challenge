FROM openjdk:18.0.1.1-slim-buster

WORKDIR /opt

# install curl to check if service is healthy
# "apt-get clean" && "rm -rf..." suggested by hadolint:
# https://github.com/codacy/codacy-hadolint/blob/master/codacy-hadolint/docs/description/DL3009.md
RUN apt-get update && \
    apt-get install --no-install-recommends -y \
        curl=7.64.0-4+deb10u2 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    addgroup --gid 1000 servergroup && \
    adduser --disabled-password --gecos "" \
        --uid 1000 --gid 1000 serveruser

USER serveruser:servergroup

COPY problem/merkle-tree-java.jar merkle-tree-java.jar
COPY problem/icons_rgb_circle.png icons_rgb_circle.png
COPY problem/example-file-merkle-tree.png example-file-merkle-tree.png

ENTRYPOINT [ "java" ]
CMD [ "-jar", "merkle-tree-java.jar", "icons_rgb_circle.png", "example-file-merkle-tree.png" ]
