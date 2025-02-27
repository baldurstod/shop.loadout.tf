package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"log"
	"math"
	"slices"
	"strconv"
	"strings"

	printfulApiModel "github.com/baldurstod/go-printful-api-model"
	"github.com/baldurstod/go-printful-api-model/responses"
	"github.com/baldurstod/go-printful-api-model/schemas"
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/gin-gonic/gin"
	"github.com/greatcloak/decimal"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/model/requests"
	"shop.loadout.tf/src/server/mongo"
)

func apiCreateProduct(c *gin.Context, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}
	//log.Println(params)
	//createProduct := params["product"].(requests.CreateProductRequest)

	createProductRequest := requests.CreateProductRequest{}
	err := mapstructure.Decode(params["product"], &createProductRequest)
	if err != nil {
		log.Println(err)
		return errors.New("error while reading params")
	}

	err = checkParams(&createProductRequest)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("invalid params: %w", err)
	}

	log.Println( /*createProductRequest.Name, */ createProductRequest.VariantID)
	products, err := createProduct(&createProductRequest)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error while creating product: %w", err)
	}

	jsonSuccess(c, map[string]interface{}{"products": products})

	return nil
}

func checkParams(request *requests.CreateProductRequest) error {
	if request.ProductID == 0 {
		return errors.New("invalid product id")
	}

	if request.VariantID == 0 {
		return errors.New("invalid variant id")
	}

	if len(request.Placements) == 0 {
		return errors.New("product have no placements")
	}

	for i, placement := range request.Placements {
		if placement.Placement == "" {
			return fmt.Errorf("placemeny %d has no id", i)
		}

		if placement.Technique == "" {
			return fmt.Errorf("placemeny %d has no technique", i)
		}

		if placement.Image == "" {
			return fmt.Errorf("placemeny %d has no image", i)
		}

		if placement.Orientation == "" {
			return fmt.Errorf("placemeny %d has no orientation", i)
		}
	}

	_, variants, err := getPrintfulProduct(request.ProductID)
	if err != nil {
		return fmt.Errorf("product %d not found", request.ProductID)
	}

	idx := slices.IndexFunc(variants, func(v printfulmodel.Variant) bool { return v.ID == request.VariantID })
	if idx == -1 {
		return fmt.Errorf("variant %d not found", request.VariantID)
	}

	/*
		templates, err := getPrintfulMockupTemplates(request.ProductID)
		if err != nil {
			return errors.New("unable to get product templates")
		}
	*/

	styles, err := getPrintfulStyles(request.ProductID)
	if err != nil {
		return errors.New("unable to get product styles")
	}
	//log.Println(product, variants, styles)

	for i, placement := range request.Placements {
		/*
			idx := slices.IndexFunc(templates, func(t printfulmodel.MockupTemplates) bool {
				if t.Orientation != placement.Orientation ||
					t.Technique != placement.Technique ||
					t.Placement != placement.Placement {
					return false
				}

				idx := slices.IndexFunc(t.CatalogVariantIDs, func(id int) bool { return id == request.VariantID })
				if idx != -1 {
					return true
				}

				return true
			})

			if idx == -1 {
				return fmt.Errorf("template not foundd for placement %d", i)
			}
		*/
		styleIdx := slices.IndexFunc(styles, func(s printfulmodel.MockupStyles) bool {
			if //s.Orientation != placement.Orientation ||
			//TODO: orientation
			s.Technique != placement.Technique ||
				s.Placement != placement.Placement {
				return false
			}

			return true
		})

		if styleIdx == -1 {
			return fmt.Errorf("style not foundd for placement %d", i)
		}

		style := styles[styleIdx]
		overSample := 2.
		styleWidth := int(math.Ceil(style.PrintAreaWidth * float64(style.Dpi) * overSample))
		styleHeight := int(math.Ceil(style.PrintAreaHeight * float64(style.Dpi) * overSample))

		// TODO: check image size
		b64data := placement.Image[strings.IndexByte(placement.Image, ',')+1:] // Remove data:image/png;base64,

		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))
		config, err := png.DecodeConfig(reader)
		if err != nil {
			return errors.New("unable to decode image")
		}

		if config.Width > 20000 || config.Height > 20000 {
			return errors.New("image too large")
		}

		if config.Width < styleWidth || config.Height < styleHeight {
			return fmt.Errorf("invalid image size: %dx%d, expected %dx%d", config.Width, config.Height, styleWidth, styleHeight)
		}

		img, err := png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data)))
		if err != nil {
			return errors.New("Error while decoding image")
		}

		placement.DecodedImage = img
	}

	return nil
}

