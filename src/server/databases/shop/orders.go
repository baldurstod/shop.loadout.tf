package shop

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/shopspring/decimal"
	"shop.loadout.tf/src/server/model"
)

var enveloped = newEnveloped(kms{})

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
	now := time.Now()
	order.DateCreated = now
	order.DateUpdated = now
	order.Status = "created"

	if err := insertOrder(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func insertOrder(order *model.Order) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	shippingAddress, err := json.Marshal(&order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal order.ShippingAddress: <%w>", err)
	}

	shippingAddressEncryptedField, shippingAddressEncryptedKey, shippingAddressEncryptedKekId, err := enveloped.EncryptAES(shippingAddress)
	if err != nil {
		return fmt.Errorf("failed to encrypt shipping address: <%w>", err)
	}

	billingAddress, err := json.Marshal(&order.BillingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal order.BillingAddress: <%w>", err)
	}

	billingAddressEncryptedField, billingAddressEncryptedKey, billingAddressEncryptedKekId, err := enveloped.EncryptAES(billingAddress)
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

	_, err = shopDb.Exec(`INSERT INTO orders (id, currency, shipping_address, shipping_address_dek, shipping_address_kek, billing_address, billing_address_dek, billing_address_kek, same_billing_address, items, shipping_infos, tax_info, shipping_method, items_price, discount_price, shipping_price, tax_price, total_price	, 	printful_order_id, paypal_order_id, status, date_created, date_updated)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
			ON CONFLICT (id) DO UPDATE SET
			currency = $2,
			shipping_address = $3,
			shipping_address_dek = $4,
			shipping_address_kek = $5,
			billing_address = $6,
			billing_address_dek = $7,
			billing_address_kek = $8,
			same_billing_address = $9,
			items = $10,
			shipping_infos = $11,
			tax_info = $12,
			shipping_method = $13,
			items_price = $14,
			discount_price = $15,
			shipping_price = $16,
			tax_price = $17,
			total_price = $18,
			printful_order_id = $19,
			paypal_order_id = $20,
			status = $21,
			date_created = $22,
			date_updated = $23`,
		order.ID,
		order.Currency,
		shippingAddressEncryptedField,
		shippingAddressEncryptedKey,
		shippingAddressEncryptedKekId,
		billingAddressEncryptedField,
		billingAddressEncryptedKey,
		billingAddressEncryptedKekId,
		order.SameBillingAddress,
		items,
		shippingInfos,
		taxInfo,
		order.ShippingMethod,
		order.ItemsPrice,
		order.DiscountPrice,
		order.ShippingPrice,
		order.TaxPrice,
		order.TotalPrice,
		order.PrintfulOrderID,
		order.PaypalOrderID,
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

type UpdateOrderFields struct {
	Currency           bool
	ShippingAddress    bool
	BillingAddress     bool
	SameBillingAddress bool
	Items              bool
	ShippingInfos      bool
	TaxInfo            bool
	PercentDiscount    bool
	PriceDiscount      bool
	ShippingMethod     bool
	ItemsPrice         bool
	DiscountPrice      bool
	ShippingPrice      bool
	TaxPrice           bool
	TotalPrice         bool
	PrintfulOrderID    bool
	PaypalOrderID      bool
	Status             bool
}

func UpdateOrder(order *model.Order, fields UpdateOrderFields) error {
	if shopDb == nil {
		return errors.New("database is not initialized. Did you forgot to init postgre ?")
	}

	queryString := make([]string, 0)
	queryParams := []any{order.ID, time.Now()}

	v := reflect.ValueOf(fields)
	typeOfS := v.Type()

	// Using reflection to list UpdateOrderFields fields
	for i := 0; i < v.NumField(); i++ {
		name := typeOfS.Field(i).Name
		value := v.Field(i).Bool()
		if !value {
			continue
		}

		addSetStatement := func(column string, value any) {
			param := "$" + strconv.Itoa(len(queryParams)+1)

			queryString = append(queryString, column+" = "+param)
			queryParams = append(queryParams, value)
		}

		switch name {
		case "Currency":
			addSetStatement("currency", order.Currency)
		case "ShippingAddress":
			shippingAddress, err := json.Marshal(&order.ShippingAddress)
			if err != nil {
				return fmt.Errorf("failed to marshal order.ShippingAddress: <%w>", err)
			}

			shippingAddressEncryptedField, shippingAddressEncryptedKey, shippingAddressKekId, err := enveloped.EncryptAES(shippingAddress)
			if err != nil {
				return fmt.Errorf("failed to encrypt shipping address: <%w>", err)
			}

			addSetStatement("shipping_address", shippingAddressEncryptedField)
			addSetStatement("shipping_address_dek", shippingAddressEncryptedKey)
			addSetStatement("shipping_address_kek", shippingAddressKekId)
		case "BillingAddress":
			billingAddress, err := json.Marshal(&order.BillingAddress)
			if err != nil {
				return fmt.Errorf("failed to marshal order.BillingAddress: <%w>", err)
			}

			billingAddressEncryptedField, billingAddressEncryptedKey, billingAddressKekId, err := enveloped.EncryptAES(billingAddress)
			if err != nil {
				return fmt.Errorf("failed to encrypt billing address: <%w>", err)
			}

			addSetStatement("billing_address", billingAddressEncryptedField)
			addSetStatement("billing_address_dek", billingAddressEncryptedKey)
			addSetStatement("billing_address_kek", billingAddressKekId)
		case "SameBillingAddress":
			addSetStatement("same_billing_address", order.SameBillingAddress)
		case "Items":

			items, err := json.Marshal(&order.Items)
			if err != nil {
				return fmt.Errorf("failed to marshal order.Items: <%w>", err)
			}

			addSetStatement("items", items)
		case "ShippingInfos":

			shippingInfos, err := json.Marshal(&order.ShippingInfos)
			if err != nil {
				return fmt.Errorf("failed to marshal order.ShippingInfos: <%w>", err)
			}

			addSetStatement("shipping_infos", shippingInfos)
		case "TaxInfo":

			taxInfo, err := json.Marshal(&order.TaxInfo)
			if err != nil {
				return fmt.Errorf("failed to marshal order.TaxInfo: <%w>", err)
			}

			addSetStatement("tax_info", taxInfo)
		case "ShippingMethod":
			addSetStatement("shipping_method", order.ShippingMethod)
		case "ItemsPrice":
			addSetStatement("items_price", order.ItemsPrice)
		case "DiscountPrice":
			addSetStatement("discount_price", order.DiscountPrice)
		case "ShippingPrice":
			addSetStatement("shipping_price", order.ShippingPrice)
		case "TaxPrice":
			addSetStatement("tax_price", order.TaxPrice)
		case "TotalPrice":
			addSetStatement("total_price", order.TotalPrice)
		case "PaypalOrderID":
			addSetStatement("paypal_order_id", order.PaypalOrderID)
		case "Status":
			addSetStatement("status", order.Status)
		default:
			return errors.New("missing field in UpdateOrder " + name)
		}
	}

	if len(queryString) == 0 {
		return errors.New("failed to update order: no field selected for update")
	}

	query := `UPDATE orders SET date_updated = $2,` + strings.Join(queryString, ",") + ` WHERE id = $1;`
	res, err := shopDb.Exec(query, queryParams...)
	if err != nil {
		return fmt.Errorf("failed to update order: <%w>", err)
	}

	if rows, err := res.RowsAffected(); rows != 1 || err != nil {
		if err != nil {
			return fmt.Errorf("failed to get rows affected %s: <%w>", order.ID, err)
		} else {
			return fmt.Errorf("failed to update user %s: %d rows affected, expected 1", order.ID, rows)
		}
	}

	return nil
}

func GetOrder(orderId string) (*model.Order, error) {
	query := `SELECT id, currency, shipping_address, shipping_address_dek, shipping_address_kek, billing_address, billing_address_dek, billing_address_kek, same_billing_address, items, shipping_infos, tax_info, shipping_method, printful_order_id, paypal_order_id, status, date_created, date_updated FROM orders WHERE id = $1;`
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
	var encryptedShippingAddressDek string
	var encryptedShippingAddressKek int64
	var encryptedBillingAddress string
	var encryptedBillingAddressDek string
	var encryptedBillingAddressKek int64
	var sameBillingAddress bool
	var items string
	var shippingInfos string
	var taxInfo string
	var percentDiscount decimal.Decimal
	var priceDiscount decimal.Decimal
	var shippingMethod string
	var printfulOrderID string
	var paypalOrderID string
	//var encryptedDek string
	var status string
	var dateCreated time.Time
	var dateUpdated time.Time

	err := row.Scan(&id, &currency, &encryptedShippingAddress, &encryptedShippingAddressDek, &encryptedShippingAddressKek, &encryptedBillingAddress, &encryptedBillingAddressDek, &encryptedBillingAddressKek, &sameBillingAddress, &items, &shippingInfos, &taxInfo, &shippingMethod, &printfulOrderID, &paypalOrderID, &status, &dateCreated, &dateUpdated)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row in GetOrder: <%w>", err)
	}

	/*
		plainShippingAddressDek, err := enveloped.DecryptDek(context.Background(), []byte(encryptedShippingAddressDek))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt shipping_address_dek: <%w>", err)
		}
	*/

	shippingAddressDecryptedField, err := enveloped.DecryptAES([]byte(encryptedShippingAddress), []byte(encryptedShippingAddressDek), encryptedShippingAddressKek)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt shipping address: <%w>", err)
	}

	shippingAddress := model.Address{}
	if err = json.Unmarshal(shippingAddressDecryptedField, &shippingAddress); err != nil {
		return nil, err
	}

	/*
		plainBillingAddressDek, err := enveloped.DecryptAES(context.Background(), []byte(encryptedBillingAddressDek))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt billing_address_dek: <%w>", err)
		}
	*/

	billingAddressDecryptedField, err := enveloped.DecryptAES([]byte(encryptedBillingAddress), []byte(encryptedBillingAddressDek), encryptedBillingAddressKek)
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
