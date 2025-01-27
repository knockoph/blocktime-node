package handlers

import (
	"blocktime-node/pkg/core"
	"html/template"
	"io/fs"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request, btcClient *core.Client, contentFS fs.FS) {
	blockchainInfo, err := core.GetBlockchainInfo(btcClient)
	if err != nil {
		http.Error(w, "Error getting blockchain info", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(contentFS, "index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, blockchainInfo); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
