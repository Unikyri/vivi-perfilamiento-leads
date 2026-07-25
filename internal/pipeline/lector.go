package pipeline

// Diccionarios de decodificación verificados contra el archivo v1.
// Las claves son los valores crudos del Excel v2 (nombres griegos);
// los valores son las categorías de negocio normalizadas.

var categoriaMap = map[string]string{
	"OMEGA": "A", "ETA": "B", "TAU": "C",
	"CHI": "SIN_DATO", "PI": "SIN_DATO", "XI": "SIN_DATO", "": "SIN_DATO",
}

var segmentoMap = map[string]string{
	"KAPPA": "BASICO", "SIGMA": "MEDIO", "NU": "JOVEN", "IOTA": "ALTO",
	"PI": "SIN_DATO", "CHI": "SIN_DATO", "XI": "SIN_DATO", "": "SIN_DATO",
}

// edadMap unifica las dos grafías que aparecen en el Excel v2
// ("20 - 35 años" vs "20 a 35 años") según doc 13 §0 corrección C4.
var edadMap = map[string]string{
	"20 - 35 años":     "20-35",
	"20 a 35 años":     "20-35",
	"36 - 45 años":     "36-45",
	"36 a 45 años":     "36-45",
	"46 - 55 años":     "46-55",
	"46 a 55 años":     "46-55",
	"Mayor de 55 años": "55+",
	"Menor de 19 años": "SIN_DATO",
	"":                 "SIN_DATO",
}
