package enterprise

type GetEnterprisesFilter struct {
  Orders map[string]string `json:"orders"`
  Filters map[string]interface{} `json:"filters"`
  Page uint64 `json:"page"`
}

type GetSuggestions struct {
  Filters map[string][]string `json:"filters"`
  Field string `json:"field"`
  Term string `json:"term"`
}