package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type RedeemCode struct {
	ID        int64
	Code      string
	Type      string
	Value     float64
	Status    string
	UsedBy    *int64
	UsedAt    *time.Time
	Notes     string
	CreatedAt time.Time

	MaxClaims    int
	ClaimedCount int
	GroupID      *int64
	ValidityDays int

	User  *User
	Group *Group
}

func (r *RedeemCode) IsUsed() bool {
	return r.Status == StatusUsed
}

func (r *RedeemCode) CanUse() bool {
	return r.Status == StatusUnused
}

func (r *RedeemCode) IsWelfare() bool {
	return r != nil && r.Type == RedeemTypeWelfare
}

func (r *RedeemCode) RemainingClaims() int {
	if r == nil || r.MaxClaims <= 0 {
		return 0
	}
	remaining := r.MaxClaims - r.ClaimedCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

func GenerateRedeemCode() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
