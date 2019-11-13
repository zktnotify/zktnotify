package viewmodel

type VerifyResult struct {
	Namespace       string `json:"namespace"`
	Field           string `json:"field"`
	StructNamespace string `json:"struct_namespace"`
	StructField     string `json:"struct_field"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actual_tag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
}
