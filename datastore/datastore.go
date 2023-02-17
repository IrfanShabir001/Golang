package datastore

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB interface {
	GetById(ctx context.Context, collectionName string, filter interface{}, dto interface{}) error
	Save(ctx context.Context, collectionName string, dto interface{}) error
	SaveMany(ctx context.Context, collectionName string, dtos []interface{}) error
	Update(ctx context.Context, collectionName string, filter, dto interface{}) error
	GetByFilter(ctx context.Context, collectionName string, filter interface{}, opt *options.FindOptions, dto interface{}) ([]byte, error)
	Delete(ctx context.Context, collectionName string, filter interface{}) error
	DeleteMany(ctx context.Context, collectionName string, filter interface{}) error
}