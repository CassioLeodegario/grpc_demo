package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/cassioleodegario/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status:", stream)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "c1",
			Name:  "Cassio1",
			Email: "c@c.com",
		},
		{
			Id:    "c2",
			Name:  "Cassio2",
			Email: "c@c.com",
		},
		{
			Id:    "c3",
			Name:  "Cassio3",
			Email: "c@c.com",
		},
		{
			Id:    "c4",
			Name:  "Cassio4",
			Email: "c@c.com",
		},
		{
			Id:    "c5",
			Name:  "Cassio5",
			Email: "c@c.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response %v", err)
	}

	fmt.Println(res)

}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "c1",
			Name:  "Cassio1",
			Email: "c@c.com",
		},
		{
			Id:    "c2",
			Name:  "Cassio2",
			Email: "c@c.com",
		},
		{
			Id:    "c3",
			Name:  "Cassio3",
			Email: "c@c.com",
		},
		{
			Id:    "c4",
			Name:  "Cassio4",
			Email: "c@c.com",
		},
		{
			Id:    "c5",
			Name:  "Cassio5",
			Email: "c@c.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status %v:\n\n", res.GetUser().GetName(), res.Status)
		}
		close(wait)
	}()

	<-wait
}
