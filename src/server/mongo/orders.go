package mongo

import (
	"context"
	"time"

	"github.com/baldurstod/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shop.loadout.tf/src/server/model"
)

func CreateOrder() (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	order := model.NewOrder()
	order.ID = randstr.String(12, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	_, err := ordersCollection.InsertOne(ctx, order)
	if mongo.IsDuplicateKeyError(err) {
		return CreateOrder() // TODO: improve that
	}

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func UpdateOrder(order *model.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	order.DateUpdated = time.Now().Unix()

	encryptedOrder, err := encryptOrder(order)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "id", Value: order.ID}}
	_, err = ordersCollection.ReplaceOne(ctx, filter, encryptedOrder, opts)
	if err != nil {
		return err
	}

	return nil
}

func FindOrder(orderID string) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "id", Value: orderID}}

	r := ordersCollection2.FindOne(ctx, filter)

	order := model.Order{}
	if err := r.Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func FindOrderByPaypalID(paypalID string) (*model.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "paypal_order_id", Value: paypalID}}

	r := ordersCollection.FindOne(ctx, filter)

	order := model.Order{}
	if err := r.Decode(&order); err != nil {
		return nil, err
	}

	return &order, nil
}

func encryptOrder(order *model.Order) (*bson.M, error) {
	shippingAddressEncryptedField, err := encryptAddress(&order.ShippingAddress)
	if err != nil {
		return nil, err
	}

	billingAddressEncryptedField, err := encryptAddress(&order.BillingAddress)
	if err != nil {
		return nil, err
	}

	return &bson.M{
		"id":                   order.ID,
		"currency":             order.Currency,
		"date_created":         order.DateCreated,
		"date_updated":         order.DateUpdated,
		"shipping_address":     shippingAddressEncryptedField,
		"billing_address":      billingAddressEncryptedField,
		"same_billing_address": order.SameBillingAddress,
		"items":                order.Items,
		"shipping_infos":       order.ShippingInfos,
		"tax_info":             order.TaxInfo,
		"shipping_method":      order.ShippingMethod,
		"printful_order_id":    order.PrintfulOrderID,
		"paypal_order_id":      order.PaypalOrderID,
		"status":               order.Status,
	}, nil
}

func encryptAddress(address *model.Address) (*bson.M, error) {
	/*
		nameRawValueType, nameRawValueData, err := bson.MarshalValue("Greg")
		if err != nil {
			panic(err)
		}
	*/

	firstNameEncryptedField, err := encryptString(address.FirstName)
	if err != nil {
		return nil, err
	}

	return &bson.M{
		"first_name":   firstNameEncryptedField,
		"last_name":    address.LastName,
		"organization": address.Organization,
		"address1":     address.Address1,
		"address2":     address.Address2,
		"city":         address.City,
		"state_code":   address.StateCode,
		"state_name":   address.StateName,
		"country_code": address.CountryCode,
		"country_name": address.CountryName,
		"postal_code":  address.PostalCode,
		"phone":        address.Phone,
		"email":        address.Email,
		"tax_number":   address.TaxNumber,
	}, nil
}

func encryptString(s string) (*primitive.Binary, error) {
	nameRawValueType, nameRawValueData, err := bson.MarshalValue(s)
	if err != nil {
		return nil, err
	}
	nameRawValue := bson.RawValue{Type: nameRawValueType, Value: nameRawValueData}
	nameEncryptionOpts := options.Encrypt().
		SetAlgorithm("AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic").
		SetKeyID(dataKeyId)

	encryptedField, err := clientEnc.Encrypt(
		context.TODO(),
		nameRawValue,
		nameEncryptionOpts)
	if err != nil {
		return nil, err
	}

	return &encryptedField, nil
}
