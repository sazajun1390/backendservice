package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ormFlag = flag.String("orm", "bun", "ORM type: bun or gorm")

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("usage: patch-models --orm=bun|gorm <file_or_dir>...")
		os.Exit(1)
	}

	for _, path := range flag.Args() {
		processPath(path)
	}
}

func processPath(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot access %s: %v\n", path, err)
		return
	}

	if info.IsDir() {
		filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if !d.IsDir() && strings.HasSuffix(p, ".go") {
				patchFile(p)
			}
			return nil
		})
	} else {
		patchFile(path)
	}
}

func patchFile(filename string) {
	src, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", filename, err)
		return
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", filename, err)
		return
	}

	// ORM用 import 追加
	switch *ormFlag {
	case "bun":
		ensureImport(node, "github.com/uptrace/bun")
	case "gorm":
		ensureImport(node, "gorm.io/gorm")
	default:
		fmt.Fprintf(os.Stderr, "unsupported orm: %s\n", *ormFlag)
		return
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 埋め込み追加
			if *ormFlag == "bun" && !hasEmbed(st, "bun.BaseModel") {
				st.Fields.List = prependField(st.Fields.List, "bun.BaseModel")
			} else if *ormFlag == "gorm" && !hasEmbed(st, "gorm.Model") {
				st.Fields.List = prependField(st.Fields.List, "gorm.Model")
			}

			// 各フィールド処理
			for _, f := range st.Fields.List {
				if len(f.Names) == 0 {
					continue
				}
				name := f.Names[0].Name
				tag := getTag(f)

				if *ormFlag == "bun" {
					// sql.NullTime -> bun.NullTime
					if sel, ok := f.Type.(*ast.SelectorExpr); ok {
						if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "sql" && sel.Sel.Name == "NullTime" {
							f.Type = ast.NewIdent("bun.NullTime")
						}
					}

					// bunタグ追加
					if !strings.Contains(tag, "bun:") {
						tag = addBunTag(tag, name)
						setTag(f, tag)
					}
				}
			}
		}
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		fmt.Fprintf(os.Stderr, "print %s: %v\n", filename, err)
		return
	}
	os.WriteFile(filename, buf.Bytes(), 0644)
}

// --- ヘルパー ---

func prependField(list []*ast.Field, typeName string) []*ast.Field {
	return append([]*ast.Field{{Type: ast.NewIdent(typeName)}}, list...)
}

func hasEmbed(st *ast.StructType, embed string) bool {
	for _, f := range st.Fields.List {
		if ident, ok := f.Type.(*ast.Ident); ok && ident.Name == embed {
			return true
		}
		if sel, ok := f.Type.(*ast.SelectorExpr); ok {
			if pkg, ok := sel.X.(*ast.Ident); ok && fmt.Sprintf("%s.%s", pkg.Name, sel.Sel.Name) == embed {
				return true
			}
		}
	}
	return false
}

func ensureImport(f *ast.File, pkg string) {
	for _, imp := range f.Imports {
		if imp.Path.Value == fmt.Sprintf("%q", pkg) {
			return
		}
	}

	// 既存 import ブロックがあれば追加
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", pkg)},
			})
			return
		}
	}

	// import ブロックがない場合は新規作成して先頭に追加
	newDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", pkg)},
			},
		},
	}
	f.Decls = append([]ast.Decl{newDecl}, f.Decls...)
}

func getTag(f *ast.Field) string {
	if f.Tag != nil {
		return strings.Trim(f.Tag.Value, "`")
	}
	return ""
}

func setTag(f *ast.Field, tag string) {
	f.Tag = &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("`%s`", tag)}
}

func addBunTag(existing, name string) string {
	col := toSnakeCase(name)
	if existing != "" {
		return existing + " " + fmt.Sprintf(`bun:"%s"`, col)
	}
	return fmt.Sprintf(`bun:"%s"`, col)
}

func toSnakeCase(s string) string {
	var out []rune
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			out = append(out, '_')
		}
		out = append(out, rune(strings.ToLower(string(r))[0]))
	}
	return string(out)
}
