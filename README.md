# go-nodestats
node statistics at-a-glance and formatted for prometheus

the main goroutine of this project will query the apiserver using incluster config for all pods and maintain a guage type prometheus metric with labels of namespace and nodename for watching pods move around the cluster.

the chart in this repository is satasfactory to set up the deployment, service, serviceaccount and clusterrolebinding (to the "view" clusterrole)

# build
this is meant to run in a container, for the docker build instructions check the dockerfile, but if you desire to build this outside of docker:

go get -v 	"github.com/prometheus/client_golang/prometheus" \
	"github.com/prometheus/client_golang/prometheus/promhttp" \
	"k8s.io/apimachinery/pkg/apis/meta/v1@kubernetes-1.20.0" \
	"k8s.io/client-go/kubernetes@kubernetes-1.20.0" \
	"k8s.io/client-go/rest@kubernetes-1.20.0" 

then

go build -o nodestats .

