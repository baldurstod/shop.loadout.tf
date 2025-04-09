package printfulapi

import (
	"errors"
	"slices"

	printfulmodel "github.com/baldurstod/go-printful-sdk/model"
	"shop.loadout.tf/src/server/mongo/printfuldb"
)

func GetVariants(productID int) ([]printfulmodel.Variant, error) {
	variants, _, err := printfuldb.FindVariants(productID)
	if err == nil {
		return variants, nil
	}

	return nil, errors.New("unable to find variants")
}

func GetVariant(variantID int) (*printfulmodel.Variant, error) {
	variant, _, err := printfuldb.FindVariant(variantID)
	if err == nil {
		return variant, nil
	}

	return nil, errors.New("unable to find variant")
}

type GetSimilarVariantsPlacement struct {
	Placement   string `json:"placement"`
	Technique   string `json:"technique"`
	Orientation string `json:"orientation"`
}

func GetSimilarVariants(variantID int, placements []GetSimilarVariantsPlacement) ([]int, error) {
	if placements == nil {
		return nil, errors.New("placement is empty")
	}

	variant, err := GetVariant(variantID)
	if err != nil {
		return nil, err
	}

	product, err := GetProduct(variant.CatalogProductID)
	if err != nil {
		return nil, err
	}

	templates, err := GetMockupTemplates(variant.CatalogProductID)
	if err != nil {
		return nil, err
	}

	variantsIDs := make(map[int]int, 0)

	for _, v := range product.CatalogVariantIDs {
		if (variantID == v) || matchTemplate(templates, variantID, v, placements) {
			variantsIDs[v] = v
		}
	}

	keys := make([]int, len(variantsIDs))
	i := 0
	for k := range variantsIDs {
		keys[i] = k
		i++
	}

	return keys, nil
}

func matchTemplate(templates []printfulmodel.MockupTemplates, v1 int, v2 int, placements []GetSimilarVariantsPlacement) bool {
	for _, placement := range placements {
		template1 := findTemplate(templates, v1, &placement)
		if template1 == nil {
			return false
		}

		template2 := findTemplate(templates, v2, &placement)
		if template2 == nil {
			return false
		}

		if template1.PrintAreaWidth == 0 || template1.PrintAreaHeight == 0 {
			return false
		}

		if template1.TemplateWidth != template2.TemplateWidth ||
			template1.TemplateHeight != template2.TemplateHeight ||
			template1.PrintAreaWidth != template2.PrintAreaWidth ||
			template1.PrintAreaHeight != template2.PrintAreaHeight ||
			template1.PrintAreaTop != template2.PrintAreaTop ||
			template1.PrintAreaLeft != template2.PrintAreaLeft {
			return false
		}
	}
	return true
}

func findTemplate(templates []printfulmodel.MockupTemplates, variantID int, placement *GetSimilarVariantsPlacement) *printfulmodel.MockupTemplates {
	idx := slices.IndexFunc(templates, func(t printfulmodel.MockupTemplates) bool {
		if t.Orientation != placement.Orientation ||
			t.Technique != placement.Technique ||
			t.Placement != placement.Placement {
			return false
		}

		idx := slices.IndexFunc(t.CatalogVariantIDs, func(id int) bool { return id == variantID })

		return idx == -1
	})

	if idx == -1 {
		return nil
	}

	return &templates[idx]
}
