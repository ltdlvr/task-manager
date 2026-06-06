package model

type BoardRole string

type BoardMember struct {
	BoardID uint64
	UserID  uint64
	Role    BoardRole
}

const (
	BoardRoleOwner  BoardRole = "owner"
	BoardRoleMember BoardRole = "member"
)
