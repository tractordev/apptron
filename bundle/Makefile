DOCKER_CMD 	?= $(shell command -v podman || command -v docker)

bundle.tgz: bundle
	tar -C bundle -czf bundle.tgz .
	cp bundle.tgz ../session/bundle.tgz
	ls -lah bundle.tgz

bundle:
	$(DOCKER_CMD) build -t apptron-build --platform linux/386 -f Dockerfile .
	$(DOCKER_CMD) rm -f apptron-export && mkdir -p bundle
	$(DOCKER_CMD) create --name apptron-export --platform linux/386 apptron-build
	$(DOCKER_CMD) export apptron-export | tar -C bundle -xf -


clean:
	rm -f bundle.tgz
	rm -f ../session/bundle.tgz
	rm -rf bundle
.PHONY: clean
