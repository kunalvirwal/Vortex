package master

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/kunalvirwal/Vortex/internal/docker"
	"github.com/kunalvirwal/Vortex/internal/dockmaster"
	"github.com/kunalvirwal/Vortex/internal/state"
	"github.com/kunalvirwal/Vortex/internal/utils"
	pb "github.com/kunalvirwal/Vortex/proto/factory"
	"github.com/kunalvirwal/Vortex/types"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

type server struct {
	pb.UnimplementedContainerFactoryServer
}

func (s *server) Apply(ctx context.Context, body *pb.RequestBody) (*pb.BoolResponse, error) {

	dep := &types.Deployment{}
	fmt.Println("Recieved Request")
	err := yaml.Unmarshal(body.GetData(), dep)
	if err != nil {
		return &pb.BoolResponse{Success: false}, err
	}
	fmt.Println("Parsed Deployment from yaml")

	// validate deployment values
	err = utils.ValidateDeployment(dep)
	if err != nil {
		fmt.Println(err)
		return &pb.BoolResponse{Success: false}, err
	}

	newDep := true
	// checks if a deployment with the same name already exists
	state.VortexDeployments.Mu.Lock()
	for _, vDeployment := range state.VortexDeployments.List {
		if dep.Version == vDeployment.Version {
			newDep = false
			break
		}
	}
	if newDep {
		state.VortexDeployments.List = append(state.VortexDeployments.List, dep)
	}
	state.VortexDeployments.Mu.Unlock()

	// Two services can have same name if they are in different deployments but can not have the same name in the same deployment
	// Checks if any service is repeated in the deployment
	var newServiceArr []*types.VService
	for _, service := range dep.Services {
		for _, vService := range newServiceArr {
			if service.Name == vService.Service.Name {
				return &pb.BoolResponse{Success: false}, errors.New("can not have two services with the same name in the same deployment, service /'" + service.Name + "/' is repeated")
			}
		}
		vService := types.VService{Service: service, Deployment: dep.Version}
		newServiceArr = append(newServiceArr, &vService)
	}

	var g errgroup.Group

	// Modify the services if they are already deployed and remove them from newServiceArr
	for _, service := range newServiceArr {
		for _, vService := range state.VortexServices.List {
			if service.Service.Name == vService.Service.Name && service.Deployment == vService.Deployment {
				g.Go(func() error {
					return dockmaster.Modify(vService, service)
				})
				for i, v := range newServiceArr {
					if v.Service.Name == service.Service.Name && v.Deployment == service.Deployment {
						newServiceArr = append(newServiceArr[:i], newServiceArr[i+1:]...)
						break
					}
				}
				break
			}
		}
	}

	// Deploy the new services
	for _, service := range newServiceArr {
		state.VortexContainers.Mu.Lock()
		state.VortexServices.List = append(state.VortexServices.List, service)
		state.VortexContainers.Mu.Unlock()
		g.Go(func() error {
			return dockmaster.Deploy(service)
		})
	}

	if err = g.Wait(); err != nil {
		go func() {
			timer := time.NewTimer(3 * time.Second)
			<-timer.C
			os.Exit(0)
		}()
		return &pb.BoolResponse{Success: false}, err
	}
	fmt.Println(state.GetState())
	return &pb.BoolResponse{Success: true}, nil
}

func (s *server) Delete(ctx context.Context, body *pb.NameHolder) (*pb.BoolResponse, error) {
	Name := body.GetName()
	query := strings.Split(Name, " ")

	// first remove the deployment, services and containers from the state variables
	// then delete the containers from docker sdk otherwise tracker will try to restart them
	// fmt.Printf("%v", query)
	var rmContainerIDs []string

	if len(query) == 1 {

		// Removing deployment from state
		dep := utils.PopDeployment(state.VortexDeployments, query[0])
		if dep == nil {
			return &pb.BoolResponse{Success: false}, errors.New("deployment not found")
		}

		// Removing services from state
		rmServices := utils.PopServicesByDepVersion(state.VortexServices, dep.Version)

		// Removing containers from VortexContainers
		for _, vService := range rmServices {
			utils.RemoveContainerConfigsByService(state.VortexContainers, vService)
			rmContainerIDs = append(rmContainerIDs, vService.ContainerIDs...)
		}

	} else if len(query) == 2 {

		// Find Deployment
		var dep *types.Deployment
		state.VortexDeployments.Mu.RLock()
		for _, deploy := range state.VortexDeployments.List {
			if deploy.Version == query[0] {
				dep = deploy
				break
			}
		}
		state.VortexDeployments.Mu.RUnlock()
		if dep == nil {
			return &pb.BoolResponse{Success: false}, errors.New("deployment not found")
		}

		// Removing service from state
		rmService := utils.PopServiceByName(state.VortexServices, query[1])
		if rmService == nil {
			return &pb.BoolResponse{Success: false}, errors.New("service not found")
		}

		// Removing containers from VortexContainers
		utils.RemoveContainerConfigsByService(state.VortexContainers, rmService)
		rmContainerIDs = append(rmContainerIDs, rmService.ContainerIDs...)

	} else {
		return &pb.BoolResponse{Success: false}, errors.New("invalid deployment or service name recieved")
	}

	// Delete the containers from docker sdk
	// Make sure the containers are delete after removing them from state variables as othervise tracker will try to restart them
	for _, id := range rmContainerIDs {
		docker.DeleteContainer(id)
	}

	return &pb.BoolResponse{Success: true}, nil
}

func (s *server) Show(ctx context.Context, body *pb.NameHolder) (*pb.ResponseBody, error) {
	query := body.GetName()
	if query == "all" {
		data, err := json.Marshal(state.GetState())
		if err != nil {
			return nil, err
		}
		return &pb.ResponseBody{Data: data}, nil
	}
	return nil, errors.New("invalid command recieved")
}

func (s *server) Down(ctx context.Context, body *pb.NameHolder) (*pb.NameHolder, error) {
	msg := body.GetName()
	if msg != "kill" {
		return nil, errors.New("invalid command recieved")
	}
	go func() {
		timer := time.NewTimer(3 * time.Second)
		<-timer.C
		os.Exit(0)
	}()
	return &pb.NameHolder{Name: "Stopping Vortex...Bye"}, nil
}

func StartGrpcServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer(
	// Can define grpc interceptors(middlewares) here
	// grpc.StreamInterceptor(func{}),
	// grpc.UnaryInterceptor(func{}),
	)
	pb.RegisterContainerFactoryServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
