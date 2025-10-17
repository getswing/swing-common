package service

import (
	"context"
	"errors"
	"fmt"

	"getswing.app/player-service/internal/repository"
	pb "getswing.app/player-service/proto"
	"gorm.io/gorm"
)

type PlayerServiceImpl struct {
	pb.UnimplementedPlayerServiceServer
	repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) *PlayerServiceImpl {
	return &PlayerServiceImpl{repo: repo}
}

func (s *PlayerServiceImpl) GetPlayer(ctx context.Context, req *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	fmt.Println("Received request for player:", req.Id)

	// todo: need to fix this
	player, err := s.repo.FindByID(ctx, req.Id)

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	players := []*pb.Player{
		{Id: player.ID, Name: player.FirstName},
	}

	resp := &pb.GetPlayerResponse{
		Players: players,
	}

	return resp, nil
}
