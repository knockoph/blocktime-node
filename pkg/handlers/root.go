package handlers

import (
	"blocktime-node/pkg/core"
	"html/template"
	"io/fs"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request, contentFS fs.FS, info *core.Info) {
	blocks, err := info.GetBlocks(true)
	if err != nil {
		http.Error(w, "Error getting blockchain info", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(contentFS, "index.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	blockchainInfo := core.BlockchainInfo{
		Blocks: blocks,
	}

	if err := tmpl.Execute(w, blockchainInfo); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
