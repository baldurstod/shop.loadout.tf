package requests

type CreateProductRequest struct {
	VariantID uint   `mapstructure:"variant_id"`
	Name      string `mapstructure:"name"`
	Type      string `mapstructure:"type"`
	Image     string `mapstructure:"image"`
}
