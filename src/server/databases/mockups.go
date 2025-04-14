package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/model"
)

func InsertMockupTasks(tasks []*model.MockupTask) error {
	for _, task := range tasks {
		err := InsertMockupTask(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertMockupTask(task *model.MockupTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var t model.MockupTask = *task
	t.Status = "created"
	t.DateCreated = time.Now().Unix()
	t.DateUpdated = time.Now().Unix()

	insertOneResult, err := mockupTasksCollection.InsertOne(ctx, t)

	if err != nil {
		return err
	}

	task.ID = insertOneResult.InsertedID.(primitive.ObjectID)

	return nil
}

func FindMockupTasks() ([]*model.MockupTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "status", Value: "created"}}

	cursor, err := mockupTasksCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := []*model.MockupTask{}

	for cursor.Next(context.TODO()) {
		task := model.MockupTask{}
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		results = append(results, &task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func UpdateMockupTask(task *model.MockupTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	task.DateUpdated = time.Now().Unix()

	filter := bson.D{primitive.E{Key: "_id", Value: task.ID}}
	_, err := mockupTasksCollection.ReplaceOne(ctx, filter, task, opts)
	if err != nil {
		return err
	}

	return nil
}
