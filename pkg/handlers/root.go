package handlers

import (
	"blocktime-node/pkg/core"
	"encoding/json"
	"html/template"
	"io/fs"
	"net/http"
)

// HandleRoot handles requests to the root path
func HandleRoot(w http.ResponseWriter, r *http.Request, btcClient *core.Client, contentFS fs.FS) {
	// Call getblockchaininfo to get the current block count
	result, err := btcClient.Call("getblockchaininfo", nil)
	if err != nil {
		http.Error(w, "Error calling Bitcoin Core: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the result to get the block count
	var blockchainInfo struct {
		Blocks int `json:"blocks"`
	}

	if err := json.Unmarshal(result, &blockchainInfo); err != nil {
		http.Error(w, "Error parsing blockchain info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Read the index.html template
	tmpl, err := template.ParseFS(contentFS, "index.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template with the blockchain info
	if err := tmpl.Execute(w, blockchainInfo); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
