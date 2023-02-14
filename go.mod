module github.com/nuagenetworks/nuage-cni

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/ccding/go-logging v0.0.0-20190618175518-0ac4cc1a6533 // indirect
	github.com/containernetworking/cni v0.3.1-0.20161010053931-d872391998fb
	github.com/coreos/go-iptables v0.1.1-0.20160907220151-5463fbac3bcc // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.0.0-00010101000000-000000000000 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nuagenetworks/libvrsdk v0.0.0-20200625144000-d7373f6f983c
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vishvananda/netlink v0.0.0-20151203164549-edcd99c0881a
	github.com/vishvananda/netns v0.0.0-20160430053723-8ba1072b58e0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab // indirect
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20200124190032-861946025e34 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d
