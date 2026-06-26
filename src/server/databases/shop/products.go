package shop

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"shop.loadout.tf/src/server/model"
)

/*
func ProductIDExist(id string) (bool, error) {
	r := productsCollection.FindOne(context.Background(), bson.D{primitive.E{Key: "id", Value: id}})

	err := r.Err()

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
*/

func ProductIDExist(productID string) (bool, error) {
	if shopDb == nil {
		return false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT id FROM products WHERE id = $1;`
	row := shopDb.QueryRow(query, productID)

	var id string
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func InsertProduct(product *model.Product) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	files, err := json.Marshal(&product.Files)
	if err != nil {
		return fmt.Errorf("failed to marshal product.Files: <%w>", err)
	}

	extraData, err := json.Marshal(&product.ExtraData)
	if err != nil {
		return fmt.Errorf("failed to marshal product.ExtraData: <%w>", err)
	}

	options, err := json.Marshal(&product.Options)
	if err != nil {
		return fmt.Errorf("failed to marshal product.Options: <%w>", err)
	}

	_, err = shopDb.Exec(`INSERT INTO products (id, name, product_name, thumbnail_url, description, is_ignored, date_created, date_updated, files, variant_ids, external_id_1, external_id_2, extra_data, options, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	ON CONFLICT (id) DO UPDATE SET
		name = $2,
		product_name = $3,
		thumbnail_url = $4,
		description = $5,
		is_ignored = $6,
		date_created = $7,
		date_updated = $8,
		files = $9,
		variant_ids = $10,
		external_id_1 = $11,
		external_id_2 = $12,
		extra_data = $13,
		options = $14,
		status = $15`,
		product.ID,
		product.Name,
		product.ProductName,
		product.ThumbnailURL,
		product.Description,
		product.IsIgnored,
		time.Now(),
		time.Now(),
		files,
		pq.Array(product.VariantIDs),
		product.ExternalID1,
		product.ExternalID2,
		extraData,
		options,
		product.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to insert mockup task : <%w>", err)
	}

	return nil
}

func GetProduct(productID string) (*model.Product, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT name, product_name, thumbnail_url, description, is_ignored, date_created, date_updated, files, variant_ids, external_id_1, external_id_2, extra_data, options, status FROM products WHERE id = $1;`
	row := shopDb.QueryRow(query, productID)

	var name string
	var productName string
	var thumbnailURL string
	var description string
	var isIgnored bool
	var dateCreated time.Time
	var dateUpdated time.Time
	var files string
	var variantIDs []string
	var externalID1 string
	var externalID2 string
	var extraData string
	var options string
	var status string

	err := row.Scan(&name, &productName, &thumbnailURL, &description, &isIgnored, &dateCreated, &dateUpdated, &files, pq.Array(&variantIDs), &externalID1, &externalID2, &extraData, &options, &status)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in FindProduct: <%w>", err)
	}

	jsonFiles := []model.File{}
	if err = json.Unmarshal([]byte(files), &jsonFiles); err != nil {
		return nil, err
	}

	jsonExtraData := map[string]any{}
	if err = json.Unmarshal([]byte(extraData), &jsonExtraData); err != nil {
		return nil, err
	}

	jsonOptions := []model.Option{}
	if err = json.Unmarshal([]byte(options), &jsonOptions); err != nil {
		return nil, err
	}
	product := model.Product{
		ID:           productID,
		Name:         name,
		ProductName:  productName,
		ThumbnailURL: thumbnailURL,
		Description:  description,
		IsIgnored:    isIgnored,
		DateCreated:  dateCreated,
		DateUpdated:  dateUpdated,
		Files:        jsonFiles,
		VariantIDs:   variantIDs,
		ExternalID1:  externalID1,
		ExternalID2:  externalID2,
		ExtraData:    jsonExtraData,
		Options:      jsonOptions,
		Variants:     []model.Variant{},
		Status:       status,
	}

	return &product, nil
}

func GetProductsByStatus(status string) ([]*model.Product, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT id, name, product_name, thumbnail_url, description, is_ignored, date_created, date_updated, files, variant_ids, external_id_1, external_id_2, extra_data, options FROM products WHERE status = $1;`
	res, err := shopDb.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query "+query+"in GetProductsByStatus: <%w>", err)
	}
	defer res.Close()

	products := make([]*model.Product, 0, 100)
	for res.Next() {
		var id string
		var name string
		var productName string
		var thumbnailURL string
		var description string
		var isIgnored bool
		var dateCreated time.Time
		var dateUpdated time.Time
		var files string
		var variantIDs []string
		var externalID1 string
		var externalID2 string
		var extraData string
		var options string

		err := res.Scan(&id, &name, &productName, &thumbnailURL, &description, &isIgnored, &dateCreated, &dateUpdated, &files, pq.Array(&variantIDs), &externalID1, &externalID2, &extraData, &options)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row in FindProduct: <%w>", err)
		}

		jsonFiles := []model.File{}
		if err = json.Unmarshal([]byte(files), &jsonFiles); err != nil {
			return nil, err
		}

		jsonExtraData := map[string]any{}
		if err = json.Unmarshal([]byte(extraData), &jsonExtraData); err != nil {
			return nil, err
		}

		jsonOptions := []model.Option{}
		if err = json.Unmarshal([]byte(options), &jsonOptions); err != nil {
			return nil, err
		}
		product := model.Product{
			ID:           id,
			Name:         name,
			ProductName:  productName,
			ThumbnailURL: thumbnailURL,
			Description:  description,
			IsIgnored:    isIgnored,
			DateCreated:  dateCreated,
			DateUpdated:  dateUpdated,
			Files:        jsonFiles,
			VariantIDs:   variantIDs,
			ExternalID1:  externalID1,
			ExternalID2:  externalID2,
			ExtraData:    jsonExtraData,
			Options:      jsonOptions,
			Variants:     []model.Variant{},
			Status:       status,
		}
		products = append(products, &product)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row in GetProductsByStatus: <%w>", err)
	}

	return products, nil
}
