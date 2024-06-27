package order

import (
	"net/http"
)

type Quantity struct {
	Type  string
	Value string
}

type Deadline struct {
	Type   string
	Value  string
	Action string
}

type CreationData struct {
	Pair      string
	Market    string
	Action    string
	Behavior  string
	Price     string
	Quantity  Quantity
	Deadlines []Deadline
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	// schema := jsonschema.NewReferenceLoaderFileSystem()
}

func Create(data CreationData) {

}
