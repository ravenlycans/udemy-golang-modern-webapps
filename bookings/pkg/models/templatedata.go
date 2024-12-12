package models

// TemplateData is the type for passing data from handlers into templates.
type TemplateData struct {
	StringMap  map[string]string
	IntMap     map[string]int
	FloatMap   map[string]float32
	Data       map[string]interface{}
	CSRFToken  string
	FlashMsg   string
	WarningMsg string
	ErrorMsg   string
}
