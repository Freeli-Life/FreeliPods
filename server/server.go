package server

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"log"

	fc "FreeliPods/crypto"
	"FreeliPods/database"
	pb "FreeliPods/podServer"
)

type PodServiceServer struct {
	pb.UnimplementedPodServiceServer
	DB         *database.Store
	PrivateKey crypto.PrivateKey
	Domain     string
}

func (s *PodServiceServer) RegisterUsername(ctx context.Context, req *pb.RegisterUsernameRequest) (*pb.RegisterUsernameResponse, error) {
	log.Printf("Received registration request for username: %s", req.Username)

	if len(req.Salt) != 16 {
		return nil, errors.New("salt must be 16 bytes")
	}
	if len(req.PublicSigningKey) != 32 {
		return nil, errors.New("signing key must be 32 bytes")
	}
	if len(req.PublicEncryptionKey) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}

	exists, err := s.DB.UserExists(req.Username)
	if err != nil {
		log.Printf("Error checking database: %v", err)
		return nil, errors.New("internal server error")
	}
	if exists {
		log.Printf("Username %s already exists.", req.Username)
		return nil, fmt.Errorf("username '%s' is already taken", req.Username)
	}

	signature, err := fc.SignRegistrationData(
		s.PrivateKey,
		s.Domain,
		req.Username,
		req.Salt,
		req.PublicSigningKey,
		req.PublicEncryptionKey,
	)
	if err != nil {
		log.Printf("Error signing data: %v", err)
		return nil, errors.New("internal server error during signing")
	}

	err = s.DB.AddUser(req.Username, req.Salt, req.PublicSigningKey, req.PublicEncryptionKey)
	if err != nil {
		log.Printf("Error adding user to database: %v", err)
		return nil, errors.New("internal server error saving user")
	}

	log.Printf("Successfully registered username: %s", req.Username)

	return &pb.RegisterUsernameResponse{ServerSignature: signature}, nil
}