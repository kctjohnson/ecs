package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func (tm *TemplateManager) LoadEnemyTemplates(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var templates []EnemyTemplate
	err = json.Unmarshal(data, &templates)
	if err != nil {
		return fmt.Errorf("error unmarshalling enemy templates: %w", err)
	}

	for _, tmpl := range templates {
		tm.EnemyTemplates[tmpl.ID] = tmpl
	}

	return nil
}

func (tm *TemplateManager) LoadItemTemplates(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var templates []ItemTemplate
	err = json.Unmarshal(data, &templates)
	if err != nil {
		return fmt.Errorf("error unmarshalling item templates: %w", err)
	}

	for _, tmpl := range templates {
		tm.ItemTemplates[tmpl.ID] = tmpl
	}

	return nil
}

func (tm *TemplateManager) LoadAllTemplates(templatesDir string) error {
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return fmt.Errorf("templates directory does not exist: %w", err)
	}

	templates := map[string]func(string) error{
		"enemy_templates.json": tm.LoadEnemyTemplates,
		"item_templates.json":  tm.LoadItemTemplates,
	}

	for filename, loadFunc := range templates {
		filePath := fmt.Sprintf("%s/%s", templatesDir, filename)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			log.Printf("Warning: template file does not exist: %s", filePath)
			continue
		}

		if err := loadFunc(filePath); err != nil {
			return fmt.Errorf("error loading %s: %w", filePath, err)
		}
	}

	return nil
}
