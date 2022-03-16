FROM golang:1.17.8
ENV GOOS linux
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go get -v 	"github.com/prometheus/client_golang/prometheus" \
	"github.com/prometheus/client_golang/prometheus/promhttp" \
	"k8s.io/apimachinery/pkg/apis/meta/v1@kubernetes-1.20.0" \
	"k8s.io/client-go/kubernetes@kubernetes-1.20.0" \
	"k8s.io/client-go/rest@kubernetes-1.20.0" 
RUN go build -o nodestats . && chmod +x nodestats
CMD ["/app/nodestats"]