func createProduct(request *requests.CreateProductRequest) ([]*model.Product, error) {
	/*
		return nil, nil
		pfVariant, err := getPrintfulVariant(request.VariantID)
		if err != nil {
			log.Println(err)
			return nil, errors.New("variant not found")
		}
	*/

	/*
		log.Println(pfVariant)
		pfProduct, _, err := getPrintfulProduct(pfVariant.CatalogProductID)
		if err != nil {
			log.Println(err)
			return nil, errors.New("product not found")
		}
	*/

	//	log.Println(pfProduct)

	placements := make([]map[string]interface{}, 0)
	for _, placement := range request.Placements {
		placements = append(placements, map[string]interface{}{
			"placement":   placement.Placement,
			"technique":   placement.Technique,
			"orientation": placement.Orientation,
		})
	}

	resp, err := fetchAPI("get-similar-variants", 1, map[string]interface{}{
		"variant_id": request.VariantID,
		"placements": placements,
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	similarVariantsResponse := printfulApiModel.SimilarVariantsResponse{}
	err = json.NewDecoder(resp.Body).Decode(&similarVariantsResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !similarVariantsResponse.Success {
		log.Println(similarVariantsResponse)
		return nil, errors.New("error while getting printful variant")
	}

	log.Println(similarVariantsResponse)

	variantCount := len(similarVariantsResponse.SimilarVariants)
	ids, err := createShopProducts(variantCount)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error while creating product: %w", err)
	}
	log.Println(ids)
	/*
		variants := make([]interface{}, 0, variantCount) //map[string]interface{}{}
		i := 0
		for i < variantCount {
			variant := map[string]interface{}{
				"variant_id":          similarVariantsResponse.SimilarVariants[i],
				"external_variant_id": ids[i],
				"retail_price":        9999,
			}

			variants = append(variants, variant)
			i += 1
		}

		log.Println(ids, err)

		/*
			resp, err = fetchAPI("create-sync-product", 1, map[string]interface{}{
				"product_id": pfVariant.CatalogProductID,
				"variants":   variants,
				"name":       request.Name,
				"image":      request.Image,
			})

			if err != nil {
				log.Println(err)
				return nil, errors.New("error while calling printful api")
			}
	*/

	/*
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
	*/
	/*
		response := CreateSyncProductResponse{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Println(err)
			return nil, errors.New("error while decoding printful response")
		}

		if !response.Success {
			log.Println(response)
			return nil, errors.New("error while creating printful product")
		}

		log.Println("createProduct", response)
	*/

	products, err := createShopProductFromPrintfulVariant(request.VariantID)
	if err != nil {
		return nil, fmt.Errorf("error while creating shop product %w", err)
	}
	log.Println(products)

	//return &variantResponse.Result.Variant, nil

	return nil, nil
}

type CreateSyncProductResponse struct {
	Success     bool                `json:"success"`
	SyncProduct schemas.SyncProduct `json:"result"`
}

type GetSyncProductResponse struct {
	Success         bool                             `json:"success"`
	SyncProductInfo printfulApiModel.SyncProductInfo `json:"result"`
}

func createShopProductFromPrintfulVariant(variantID int) (*model.Product, error) {
	log.Println("creating product for printful variant id:", variantID)

	pfVariant, err := getPrintfulVariant(variantID)
	if err != nil {
		return nil, fmt.Errorf("error while creating product from variant: %w", err)
	}

	log.Println(pfVariant)

	/*
		resp, err := fetchAPI("get-sync-product", 1, map[string]interface{}{
			"sync_product_id": variantID,
		})

		if err != nil {
			return nil, err
		}
	*/

	/*body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))*/
	/*
		response := GetSyncProductResponse{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Println(err)
			return nil, errors.New("error while decoding printful response")
		}

		if !response.Success {
			log.Println(response)
			return nil, errors.New("error while creating printful product")
		}
	*/

	//log.Println("createShopProduct", response)

	/*
		type SyncProductInfo struct {
			SyncProduct  SyncProduct   `json:"sync_product" bson:"sync_product"`
			SyncVariants []SyncVariant `json:"sync_variants" bson:"sync_variants"`
		}
	*/

	//	syncProduct := response.SyncProductInfo.SyncProduct
	//syncVariants := response.SyncProductInfo.SyncVariants

	/*
		variantIDs := []string{}
		for _, syncVariant := range syncVariants {
			variantIDs = append(variantIDs, syncVariant.ExternalID)
		}

		shopProducts := []*model.Product{}
		for _, syncVariant := range syncVariants {
			//v = append(v, key)
			shopProduct, err := createShopProduct2(syncProduct, syncVariant, variantIDs)
			shopProducts = append(shopProducts, shopProduct)

			if err != nil {
				log.Println(err)
				return nil, errors.New("error while creating shop product")
			}

			log.Println(shopProduct)
			/*
				const shoProduct = await this.#createShopProduct2(syncProduct, syncVariants);
				productsIds.push(shoProduct.id);
				products.push(shoProduct);
	*/
	//}
	/*
		if (productsIds.length > 1) {
			for (const productId of productsIds) {
				const updateOneResult = await this.#productsCollection.updateOne({ _id: productId }, { $set: { variantIds: productsIds }});
			}
		}
	*/

	// return the first product created
	//return products[0];

	return nil, nil
}

func createShopProduct2(syncProduct schemas.SyncProduct, syncVariant schemas.SyncVariant, variantIDs []string) (*model.Product, error) {
	product := model.NewProduct()
	product.Name = syncVariant.Name
	product.ProductName = syncProduct.Name
	//product.Currency = syncVariant.Currency
	product.ThumbnailURL = syncProduct.ThumbnailURL
	product.ExternalID1 = strconv.FormatInt(syncVariant.ID, 10)
	product.Status = "completed"
	product.VariantIDs = variantIDs

	retailPrice, err := decimal.NewFromString(syncVariant.RetailPrice)
	if err != nil {
		return nil, err
	}
	//product.RetailPrice = retailPrice
	log.Println("retailPrice", retailPrice)
	panic("add retail price")

	id, err := primitive.ObjectIDFromHex(syncVariant.ExternalID)
	if err != nil {
		return nil, err
	}
	product.ID = id

	pfVariant, err := getPrintfulVariant(syncVariant.VariantID)
	if err != nil {
		return nil, err
	}

	log.Println(pfVariant)

	if pfVariant.ColorCode != "" {
		product.AddOption("color", "color", pfVariant.ColorCode)
	}
	if pfVariant.ColorCode2 != "" {
		product.AddOption("color2", "color", pfVariant.ColorCode2)
	}
	if pfVariant.Size != "" {
		product.AddOption("size", "size", pfVariant.Size)
	}
	if pfVariant.Image != "" {
		product.AddFile("product", pfVariant.Image)
	}

	pfProduct, _, err := getPrintfulProduct(pfVariant.CatalogProductID)
	if err != nil {
		return nil, err
	}

	if pfProduct.Description != "" {
		product.Description = pfProduct.Description
	}

	for _, file := range syncVariant.Files {
		//v = append(v, key)
		product.AddFile(file.Type, file.URL)
	}

	err = mongo.UpdateProduct(&product)
	if err != nil {
		return nil, err
	}

	/*
		const shopProduct = new ShopProduct();
		const printfulProductReference = syncVariant.product;

		const printfulVariant = await this.#getPrintfulVariant(printfulProductReference.productId, printfulProductReference.variantId);
		if (!printfulVariant) {
			throw new Error(`Printful variant not found productId: ${printfulProductReference.productId} variantId: ${printfulProductReference.variantId}`);
		}

		const description = await this.#getPrintfulProductDescription(syncVariant?.product?.productId)
		if (description) {
			shopProduct.description = description;
		}

		const syncVariantFiles = syncVariant.files;
		if (syncVariantFiles) {
			for (const syncVariantFile of syncVariantFiles) {
				shopProduct.addFile(syncVariantFile.type, syncVariantFile.url);
			}
		}

		//console.log('createShopProduct2 printfulVariant', printfulVariant);
		const replaceOneResult = await this.#productsCollection.replaceOne({ _id: shopProduct.id }, shopProduct.toJSON());
		if (!replaceOneResult?.acknowledged) {
			winston.error('Error in #createShopProduct2 : replaceOne failed', { replaceOneResult: replaceOneResult, shopProduct: shopProduct.toJSON() });
			throw new Error('Error in #createShopProduct2 : replaceOne failed');
		}
		return shopProduct;
	*/
	/*
		if (productsIds.length > 1) {
			for (const productId of productsIds) {
				const updateOneResult = await this.#productsCollection.updateOne({ _id: productId }, { $set: { variantIds: productsIds }});
				//console.log(updateOneResult);
			}
		}
	*/

	return &product, nil
}

func createShopProducts(count int) ([]string, error) {
	ret := make([]string, 0, count)
	i := 0
	for i < count {
		product, err := mongo.CreateProduct()
		if err != nil {
			return nil, err
		}

		ret = append(ret, product.ID.Hex())

		i += 1
	}

	return ret, nil
}

func getPrintfulVariant(variantID int) (*printfulmodel.Variant, error) {
	/*u, err := url.JoinPath(printfulConfig.Endpoint, "/products/variant/", strconv.Itoa(int(variantID)))
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while getting printful url")
	}

	log.Println(u)
	resp, err := http.Get(u)*/
	resp, err := fetchAPI("get-variant", 1, map[string]interface{}{
		"variant_id": variantID,
	})

	//body, _ := ioutil.ReadAll(resp.Body)
	//log.Println(string(body))

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	variantResponse := responses.GetVariantResponse{}
	err = json.NewDecoder(resp.Body).Decode(&variantResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !variantResponse.Success {
		log.Println(variantResponse)
		return nil, errors.New("error while getting printful variant")
	}
	//log.Println("variantResponse", variantResponse)

	return &variantResponse.Result, nil
}

func getPrintfulProduct(productID int) (*printfulmodel.Product, []printfulmodel.Variant, error) {

	/*u, err := url.JoinPath(printfulConfig.Endpoint, "/product/", strconv.Itoa(int(productID)))
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while getting printful url")
	}

	log.Println(u)
	resp, err := http.Get(u)*/
	resp, err := fetchAPI("get-product", 1, map[string]interface{}{
		"product_id": productID,
	})

	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("error while calling printful api")
	}

	productResponse := responses.GetProductResponse{}
	err = json.NewDecoder(resp.Body).Decode(&productResponse)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("error while decoding printful response")
	}

	if !productResponse.Success {
		log.Println(productResponse)
		return nil, nil, errors.New("error while getting printful variant")
	}

	return &productResponse.Result.Product, productResponse.Result.Variants, nil
}

func getPrintfulMockupTemplates(productID int) ([]printfulmodel.MockupTemplates, error) {
	resp, err := fetchAPI("get-mockup-templates", 1, map[string]interface{}{
		"product_id": productID,
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	productResponse := responses.GetMockupTemplatesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&productResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !productResponse.Success {
		log.Println(productResponse)
		return nil, errors.New("error while getting printful variant")
	}

	return productResponse.Result.Templates, nil
}

func getPrintfulStyles(productID int) ([]printfulmodel.MockupStyles, error) {
	resp, err := fetchAPI("get-mockup-styles", 1, map[string]interface{}{
		"product_id": productID,
	})

	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	stylesResponse := responses.GetMockupStylesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&stylesResponse)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while decoding printful response")
	}

	if !stylesResponse.Success {
		log.Println(stylesResponse)
		return nil, errors.New("error while getting printful styles")
	}

	return stylesResponse.Result.Styles, nil
}
