package handlers

import (
	pb "question-service/proto"

	"cloud.google.com/go/firestore"
)

type Service struct {
	Client *firestore.Client
}

type GrpcServer struct {
	pb.UnimplementedQuestionMatchingServiceServer // Embed the unimplemented service
	Client                                        *firestore.Client
}
