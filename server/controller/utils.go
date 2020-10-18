package controller

import (
	"html/template"
	"net/http"
	"path"
	"path/filepath"
)

func (c *Controller) renderTemplate(w http.ResponseWriter, templateName, contentType string, values interface{}) {
	bundlePath, err := c.api.GetBundlePath()
	if err != nil {
		c.api.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	templateDir := filepath.Join(bundlePath, "assets", "templates")
	tmplPath := path.Join(templateDir, templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	if err = tmpl.Execute(w, values); err != nil {
		http.Error(w, "failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
