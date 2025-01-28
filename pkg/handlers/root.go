package handlers

import (
	"blocktime-node/pkg/core"
	"blocktime-node/pkg/utils"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request, contentFS fs.FS, info *core.Info) {
	message, err := utils.Message(info, true)
	if err != nil {
		log.Println(fmt.Errorf("error in handle root: %w", err))
	}

	tmpl, err := template.ParseFS(contentFS, "index.html")
	if err != nil {
		log.Println(fmt.Errorf("error in handle root: %w", err))
		http.Error(w, "error parsing template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, message); err != nil {
		log.Println(fmt.Errorf("error in handle root: %w", err))
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
