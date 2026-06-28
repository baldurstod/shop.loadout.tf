package shop

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"shop.loadout.tf/src/server/encryption"
	"shop.loadout.tf/src/server/model"
)

var enveloped = encryption.NewEnveloped(encryption.Kms{})

func CreateOrder() (*model.Order, error) {
	var id string
	ok := false
	for range maxCreationAttempts {
		id = createRandID()
		exist, err := orderIDExist(id)
		if err != nil {
			return nil, err
		}

		if !exist {
			ok = true
			break
		}
	}

	if !ok {
		return nil, errors.New("unable to create an id")
	}

	order := model.NewOrder()
	order.ID = id

	if err := insertOrder(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func insertOrder(order *model.Order) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	dekPlain, dekCipher, err := enveloped.GenerateDek(context.Background())
	if err != nil {
		return fmt.Errorf("failed to generate DEK: <%w>", err)
	}

	shippingAddress, err := json.Marshal(&order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal order.ShippingAddress: <%w>", err)
	}

	shippingAddressEncryptedField, err := encryption.EncryptAES(shippingAddress, dekPlain)
	if err != nil {
		return fmt.Errorf("failed to encrypt shipping address: <%w>", err)
	}

	billingAddress, err := json.Marshal(&order.BillingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal order.BillingAddress: <%w>", err)
	}

	billingAddressEncryptedField, err := encryption.EncryptAES(billingAddress, dekPlain)
	if err != nil {
		return fmt.Errorf("failed to encrypt billing address: <%w>", err)
	}

	items, err := json.Marshal(&order.Items)
	if err != nil {
		return fmt.Errorf("failed to marshal order.Items: <%w>", err)
	}

	shippingInfos, err := json.Marshal(&order.ShippingInfos)
	if err != nil {
		return fmt.Errorf("failed to marshal order.ShippingInfos: <%w>", err)
	}

	taxInfo, err := json.Marshal(&order.TaxInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal order.TaxInfo: <%w>", err)
	}

	_, err = shopDb.Exec(`INSERT INTO orders (id, currency, shipping_address, billing_address, same_billing_address, items, shipping_infos, tax_info, shipping_method, printful_order_id, paypal_order_id, dek, status, date_created, date_updated)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			ON CONFLICT (id) DO UPDATE SET
			currency = $2,
			shipping_address = $3,
			billing_address = $4,
			same_billing_address = $5,
			items = $6,
			shipping_infos = $7,
			tax_info = $8,
			shipping_method = $9,
			printful_order_id = $10,
			paypal_order_id = $11,
			dek = $12,
			status = $13,
			date_created = $14,
			date_updated = $15`,
		order.ID,
		order.Currency,
		shippingAddressEncryptedField,
		billingAddressEncryptedField,
		order.SameBillingAddress,
		items,
		shippingInfos,
		taxInfo,
		order.ShippingMethod,
		order.PrintfulOrderID,
		order.PaypalOrderID,
		dekCipher,
		order.Status,
		order.DateCreated,
		order.DateUpdated,
	)

	if err != nil {
		return fmt.Errorf("failed to insert order: <%w>", err)
	}

	return nil
}

func orderIDExist(orderId string) (bool, error) {
	if shopDb == nil {
		return false, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	query := `SELECT id FROM orders WHERE id = $1;`
	row := shopDb.QueryRow(query, orderId)

	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func UpdateOrder(order *model.Order) error {
	order.DateUpdated = time.Now()

	err := insertOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func GetOrder(orderId string) (*model.Order, error) {
	query := `SELECT id, currency, shipping_address, billing_address, same_billing_address, items, shipping_infos, tax_info, shipping_method, printful_order_id, paypal_order_id, dek, status, date_created, date_updated FROM orders WHERE id = $1;`
	return getOrder(query, orderId)
}

func getOrder(query string, args ...any) (*model.Order, error) {
	if shopDb == nil {
		return nil, errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	row := shopDb.QueryRow(query, args...)

	var id string
	var currency string
	var encryptedShippingAddress string
	var encryptedBillingAddress string
	var sameBillingAddress bool
	var items string
	var shippingInfos string
	var taxInfo string
	var percentDiscount primitive.Decimal128
	var priceDiscount primitive.Decimal128
	var shippingMethod string
	var printfulOrderID string
	var paypalOrderID string
	var encryptedDek string
	var status string
	var dateCreated time.Time
	var dateUpdated time.Time

	err := row.Scan(&id, &currency, &encryptedShippingAddress, &encryptedBillingAddress, &sameBillingAddress, &items, &shippingInfos, &taxInfo, &shippingMethod, &printfulOrderID, &paypalOrderID, &encryptedDek, &status, &dateCreated, &dateUpdated)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in GetOrder: <%w>", err)
	}

	dek, err := enveloped.DecryptDek(context.Background(), []byte(encryptedDek))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt DEK: <%w>", err)
	}

	shippingAddressDecryptedField, err := encryption.DecryptAES([]byte(encryptedShippingAddress), dek)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt shipping address: <%w>", err)
	}

	shippingAddress := model.Address{}
	if err = json.Unmarshal(shippingAddressDecryptedField, &shippingAddress); err != nil {
		return nil, err
	}

	billingAddressDecryptedField, err := encryption.DecryptAES([]byte(encryptedBillingAddress), dek)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt billing address: <%w>", err)
	}

	billingAddress := model.Address{}
	if err = json.Unmarshal(billingAddressDecryptedField, &billingAddress); err != nil {
		return nil, err
	}

	jsonItems := []model.OrderItem{}
	if err = json.Unmarshal([]byte(items), &jsonItems); err != nil {
		return nil, err
	}

	jsonShippingInfos := []printfulmodel.ShippingRate{}
	if err = json.Unmarshal([]byte(shippingInfos), &jsonShippingInfos); err != nil {
		return nil, err
	}

	jsonTaxInfo := model.TaxInfo{}
	if err = json.Unmarshal([]byte(taxInfo), &jsonTaxInfo); err != nil {
		return nil, err
	}

	order := model.Order{
		ID:                 id,
		Currency:           currency,
		ShippingAddress:    shippingAddress,
		BillingAddress:     billingAddress,
		SameBillingAddress: sameBillingAddress,
		Items:              jsonItems,
		ShippingInfos:      jsonShippingInfos,
		TaxInfo:            jsonTaxInfo,
		PercentDiscount:    percentDiscount,
		PriceDiscount:      priceDiscount,
		ShippingMethod:     shippingMethod,
		PrintfulOrderID:    printfulOrderID,
		PaypalOrderID:      paypalOrderID,
		Status:             status,
		DateCreated:        dateCreated,
		DateUpdated:        dateUpdated,
	}

	return &order, nil

}

func GetOrderByPaypalID(paypalId string) (*model.Order, error) {
	query := `SELECT id, currency, shipping_address, billing_address, same_billing_address, items, shipping_infos, tax_info, shipping_method, printful_order_id, paypal_order_id, dek, status, date_created, date_updated FROM orders WHERE paypal_order_id = $1;`
	return getOrder(query, paypalId)
}
