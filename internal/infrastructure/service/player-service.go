package service

import (
	"context"
	"errors"
	"fmt"

	"getswing.app/player-service/internal/app/models"
	"getswing.app/player-service/internal/app/repository"
	"getswing.app/player-service/internal/pkg/utils"
	pb "getswing.app/player-service/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type PlayerServiceImpl struct {
	pb.UnimplementedPlayerServiceServer
	repo repository.PlayerRepository
}

func ModelToProtoPlayer(player *models.Player) *pb.Player {
	if player == nil {
		return nil
	}

	return &pb.Player{
		Id:                   player.ID,
		FirstName:            player.FirstName,
		LastName:             player.LastName,
		Username:             player.Username,
		Email:                player.Email,
		Dialcode:             player.DialCode,
		PhoneNumber:          player.PhoneNumber,
		Nationality:          player.Nationality,
		Gender:               string(player.Gender),
		ProfilePath:          wrapperspb.String(player.ProfilePath),
		Experienced:          player.Experienced,
		IsBetaTester:         utils.BoolPointerToBoolValue(player.IsBetaTester),
		IsOnboarding:         utils.BoolPointerToBoolValue(player.IsOnboarding),
		IsFirstBuy:           utils.BoolPointerToBoolValue(player.IsFirstBuy),
		IsDefaultUsername:    utils.BoolPointerToBoolValue(player.IsDefaultUsername),
		IsFirstBuyGolfCourse: utils.BoolPointerToBoolValue(player.IsFirstBuyGolfCourse),
		Dob:                  utils.TimePointerToTimestamp(player.Dob),
		CreatedAt:            utils.TimePointerToTimestamp(player.CreatedAt),
		UpdatedAt:            utils.TimePointerToTimestamp(player.UpdatedAt),
	}
}

// END OF TODO

func NewPlayerService(repo repository.PlayerRepository) *PlayerServiceImpl {
	return &PlayerServiceImpl{repo: repo}
}

func (s *PlayerServiceImpl) GetPlayer(ctx context.Context, req *pb.GetPlayerRequest) (*pb.GetPlayerResponse, error) {
	fmt.Println("Received request for player:", req.Key)
	var player *models.Player
	var err error

	if req == nil || req.FindBy == "" {
		return  nil, fmt.Errorf("param is required")
	}

	if req.FindBy == "email" {
		player, err = s.repo.FindByEmail(ctx, req.Key)
	}

	if req.FindBy == "id" {
		player, err = s.repo.FindByID(ctx, req.Key)

	}

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	players := []*pb.Player{
		ModelToProtoPlayer(player),
	}

	resp := &pb.GetPlayerResponse{
		Players: players,
	}

	return resp, nil
}

func (s *PlayerServiceImpl) GetPlayerByPhone(ctx context.Context, req *pb.GetPlayerByPhoneRequest) (*pb.GetPlayerResponse, error) {
	player, err := s.repo.FindByPhone(ctx, req.PhoneNumber, req.PhoneCountryCode)

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	players := []*pb.Player{
		ModelToProtoPlayer(player),
	}

	resp := &pb.GetPlayerResponse{
		Players: players,
	}

	return resp, nil
}
