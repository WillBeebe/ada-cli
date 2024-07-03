package projects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Entity represents a code entity such as a class, function, or interface.
type Entity struct {
	Label       string
	Name        string
	FileName    string
	Line        int
	Parameters  []string
	ReturnTypes []string
	Language    string
}

// Relationship represents a relationship between two entities.
type Relationship struct {
	From string
	Type string
	To   string
}

// CodebaseParser is used to parse the codebase.
type CodebaseParser struct {
	Entities      []Entity
	Relationships []Relationship
	FileSet       *token.FileSet
	FileName      string
	IgnoreList    []string
	PkgAliases    map[string]string // For storing package aliases
	Language      string
}

func (p *CodebaseParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		if n.Tok == token.TYPE {
			for _, spec := range n.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				position := p.FileSet.Position(typeSpec.Pos())
				entity := Entity{
					Label:    "Class",
					Name:     typeSpec.Name.Name,
					FileName: p.FileName,
					Line:     position.Line,
					Language: p.Language,
				}
				p.Entities = append(p.Entities, entity)
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					for _, field := range structType.Fields.List {
						if ident, ok := field.Type.(*ast.Ident); ok {
							p.Relationships = append(p.Relationships, Relationship{From: typeSpec.Name.Name, Type: "HAS_FIELD", To: ident.Name})
						}
					}
				} else if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
					for _, method := range interfaceType.Methods.List {
						if funcType, ok := method.Type.(*ast.FuncType); ok {
							p.Relationships = append(p.Relationships, Relationship{From: typeSpec.Name.Name, Type: "DECLARES_METHOD", To: method.Names[0].Name})
							// Capture parameters and return types for interface methods
							params := captureParams(funcType.Params)
							returns := captureParams(funcType.Results)
							p.Entities = append(p.Entities, Entity{
								Label:       "Function",
								Name:        method.Names[0].Name,
								FileName:    p.FileName,
								Line:        position.Line,
								Parameters:  params,
								ReturnTypes: returns,
								Language:    p.Language,
							})
						}
					}
				}
			}
		}
	case *ast.FuncDecl:
		position := p.FileSet.Position(n.Pos())
		params := captureParams(n.Type.Params)
		returns := captureParams(n.Type.Results)
		entity := Entity{
			Label:       "Function",
			Name:        n.Name.Name,
			FileName:    p.FileName,
			Line:        position.Line,
			Parameters:  params,
			ReturnTypes: returns,
			Language:    p.Language,
		}
		p.Entities = append(p.Entities, entity)
		if n.Recv != nil {
			for _, recv := range n.Recv.List {
				if starExpr, ok := recv.Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						p.Relationships = append(p.Relationships, Relationship{From: ident.Name, Type: "CONTAINS", To: n.Name.Name})
					}
				}
			}
		}
		for _, stmt := range n.Body.List {
			ast.Inspect(stmt, func(node ast.Node) bool {
				switch cn := node.(type) {
				case *ast.ExprStmt:
					if callExpr, ok := cn.X.(*ast.CallExpr); ok {
						p.addCallRelationship(n.Name.Name, callExpr)
					}
				case *ast.CallExpr:
					p.addCallRelationship(n.Name.Name, cn)
				}
				return true
			})
		}
	case *ast.ImportSpec:
		// Capture imports and dependencies
		pkgName := strings.Trim(n.Path.Value, `"`)
		alias := ""
		if n.Name != nil {
			alias = n.Name.Name
		} else {
			parts := strings.Split(pkgName, "/")
			alias = parts[len(parts)-1]
		}
		p.PkgAliases[alias] = pkgName
		entity := Entity{
			Label:    "Package",
			Name:     pkgName,
			FileName: p.FileName,
			Line:     p.FileSet.Position(n.Pos()).Line,
			Language: p.Language,
		}
		p.Entities = append(p.Entities, entity)
	}
	return p
}

// captureParams captures the parameter or return type names
func captureParams(fieldList *ast.FieldList) []string {
	if fieldList == nil {
		return nil
	}
	var params []string
	for _, field := range fieldList.List {
		for _, name := range field.Names {
			params = append(params, name.Name)
		}
	}
	return params
}

