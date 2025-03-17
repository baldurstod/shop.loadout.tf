package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
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
	printfulsdk "github.com/baldurstod/go-printful-sdk"
	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"github.com/baldurstod/randstr"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	printfulapi "shop.loadout.tf/src/server/api/printful"
	"shop.loadout.tf/src/server/config"
	"shop.loadout.tf/src/server/model"
	"shop.loadout.tf/src/server/model/requests"
	"shop.loadout.tf/src/server/mongo"
	"shop.loadout.tf/src/server/mongo/printfuldb"
)

var imagesConfig config.Images

func SetImagesConfig(config config.Images) {
	imagesConfig = config
}

func apiCreateProduct(c *gin.Context, params map[string]interface{}) error {
	if params == nil {
		return errors.New("no params provided")
	}

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

	styles, err := printfulapi.GetMockupStyles(request.ProductID) //getPrintfulStyles(request.ProductID)
	if err != nil {
		return errors.New("unable to get product styles")
	}

	for i, placement := range request.Placements {
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
	placements := make([]printfulapi.GetSimilarVariantsPlacement, 0)
	for _, placement := range request.Placements {
		placements = append(placements, printfulapi.GetSimilarVariantsPlacement{
			Placement:   placement.Placement,
			Technique:   placement.Technique,
			Orientation: placement.Orientation,
		})
	}

	similarVariants, err := printfulapi.GetSimilarVariants(request.VariantID, placements)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error while calling printful api")
	}

	/*
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
	*/

	extraDataPlacements := make([]map[string]any, 0, len(request.Placements))
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

		imageURL, err := url.JoinPath(imagesConfig.BaseURL, "/image/", filename)
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

	products := make([]*model.Product, 0, len(similarVariants))
	log.Println(similarVariants)

	mockupTemplates, err := printfulapi.GetMockupTemplates(request.ProductID) //getPrintfulMockupTemplates(request.ProductID)
	if err != nil {
		return nil, err
	}

	cache := make(map[image.Image]map[int]*model.MockupTask)
	imageCache := make(map[image.Image]string)
	tasks := make([]*model.MockupTask, 0, len(similarVariants))
	for _, similarVariant := range similarVariants {
		product, err := createShopProductFromPrintfulVariant(similarVariant, extraData, request.Technique, request.Placements, mockupTemplates, cache, imageCache, &tasks)
		if err != nil {
			return nil, fmt.Errorf("error while creating shop product %w", err)
		}
		products = append(products, product)
	}

	err = updateProductsVariants(products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = mongo.InsertMockupTasks(tasks)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = initMockupTasks(tasks)
	if err != nil {
		log.Println(err)
		return nil, err
	}

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

func computeProductPrice(productID int, variantID int, technique string, placements []*requests.CreateProductRequestPlacement, currency string) (decimal.Decimal, error) {
	productPrices, err := printfulapi.GetProductPrices(productID, currency, printfulConfig.Markup)
	if err != nil {
		log.Println(err)
		return decimal.NewFromInt(0), err
	}

	prices := map[string]decimal.Decimal{}
	for _, placement := range placements {
		for _, pricePlacement := range productPrices.Product.Placements {
			if pricePlacement.ID == placement.Placement && pricePlacement.TechniqueKey == technique {
				dec, err := decimal.NewFromString(pricePlacement.Price)

				if err != nil {
					log.Println(err)
					return decimal.NewFromInt(0), fmt.Errorf("can't convert string to decimal %s", pricePlacement.Price)
				}

				prices[pricePlacement.ID] = dec
			}
		}
	}

	maxPrice := decimal.NewFromInt(0)
	maxPricePlacement := ""
	for placement, price := range prices {
		if price.Compare(maxPrice) > 0 {
			maxPrice = price
			maxPricePlacement = placement
		}
	}

	if maxPricePlacement != "" {
		prices[maxPricePlacement] = decimal.NewFromInt(0)
	}

	//for _ :=range productPricesResponse.Result.Variants
	idx := slices.IndexFunc(productPrices.Variants, func(v printfulmodel.VariantsPriceData) bool { return v.ID == variantID })
	if idx == -1 {
		return decimal.NewFromInt(0), fmt.Errorf("variant %d not found", variantID)
	}

	variant := productPrices.Variants[idx]
	idx2 := slices.IndexFunc(variant.Techniques, func(v printfulmodel.TechniquePriceInfo) bool { return v.TechniqueKey == technique })
	if idx2 == -1 {
		return decimal.NewFromInt(0), fmt.Errorf("technique %s not found", technique)
	}

	techniquePriceInfo := variant.Techniques[idx2]

	variantPrice, err := decimal.NewFromString(techniquePriceInfo.Price)
	if err != nil {
		log.Println(err)
		return decimal.NewFromInt(0), fmt.Errorf("can't convert string to decimal %s", techniquePriceInfo.Price)
	}

	totalPrice := variantPrice
	for _, price := range prices {
		totalPrice = totalPrice.Add(price)
	}

	/*
		for (const placementPrice of productPrices.product.placements) {
			if (placementPrice.techniqueKey == technique && placements.has(placementPrice.id)) {
				placementsPrices.set(placementPrice.id, Number(placementPrice.price));
			}
		}
	*/

	return totalPrice, nil
}

func createShopProductFromPrintfulVariant(variantID int, extraData map[string]any, technique string, placements []*requests.CreateProductRequestPlacement, mockupTemplates []printfulmodel.MockupTemplates, cache map[image.Image]map[int]*model.MockupTask, imageCache map[image.Image]string, tasks *[]*model.MockupTask) (*model.Product, error) {
	log.Println("creating product for printful variant id:", variantID)

	pfVariant, _, err := printfuldb.FindVariant(variantID) //getPrintfulVariant(variantID)
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
	product.ThumbnailURL = pfVariant.Image
	product.ExternalID1 = strconv.FormatInt(int64(variantID), 10)
	product.Status = "created"
	product.ExtraData = extraData

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
		product.SetFile("product", pfVariant.Image, "")
	}

	if pfProduct.Description != "" {
		product.Description = pfProduct.Description
	}

	err = createMockupTasks(product.ID, pfVariant.ID, placements, mockupTemplates, cache, tasks)
	if err != nil {
		return nil, err
	}

	err = mongo.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	currency := "USD" //TODO: create currency variable
	price, err := computeProductPrice(pfVariant.CatalogProductID, variantID, technique, placements, currency)
	if err != nil {
		return nil, err
	}
	err = mongo.SetRetailPrice(product.ID, currency, price)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func createMockupTasks(productID string, variantID int, placements []*requests.CreateProductRequestPlacement, mockupTemplates []printfulmodel.MockupTemplates, cache map[image.Image]map[int]*model.MockupTask, tasks *[]*model.MockupTask) error {
	for i, placement := range placements {
		//log.Println(placement)
		idx := slices.IndexFunc(mockupTemplates, func(t printfulmodel.MockupTemplates) bool {
			if t.Orientation != placement.Orientation ||
				t.Technique != placement.Technique ||
				t.Placement != placement.Placement {
				return false
			}

			idx := slices.IndexFunc(t.CatalogVariantIDs, func(id int) bool { return id == variantID })
			return idx != -1
		})

		if idx == -1 {
			return fmt.Errorf("template not foundd for placement %d", i)
		}

		mockupTemplate := mockupTemplates[idx]

		cache1, found := cache[placement.DecodedImage]
		if !found {
			cache1 = make(map[int]*model.MockupTask)
			cache[placement.DecodedImage] = cache1
		}

		cache2, found := cache1[idx]
		//var img string //image.Image
		if found {
			//img = cache2
			//images[placement.Placement] = img
			/*
				taskID, err := mongo.InsertMockupTask(productID, "", nil, cache2)
				if err != nil {
					log.Printf("error while generating mockup template fro placement %s: %v", placement.Placement, err)
				} else {
					tasks[taskID] = true
				}
			*/
			cache2.AddProduct(productID)
		} else {
			//task, err := mongo.InsertMockupTask(productID, placement.Image, &mockupTemplate, nil)
			task := model.MockupTask{
				ProductIDs:  []string{productID},
				SourceImage: placement.Image,
				Template:    &mockupTemplate,
				//Status:      "created",
				//DateCreated: time.Now().Unix(),
				//DateUpdated: time.Now().Unix(),
			}
			/*if err != nil {
				log.Printf("error while generating mockup template fro placement %s: %v", placement.Placement, err)
			} else {
				//images[placement.Placement] = img
				cache1[idx] = task
				tasks[task] = true
			}*/
			cache1[idx] = &task
			*tasks = append(*tasks, &task)
		}
	}
	return nil
}

const (
	PositioningOverlay    string = "overlay"
	PositioningBackground string = "background"
)

func generateMockupTemplates(variantID int, placements []*requests.CreateProductRequestPlacement, mockupTemplates []printfulmodel.MockupTemplates, cache map[image.Image]map[int]image.Image) (map[string]image.Image, error) {
	images := make(map[string]image.Image)

	for i, placement := range placements {
		log.Println(placement)
		idx := slices.IndexFunc(mockupTemplates, func(t printfulmodel.MockupTemplates) bool {
			if t.Orientation != placement.Orientation ||
				t.Technique != placement.Technique ||
				t.Placement != placement.Placement {
				return false
			}

			idx := slices.IndexFunc(t.CatalogVariantIDs, func(id int) bool { return id == variantID })
			return idx != -1
		})

		if idx == -1 {
			return nil, fmt.Errorf("template not foundd for placement %d", i)
		}

		mockupTemplate := mockupTemplates[idx]

		cache1, found := cache[placement.DecodedImage]
		if !found {
			cache1 = make(map[int]image.Image)
			cache[placement.DecodedImage] = cache1
		}

		cache2, found := cache1[idx]
		var img image.Image
		if found {
			img = cache2
			images[placement.Placement] = img
		} else {
			img, err := printfulsdk.GenerateMockup(placement.DecodedImage, &mockupTemplate)
			if err != nil {
				log.Printf("error while generating mockup template fro placement %s: %v", placement.Placement, err)
			} else {
				images[placement.Placement] = img
				cache1[idx] = img
			}
		}
	}

	return images, nil
}

func updateProductsVariants(products []*model.Product) error {
	variantIDs := make([]string, 0, len(products))
	for _, product := range products {
		variantIDs = append(variantIDs, product.ID)
	}

	for _, product := range products {
		product.VariantIDs = variantIDs
		product.Status = "completed"

		err := mongo.UpdateProduct(product)
		if err != nil {
			return err
		}
	}
	return nil
}

func createShopProducts(count int) ([]string, error) {
	ret := make([]string, 0, count)
	i := 0
	for i < count {
		product, err := mongo.CreateProduct()
		if err != nil {
			return nil, err
		}

		ret = append(ret, product.ID)

		i += 1
	}

	return ret, nil
}

/*
func getPrintfulVariant(variantID int) (*printfulmodel.Variant, error) {
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
*/

func getPrintfulProduct(productID int) (*printfulmodel.Product, []printfulmodel.Variant, error) {
	product, err := printfulapi.GetProduct(productID)

	if err != nil {
		return nil, nil, err
	}

	variants, err := printfulapi.GetVariants(productID)

	if err != nil {
		return nil, nil, err
	}

	return product, variants, nil
	/*
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
	*/

}

/*
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
*/

/*
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
*/

func getPrintfulPrices(productID int) ([]printfulmodel.MockupStyles, error) {
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
