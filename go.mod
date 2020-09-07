module github.com/form3tech-oss/cilium-etcd-watchdog

go 1.14

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1

replace k8s.io/api => k8s.io/api v0.15.12

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.15.12

replace k8s.io/apimachinery => k8s.io/apimachinery v0.15.13-beta.0

replace k8s.io/apiserver => k8s.io/apiserver v0.15.12

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.15.12

replace k8s.io/client-go => k8s.io/client-go v0.15.12

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.15.12

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.15.12

replace k8s.io/code-generator => k8s.io/code-generator v0.15.13-beta.0

replace k8s.io/component-base => k8s.io/component-base v0.15.12

replace k8s.io/cri-api => k8s.io/cri-api v0.15.13-beta.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.15.12

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.15.12

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.15.12

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.15.12

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.15.12

replace k8s.io/kubectl => k8s.io/kubectl v0.15.13-beta.0

replace k8s.io/kubelet => k8s.io/kubelet v0.15.12

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.15.12

replace k8s.io/metrics => k8s.io/metrics v0.15.12

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.15.12

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.15.12

replace k8s.io/sample-controller => k8s.io/sample-controller v0.15.12

require (
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/etcd-operator v0.9.4
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sys v0.0.0-20200622214017-ed371f2e16b4 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/klog v1.0.0
	k8s.io/kubernetes v1.15.12
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace k8s.io/node-api => k8s.io/node-api v0.15.12
