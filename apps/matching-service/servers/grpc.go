package servers

import (
	"log"
	pb "matching-service/proto"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	questionMatchingClient pb.QuestionMatchingServiceClient
)

func InitGrpcServer() *grpc.ClientConn {
	questionServiceAddr := os.Getenv("QUESTION_SERVICE_GRPC_URL")
	// Dial the server
	conn, err := grpc.NewClient(questionServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to %v: %v", questionServiceAddr, err)
	} else {
		log.Println("Connected to Grpc server at %v", questionServiceAddr)
	}

	// Create a new client for the ExampleService
	questionMatchingClient = pb.NewQuestionMatchingServiceClient(conn)

	return conn
}

func GetGrpcClient() pb.QuestionMatchingServiceClient {
	return questionMatchingClient
}
