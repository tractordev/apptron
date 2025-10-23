DOCKER_CMD 	?= $(shell command -v podman || command -v docker)

bzImage:
	$(DOCKER_CMD) build -t ghcr.io/tractordev/apptron:kernel -f Dockerfile .
	$(DOCKER_CMD) rm -f apptron-kernel-export
	$(DOCKER_CMD) create --name apptron-kernel-export ghcr.io/tractordev/apptron:kernel
	$(DOCKER_CMD) cp apptron-kernel-export:/bzImage ./bzImage
	ls -lah bzImage

clean:
	rm -f bzImage
.PHONY: clean