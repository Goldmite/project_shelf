package services

import (
	"database/sql"
	"fmt"

	"github.com/Goldmite/project_shelf/internal/models"
	"github.com/Goldmite/project_shelf/internal/models/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	database *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{database: db}
}

func CheckHash(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}

func (userService *UserService) CreateUser(r models.CreateUserRequest) (*string, error) {
	u, err := models.NewUserFromRequest(r)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password %w", err)
	}
	query := "INSERT INTO users (id, name, email, pw_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = userService.database.Exec(query, u.ID, u.Name, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user %w", err)
	}

	return &u.ID, nil
}

func (userService *UserService) GetUserByIdOrEmail(id, email string) (*dto.GetUserResponse, error) {
	identifyBy := email
	queryBy := "email"
	if email == "" {
		identifyBy = id
		queryBy = "id"
	}
	query := "SELECT id, name, email FROM users WHERE " + queryBy + " = ?"

	var res dto.GetUserResponse
	err := userService.database.QueryRow(query, identifyBy).Scan(&res.ID, &res.Name, &res.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found %w", err)
	}

	return &res, nil
}

func (userService *UserService) GetUserByLogin(email, password string) (*dto.GetUserResponse, error) {
	query := "SELECT id, name, email, pw_hash FROM users WHERE email = ?"

	var res dto.GetUserResponse
	var hashed string
	err := userService.database.QueryRow(query, email).Scan(&res.ID, &res.Name, &res.Email, &hashed)
	if err != nil {
		return nil, fmt.Errorf("user not found %w", err)
	}

	if err := CheckHash(hashed, password); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &res, nil
}

func (userService *UserService) GetUserInvites(userId string) ([]dto.InviteResponse, error) {
	query :=
		"SELECT i.group_id, g.name, i.invited_by " +
			"FROM invitations i " +
			"JOIN users u ON u.email = i.email_to " +
			"JOIN groups g ON g.id = i.group_id " +
			"WHERE u.id = ?"
	rows, err := userService.database.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user invites %w", err)
	}
	defer rows.Close()

	var groupInvites []dto.InviteResponse
	for rows.Next() {
		var inv dto.InviteResponse
		var invitedByUser string
		err := rows.Scan(&inv.GroupId, &inv.GroupName, &invitedByUser)
		if err != nil {
			return nil, err
		}

		inviter, err := userService.GetUserByIdOrEmail(invitedByUser, "")
		if err != nil {
			return nil, err
		}
		inv.InviterName = inviter.Name
		groupInvites = append(groupInvites, inv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupInvites, nil
}

func (userService *UserService) GetGroupMembers(groupId string) ([]dto.GetUserResponse, error) {
	query := "SELECT u.id, u.name, u.email FROM users u JOIN members m ON u.id = m.user_id WHERE m.group_id = ?"

	rows, err := userService.database.Query(query, groupId)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	defer rows.Close()

	var groupMembers []dto.GetUserResponse
	for rows.Next() {
		var member dto.GetUserResponse
		err := rows.Scan(&member.ID, &member.Name, &member.Email)
		if err != nil {
			return nil, err
		}

		groupMembers = append(groupMembers, member)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupMembers, nil
}
