all: build_docker_image run_container
build_all: build_nuage_cni

build_docker_image:
	docker build --build-arg http_proxy=${http_proxy} --build-arg https_proxy=${https_proxy} --build-arg no_proxy=${no_proxy} -t ${REGISTRY_URL}/nuage-cni -f Dockerfile.build .

run_container:
	./docker_build_script.sh build_all

build_nuage_cni:
	./scripts/buildRPM.sh
	sh scripts/build-nuage-cni-docker.sh
