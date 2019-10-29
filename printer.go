package jsonschema2go

import (
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"sort"
	"text/template"
)

var (
	structTmpl = template.Must(fileTmplWithFuncs("templates/struct.tmpl"))
	valueTmpl  = template.Must(template.New("values.tmpl").ParseGlob("templates/*.tmpl"))
)

type Import struct {
	GoPath, Alias string
}

type Imports struct {
	currentGoPath string
	aliases       map[string]string
}

func newImports(currentGoPath string, importGoPaths []string) *Imports {
	baseName := make(map[string][]string)
	for _, i := range importGoPaths {
		if i != "" && i != currentGoPath {
			pkg := path.Base(i)
			baseName[pkg] = append(baseName[pkg], i)
		}
	}

	aliases := make(map[string]string)
	for k, v := range baseName {
		if len(v) == 1 {
			aliases[v[0]] = ""
			continue
		}
		sort.Strings(v)

		for i, path := range v {
			if i == 0 {
				aliases[path] = ""
				continue
			}
			aliases[path] = fmt.Sprintf("%s%d", k, i+1)
		}
	}

	return &Imports{currentGoPath, aliases}
}

func (i *Imports) CurPackage() string {
	return path.Base(i.currentGoPath)
}

func (i *Imports) List() (imports []Import) {
	for path, alias := range i.aliases {
		imports = append(imports, Import{path, alias})
	}
	return
}

func (i *Imports) QualName(info TypeInfo) string {
	if info.BuiltIn() || info.GoPath == i.currentGoPath {
		return info.Name
	}
	qual := path.Base(info.GoPath)
	if alias := i.aliases[info.GoPath]; alias != "" {
		qual = alias
	}
	return fmt.Sprintf("%s.%s", qual, info.Name)
}

func fileTmplWithFuncs(fName string) (*template.Template, error) {
	return template.New(filepath.Base(fName)).ParseFiles(fName)
}

func PrintFile(ctx context.Context, w io.Writer, goPath string, plans []Plan) error {
	var imports *Imports
	{
		var depPaths []string
		for _, p := range plans {
			for _, d := range p.Deps() {
				depPaths = append(depPaths, d.GoPath)
			}
		}
		imports = newImports(goPath, depPaths)
	}

	return valueTmpl.Execute(w, &Plans{imports, plans})
}

type Plans struct {
	Imports *Imports
	plans   []Plan
}

func (ps *Plans) Structs() (structs []structPlanContext) {
	for _, p := range ps.plans {
		if s, ok := p.(*StructPlan); ok {
			structs = append(structs, structPlanContext{ps.Imports, s})
		}
	}
	sort.Slice(structs, func(i, j int) bool {
		return structs[i].Type().Name < structs[j].Type().Name
	})
	return
}

func (ps *Plans) Arrays() (arrays []arrayPlanContext) {
	for _, p := range ps.plans {
		if a, ok := p.(*ArrayPlan); ok {
			arrays = append(arrays, arrayPlanContext{ps.Imports, a})
		}
	}
	sort.Slice(arrays, func(i, j int) bool {
		return arrays[i].Type().Name < arrays[j].Type().Name
	})
	return
}

func (ps *Plans) Enums() (enums []enumPlanContext) {
	for _, p := range ps.plans {
		if e, ok := p.(*EnumPlan); ok {
			enums = append(enums, enumPlanContext{ps.Imports, e})
		}
	}
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].Type().Name < enums[j].Type().Name
	})
	return
}

type structPlanContext struct {
	*Imports
	*StructPlan
}

type arrayPlanContext struct {
	*Imports
	*ArrayPlan
}

type enumPlanContext struct {
	*Imports
	*EnumPlan
}