package domain

import "context"

type CommentedEntityRepository interface {

	 FindByEntityTypeAndId(ctx context.Context, companyId, entityType, entityId string) (*CommentedEntity, error)
}