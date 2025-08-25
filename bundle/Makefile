DOCKER_CMD 	?= $(shell command -v podman || command -v docker)

bundle.tgz: bundle
	tar -C bundle -czf bundle.tgz rootfs init.js
	cp bundle.tgz ../session/bundle.tgz

bundle:
	mkdir -p bundle/rootfs
	$(DOCKER_CMD) build -t apptron-build --platform linux/386 -f Dockerfile .
	$(DOCKER_CMD) create --name apptron-export --platform linux/386 apptron-build
	$(DOCKER_CMD) export apptron-export | tar -C bundle/rootfs -xf -
	$(DOCKER_CMD) rm apptron-export
	cp init.js bundle/init.js

clean:
	rm -f bundle.tgz
	rm -f ../session/bundle.tgz
	rm -rf bundle
.PHONY: clean
