package repo

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type BoardMembers struct{}

func NewBoardMembers() *BoardMembers {
	return &BoardMembers{}
}

func (r *BoardMembers) Create(ctx context.Context, client db.DB, bm *model.BoardMember) error {
	_, err := client.ExecContext(ctx, `
		INSERT INTO board_members (board_id, user_id, role)
		VALUES ($1, $2, $3)
	`, bm.BoardID, bm.UserID, bm.Role)

	return db.MapError(err)
}

func (r *BoardMembers) GetRole(ctx context.Context, client db.DB, boardID, userID uint64) (model.BoardRole, error) {
	row := client.QueryRowContext(ctx, `
		SELECT role FROM board_members 
		WHERE board_id = $1 AND user_id = $2
	`, boardID, userID)

	var role model.BoardRole
	if err := row.Scan(&role); err != nil {
		return "", db.MapError(err)
	}

	return role, nil
}
