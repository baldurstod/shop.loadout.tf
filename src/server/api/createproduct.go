package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"log"
	"math"
	"net/url"
	"slices"
	"strconv"
	"strings"

	printfulApiModel "github.com/baldurstod/go-printful-api-model"
	"github.com/baldurstod/go-printful-api-model/responses"
	"github.com/baldurstod/go-printful-api-model/schemas"
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/baldurstod/randstr"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/model/requests"
	"shop.loadout.tf/src/server/mongo"
)

var imagesConfig config.Images

func SetImagesConfig(config config.Images) {
	imagesConfig = config
}

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
		return nil, errors.New("error while getting similar variants")
	}

	//placementURLs := make(map[requests.CreateProductRequestPlacement]string)
	extraDataPlacements := make([]map[string]any, 0, len(request.Placements)) //extraData["printful"].(map[string]any)["placements"].([]map[string]any)
	for _, placement := range request.Placements {
		if placement.DecodedImage == nil {
			return nil, errors.New("decodedImage is nil")
		}

		filename := randstr.String(32)
		err = mongo.UploadImage(filename, placement.DecodedImage)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		imageURL, err := url.JoinPath(imagesConfig.BaseURL, "/", filename)
		if err != nil {
			return nil, errors.New("unable to create image url")
		}

		thumbnailURL, err := url.JoinPath(imagesConfig.BaseURL, "/", filename+"_thumb")
		if err != nil {
			return nil, errors.New("unable to create thumbnail url")
		}

		extraDataPlacement := map[string]any{
			"placement":   placement.Placement,
			"technique":   placement.Technique,
			"orientation": placement.Orientation,
			"image_url":   imageURL,
			"thumb_url":   thumbnailURL,
		}

		extraDataPlacements = append(extraDataPlacements, extraDataPlacement)
	}

	extraData := map[string]any{"printful": map[string]any{"placements": extraDataPlacements}}

	products := make([]*model.Product, 0, len(similarVariantsResponse.SimilarVariants))
	log.Println(similarVariantsResponse)
	for _, similarVariant := range similarVariantsResponse.SimilarVariants {
		product, err := createShopProductFromPrintfulVariant(similarVariant, extraData)
		if err != nil {
			return nil, fmt.Errorf("error while creating shop product %w", err)
		}
		products = append(products, product)
	}

	/*
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

	log.Println(products)

	//return &variantResponse.Result.Variant, nil

	return products, nil
}

type CreateSyncProductResponse struct {
	Success     bool                `json:"success"`
	SyncProduct schemas.SyncProduct `json:"result"`
}

type GetSyncProductResponse struct {
	Success         bool                             `json:"success"`
	SyncProductInfo printfulApiModel.SyncProductInfo `json:"result"`
}

func createShopProductFromPrintfulVariant(variantID int, extraData map[string]any) (*model.Product, error) {
	log.Println("creating product for printful variant id:", variantID)

	pfVariant, err := getPrintfulVariant(variantID)
	if err != nil {
		return nil, fmt.Errorf("error while creating product from variant: %w", err)
	}

	pfProduct, _, err := getPrintfulProduct(pfVariant.CatalogProductID)
	if err != nil {
		return nil, err
	}

	log.Println(pfVariant)

	product, err := mongo.CreateProduct()
	if err != nil {
		return nil, err
	}

	product.Name = pfVariant.Name
	product.ProductName = pfProduct.Name
	//product.Currency = syncVariant.Currency
	product.ThumbnailURL = pfProduct.Image
	product.ExternalID1 = strconv.FormatInt(int64(variantID), 10)
	product.Status = "created"
	product.ExtraData = extraData
	//product.VariantIDs = variantIDs
	/*
		retailPrice, err := decimal.NewFromString(syncVariant.RetailPrice)
		if err != nil {
			return nil, err
		}
	*/
	//product.RetailPrice = retailPrice
	//log.Println("retailPrice", retailPrice)
	//panic("add retail price")

	/*
		id, err := primitive.ObjectIDFromHex(syncVariant.ExternalID)
		if err != nil {
			return nil, err
		}
		product.ID = id
	*/
	/*
		pfVariant, err := getPrintfulVariant(syncVariant.VariantID)
		if err != nil {
			return nil, err
		}
	*/

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
	/*
		pfProduct, _, err := getPrintfulProduct(pfVariant.CatalogProductID)
		if err != nil {
			return nil, err
		}
	*/

	if pfProduct.Description != "" {
		product.Description = pfProduct.Description
	}

	/*
		for _, file := range pfVariant.Files {
			//v = append(v, key)
			product.AddFile(file.Type, file.URL)
		}
	*/

	err = mongo.UpdateProduct(product)
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

	return product, nil
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
		return nil, nil, errors.New("error while getting printful product")
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
		return nil, errors.New("error while getting mockup templates")
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
