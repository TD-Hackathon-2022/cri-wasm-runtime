package backend

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/util"

	"github.com/diannaowa/cri-template/core"
)

// maxMsgSize use 16MB as the default message size limit.
// grpc library default is 4MB
const maxMsgSize = 1024 * 1024 * 16

// CriTemplateServer is the grpc backend of cri-dockerd.
type CriTemplateServer struct {
	// endpoint is the endpoint to serve on.
	endpoint string
	// service is the docker service which implements runtime and image services.
	service core.CRIService
	// server is the grpc server.
	server *grpc.Server
}

// NewCriTemplateServer creates the cri-dockerd grpc backend.
func NewCriTemplateServer(endpoint string, s core.CRIService) *CriTemplateServer {
	return &CriTemplateServer{
		endpoint: endpoint,
		service:  s,
	}
}

func getListener(addr string) (net.Listener, error) {
	addrSlice := strings.SplitN(addr, "://", 2)
	proto := addrSlice[0]
	listenAddr := addrSlice[1]
	switch proto {
	case "fd":
		return listenFD(listenAddr)
	default:
		return util.CreateListener(addr)
	}
}

// Start starts the cri-template grpc backend.
func (s *CriTemplateServer) Start() error {
	// Start the internal service.
	if err := s.service.Start(); err != nil {
		logrus.Error(err, "Unable to start cri-template service")
		return err
	}

	logrus.Info("Start cri-template grpc backend")
	l, err := getListener(s.endpoint)
	if err != nil {
		return fmt.Errorf("cri-template failed to listen on %q: %v", s.endpoint, err)
	}
	// Create the grpc backend and register runtime and image services.
	s.server = grpc.NewServer(
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)
	runtimeapi.RegisterRuntimeServiceServer(s.server, s.service)
	runtimeapi.RegisterImageServiceServer(s.server, s.service)
	go func() {
		if err := s.server.Serve(l); err != nil {
			logrus.Error(err, "Failed to serve connections from cri-template")
			os.Exit(1)
		}
	}()
	handleNotify()
	return nil
}
