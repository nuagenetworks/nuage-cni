all: build_docker_image run_container
build_all: build_nuage_cni

build_docker_image:
	docker build --build-arg REGISTRY_URL=${REGISTRY_URL} -t ${REGISTRY_URL}/nuage-cni -f Dockerfile.build .

run_container:
	./docker_build_script.sh build_all

build_nuage_cni:
	./scripts/buildRPM.sh
	sh scripts/build-nuage-cni-docker.sh