// addCallRelationship adds a CALLS relationship between functions, including cross-package calls
func (p *CodebaseParser) addCallRelationship(callingFunc string, callExpr *ast.CallExpr) {
	if ident, ok := callExpr.Fun.(*ast.Ident); ok {
		p.Relationships = append(p.Relationships, Relationship{From: callingFunc, Type: "CALLS", To: ident.Name})
	} else if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		if pkgIdent, ok := selExpr.X.(*ast.Ident); ok {
			pkgAlias := pkgIdent.Name
			funcName := selExpr.Sel.Name
			pkgName, exists := p.PkgAliases[pkgAlias]
			if exists {
				p.Relationships = append(p.Relationships, Relationship{From: callingFunc, Type: "CALLS", To: pkgName + "." + funcName})
				p.Relationships = append(p.Relationships, Relationship{From: p.FileName, Type: "DEPENDS_ON", To: pkgName})
			} else {
				p.Relationships = append(p.Relationships, Relationship{From: callingFunc, Type: "CALLS", To: pkgAlias + "." + funcName})
			}
		}
	}
}

func parseCodebase(path string, ignoreList []string) ([]Entity, []Relationship, error) {
	fset := token.NewFileSet()

	entities := []Entity{}
	relationships := []Relationship{}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("path: %s\n", path)

		// Check if the path should be ignored
		for _, ignore := range ignoreList {
			if strings.Contains(path, ignore) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			cbparser := &CodebaseParser{FileSet: fset, IgnoreList: ignoreList, PkgAliases: make(map[string]string), FileName: path, Language: "Go"}

			file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return err
			}
			ast.Walk(cbparser, file)

			// Append Go entities and relationships
			entities = append(entities, cbparser.Entities...)
			relationships = append(relationships, cbparser.Relationships...)

		} else if !info.IsDir() && strings.HasSuffix(info.Name(), ".py") {
			// Parse Python file
			pythonEntities, pythonRelationships, err := parsePythonFile(path)
			if err != nil {
				return err
			}
			for _, entity := range pythonEntities {
				entity.Language = "Python" // Set language to Python
				entities = append(entities, entity)
			}
			for _, relationship := range pythonRelationships {
				relationships = append(relationships, relationship)
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return entities, relationships, nil
}

// parsePythonFile parses a Python file and returns entities and relationships.
func parsePythonFile(path string) ([]Entity, []Relationship, error) {
	cmd := exec.Command("python3", "parse_python.py", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, nil, err
	}

	var result struct {
		Entities      []Entity
		Relationships []Relationship
	}
	err = json.Unmarshal(out.Bytes(), &result)
	if err != nil {
		return nil, nil, err
	}

	return result.Entities, result.Relationships, nil
}

func createGraph(driver neo4j.Driver, entities []Entity, relationships []Relationship) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		for _, entity := range entities {
			_, err := tx.Run("CREATE (n:"+entity.Label+" {name: $name, file: $file, line: $line, params: $params, returns: $returns, language: $language})", map[string]interface{}{
				"name":     entity.Name,
				"file":     entity.FileName,
				"line":     entity.Line,
				"params":   entity.Parameters,
				"returns":  entity.ReturnTypes,
				"language": entity.Language,
			})
			if err != nil {
				return nil, err
			}
		}
		for _, rel := range relationships {
			_, err := tx.Run(
				"MATCH (a {name: $from}), (b {name: $to}) CREATE (a)-[:"+rel.Type+"]->(b)",
				map[string]interface{}{"from": rel.From, "to": rel.To})
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})

	return err
}

func Graph(codebasePath string) {
	ignoreList := []string{"vendor", "testdata", "ignore.go"}
	entities, relationships, err := parseCodebase(codebasePath, ignoreList)
	if err != nil {
		log.Fatalf("Error parsing codebase: %v", err)
	}

	neo4jUri := "bolt://localhost:7687"
	neo4jUsername := "neo4j"
	neo4jPassword := "test"

	driver, err := neo4j.NewDriver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))
	if err != nil {
		log.Fatalf("Error creating Neo4j driver: %v", err)
	}
	defer driver.Close()

	err = createGraph(driver, entities, relationships)
	if err != nil {
		log.Fatalf("Error creating graph in Neo4j: %v", err)
	}

	fmt.Println("Codebase successfully imported into Neo4j!")
}
