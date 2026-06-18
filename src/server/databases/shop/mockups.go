package shop

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/lib/pq"
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

/*
func InsertMockupTask(task *model.MockupTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
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
*/

func InsertMockupTask(task *model.MockupTask) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	template, err := json.Marshal(&task.Template)
	if err != nil {
		return fmt.Errorf("failed to marshal task.Template: <%w>", err)
	}

	_, err = shopDb.Exec(`INSERT INTO mockup_tasks (product_ids, source_image, template, date_created, date_updated, status)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		pq.Array(task.ProductIDs),
		task.SourceImage,
		template,
		task.DateCreated,
		task.DateUpdated,
		task.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to insert mockup task : <%w>", err)
	}

	return nil
}

/*
func FindMockupTasks() ([]*model.MockupTask, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
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
*/

func FindMockupTasks() ([]*model.MockupTask, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT id, product_ids, source_image, template, date_created, date_updated, status FROM mockup_tasks WHERE status = 'created';`
	res, err := shopDb.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query "+query+"in FindMockupTasks: <%w>", err)
	}
	defer res.Close()

	tasks := make([]*model.MockupTask, 0, 20)
	for res.Next() {
		var id int64
		var productIDs []string
		var sourceImage string
		var template string
		var status string
		var dateCreated time.Time
		var dateUpdated time.Time

		err = res.Scan(&id, pq.Array(&productIDs), &sourceImage, &template, &dateCreated, &dateUpdated, &status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in FindMockupTasks: <%w>", err)
		}

		mockupTemplates := printfulmodel.MockupTemplates{}
		if err = json.Unmarshal([]byte(template), &mockupTemplates); err != nil {
			return nil, err
		}

		task := model.MockupTask{ID: id, ProductIDs: productIDs, SourceImage: sourceImage, Template: &mockupTemplates, Status: status, DateCreated: dateCreated, DateUpdated: dateUpdated}

		tasks = append(tasks, &task)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row in FindMockupTasks: <%w>", err)
	}

	return tasks, nil
}

/*
func UpdateMockupTask(task *model.MockupTask) error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoTimeout)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	task.DateUpdated = time.Now()

	filter := bson.D{primitive.E{Key: "_id", Value: task.ID}}
	_, err := mockupTasksCollection.ReplaceOne(ctx, filter, task, opts)
	if err != nil {
		return err
	}

	return nil
}
*/

func UpdateMockupTask(task *model.MockupTask) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	_, err := shopDb.Exec(`UPDATE mockup_tasks SET date_updated = $2, status = $3 WHERE id = $1`,
		task.ID,
		task.DateUpdated,
		task.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to update mockup task : <%w>", err)
	}

	return nil
}
