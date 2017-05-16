all:
	go build -o nuage-cni-mesos nuage-cni.go
	go build -o nuage-cni-k8s nuage-cni.go
	go build -o nuage-cni-openshift nuage-cni.go

fmt: 
	go fmt ./...

lint:
	cd daemon; go install; cd ..
	cd client; go install; cd ..
	cd k8s; go install; cd ..
	go install
	gometalinter --disable=dupl --disable=gocyclo --disable=aligncheck --disable=staticcheck --disable=gas --deadline 300s client
	gometalinter --disable=dupl --disable=gocyclo --disable=aligncheck --disable=staticcheck --disable=gas --deadline 300s daemon
	gometalinter --disable=dupl --disable=gocyclo --disable=aligncheck --disable=staticcheck --disable=gas --deadline 300s k8s
	gometalinter --disable=dupl --disable=gocyclo --disable=aligncheck --disable=staticcheck --disable=gas --deadline 300s .
