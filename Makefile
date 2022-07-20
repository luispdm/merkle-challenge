export CFG ?= ../merkle.yml
COMPOSE	   := docker/docker-compose.yml

server:
	@ java \
		-jar \
		problem/merkle-tree-java.jar \
		problem/icons_rgb_circle.png \
		problem/example-file-merkle-tree.png

tests:
	@ ginkgo \
		-v \
		-p \
		-race \
		suites \
		-- \
		-c \
		${CFG}

tests-binary:
	@ cd suites && \
		CGO_ENABLED=0 \
		GOOS=linux \
		go test \
		-c \
		-v \
		-o \
		../merkle.test

dev:
	@ docker-compose \
		-f \
		${COMPOSE} \
		up

nodev:
	@ docker-compose \
		-f \
		${COMPOSE} \
		down

devrebuild:
	@ docker-compose \
		-f \
		${COMPOSE} \
		up \
		--build \
		--remove-orphans
