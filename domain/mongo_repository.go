package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const CommentedEntityCollection = "commentedEntities"

type CommentedEntityRepositoryImpl struct {
	Client  *mongo.Client
}

func (r *CommentedEntityRepositoryImpl) FindByEntityTypeAndId(ctx context.Context,  companyId, entityType, entityId string) (*CommentedEntity, error){

	res := r.Client.Database("comments_default").Collection(CommentedEntityCollection).FindOne(ctx, bson.M{"companyId": companyId, "entityType": entityType, "entityId": entityId})
	if res.Err() != nil {
		return nil, res.Err()
	}

	var commentedEntity CommentedEntity
	res.Decode(&commentedEntity)

	return &commentedEntity, nil
}