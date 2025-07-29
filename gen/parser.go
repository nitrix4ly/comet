package gen

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nitrix4ly/comet/core"
)

type Parser struct {
	schema *core.Schema
}

func NewParser() *Parser {
	return &Parser{
		schema: &core.Schema{},
	}
}

func (p *Parser) ParseFile(filename string) (*core.Schema, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentModel *core.ModelSchema
	var inModel bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "model ") {
			if currentModel != nil {
				p.schema.Models = append(p.schema.Models, *currentModel)
			}
			
			modelName := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "model "), "{"))
			currentModel = &core.ModelSchema{
				Name:      modelName,
				TableName: core.GetTableName(modelName),
				Fields:    []core.FieldSchema{},
				Relations: []core.Relation{},
			}
			inModel = true
			continue
		}

		if line == "}" && inModel {
			if currentModel != nil {
				p.schema.Models = append(p.schema.Models, *currentModel)
				currentModel = nil
			}
			inModel = false
			continue
		}

		if inModel && currentModel != nil {
			if err := p.parseField(line, currentModel); err != nil {
				return nil, fmt.Errorf("error parsing field '%s': %v", line, err)
			}
		}
	}

	if currentModel != nil {
		p.schema.Models = append(p.schema.Models, *currentModel)
	}

	return p.schema, scanner.Err()
}

func (p *Parser) parseField(line string, model *core.ModelSchema) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("invalid field definition")
	}

	fieldName := parts[0]
	fieldType := parts[1]
	
	field := core.FieldSchema{
		Name:     fieldName,
		Type:     strings.TrimSuffix(fieldType, "?"),
		Optional: strings.HasSuffix(fieldType, "?"),
	}

	if strings.HasSuffix(fieldType, "[]") {
		return p.parseRelation(line, model)
	}

	attributeStr := strings.Join(parts[2:], " ")
	if err := p.parseAttributes(attributeStr, &field); err != nil {
		return err
	}

	model.Fields = append(model.Fields, field)
	return nil
}

func (p *Parser) parseRelation(line string, model *core.ModelSchema) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("invalid relation definition")
	}

	fieldName := parts[0]
	fieldType := strings.TrimSuffix(parts[1], "[]")
	
	relation := core.Relation{
		Name:  fieldName,
		Type:  "hasMany",
		Model: fieldType,
	}

	attributeStr := strings.Join(parts[2:], " ")
	if err := p.parseRelationAttributes(attributeStr, &relation); err != nil {
		return err
	}

	model.Relations = append(model.Relations, relation)
	return nil
}

func (p *Parser) parseAttributes(attributeStr string, field *core.FieldSchema) error {
	re := regexp.MustCompile(`@(\w+)(?:\(([^)]*)\))?`)
	matches := re.FindAllStringSubmatch(attributeStr, -1)

	for _, match := range matches {
		attrName := match[1]
		attrValue := ""
		if len(match) > 2 {
			attrValue = match[2]
		}

		switch attrName {
		case "id":
			field.Primary = true
		case "auto":
			field.AutoGen = true
		case "unique":
			field.Unique = true
		case "default":
			field.Default = p.parseDefaultValue(attrValue)
		case "updatedAt":
			field.Type = "DateTime"
			field.Default = "now()"
		}
	}

	return nil
}

func (p *Parser) parseRelationAttributes(attributeStr string, relation *core.Relation) error {
	re := regexp.MustCompile(`@relation\("([^"]*)"(?:,\s*fields:\s*\[([^\]]*)\])?(?:,\s*references:\s*\[([^\]]*)\])?\)`)
	match := re.FindStringSubmatch(attributeStr)

	if len(match) > 1 {
		relation.Name = match[1]
	}
	if len(match) > 2 && match[2] != "" {
		relation.Fields = strings.Split(strings.ReplaceAll(match[2], " ", ""), ",")
	}
	if len(match) > 3 && match[3] != "" {
		relation.References = strings.Split(strings.ReplaceAll(match[3], " ", ""), ",")
	}

	if len(relation.Fields) > 0 && len(relation.References) > 0 {
		relation.Type = "belongsTo"
	}

	return nil
}

func (p *Parser) parseDefaultValue(value string) interface{} {
	value = strings.Trim(value, `"'`)
	
	switch value {
	case "now()":
		return "CURRENT_TIMESTAMP"
	case "true":
		return true
	case "false":
		return false
	default:
		if strings.Contains(value, ".") {
			return value
		}
		return value
	}
}
