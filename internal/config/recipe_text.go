package config

import _ "embed"

//go:embed recipe_text_config.json
var recipeTextConfigRaw []byte

func GetRecipeTextConfigRaw() []byte {
	return recipeTextConfigRaw
}
