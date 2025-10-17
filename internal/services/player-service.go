package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"getswing.app/player-service/internal/models"
	"getswing.app/player-service/internal/repository"
	pb "getswing.app/player-service/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type PlayerServiceImpl struct {
	pb.UnimplementedPlayerServiceServer
	repo repository.PlayerRepository
}

// TODO MOVE THIS INTO UTIL PACKAGE
func UuidToStringValue(id *uuid.UUID) *wrapperspb.StringValue {
	if id == nil {
		return nil
	}
	return wrapperspb.String(id.String())
}

func BoolPointerToBoolValue(b *bool) *wrapperspb.BoolValue {
	if b == nil {
		return nil
	}
	return wrapperspb.Bool(*b)
}

func TimePointerToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func DeletedAtToTimestamp(deletedAt gorm.DeletedAt) *timestamppb.Timestamp {
	if deletedAt.Valid {
		return timestamppb.New(deletedAt.Time)
	}
	return nil
}

func ModelToProtoPlayer(player *models.Player) *pb.Player {
	if player == nil {
		return nil
	}

	return &pb.Player{
		Id:                    player.ID,
		FirstName:             player.FirstName,
		LastName:              player.LastName,
		Username:              player.Username,
		Email:                 player.Email,
		PhoneCountryCode:      player.PhoneCountryCode,
		PhoneNumber:           player.PhoneNumber,
		Nationality:           player.Nationality,
		Gender:                player.Gender,
		ImageFileId:           UuidToStringValue(player.ImageFileID),
		OneSignalId:           UuidToStringValue(player.OneSignalID),
		ReferralCode:          player.ReferralCode,
		Experienced:           player.Experienced,
		SubscriptionId:        player.SubscriptionID,
		PartnerSubscriptions:  player.PartnerSubscriptions,
		IsBetaTester:          BoolPointerToBoolValue(player.IsBetaTester),
		IsOnboarding:          BoolPointerToBoolValue(player.IsOnboarding),
		IsFirstBuy:            BoolPointerToBoolValue(player.IsFirstBuy),
		IsDefaultUsername:     BoolPointerToBoolValue(player.IsDefaultUsername),
		IsFirstBuyGolfCourse:  BoolPointerToBoolValue(player.IsFirstBuyGolfCourse),
		Dob:                   TimePointerToTimestamp(player.Dob),
		CreatedAt:             TimePointerToTimestamp(player.CreatedAt),
		UpdatedAt:             TimePointerToTimestamp(player.UpdatedAt),
		DeletedAt:             DeletedAtToTimestamp(player.DeletedAt),
		Password:              player.Password,
		EmailVerificationCode: player.EmailVerificationCode,
		EmailTemporary:        player.EmailTemporary,
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
