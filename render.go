package jsonschema2go

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func NewRenderer() *Renderer {
	return &Renderer{
		loader:  newCachingLoader(),
		planner: newPlanner(),
		printer: newPrinter(),
	}
}

type Renderer struct {
	loader  Loader
	planner *Planners
	printer *Printer
}

func (r *Renderer) Render(ctx context.Context, fileNames []string, prefixes [][2]string) error {
	seen := make(map[TypeInfo]bool)
	grouped := make(map[string][]Plan)
	for _, fileName := range fileNames {
		u, err := url.Parse(fileName)
		if err != nil {
			return err
		}

		schema, err := r.loader.Load(ctx, u)
		if err != nil {
			return fmt.Errorf("unable to resolve schema from %q: %w", fileName, err)
		}

		newPlans, err := r.planner.Plan(context.Background(), schema, r.loader)
		if err != nil {
			return fmt.Errorf("unable to create plans from schema %q: %w ", fileName, err)
		}

		for _, plan := range newPlans {
			if typ := plan.Type(); !seen[typ] {
				seen[typ] = true

				grouped[typ.GoPath] = append(grouped[typ.GoPath], plan)
			}
		}
	}

	sort.Slice(prefixes, func(i, j int) bool {
		return prefixes[i][0] < prefixes[j][0]
	})

	for k, group := range grouped {
		path := mapPath(k, prefixes)
		if path == "" {
			return fmt.Errorf("unable to map go path: %q", k[0])
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("unable to dir %q: %w", path, err)
		}
		if err := func() error {
			f, err := os.Create(filepath.Join(path, "values.gen.go"))
			if err != nil {
				return fmt.Errorf("unable to open: %w", err)
			}
			defer f.Close()

			if err := r.printer.Print(ctx, f, k, group); err != nil {
				return fmt.Errorf("unable to print: %w", err)
			}
			return nil
		}(); err != nil {
			return err
		}
	}
	return nil
}

func mapPath(path string, sortedPrefixes [][2]string) string {
	i := sort.Search(len(sortedPrefixes), func(i int) bool {
		return sortedPrefixes[i][0] > path
	})
	for i = i - 1; i >= 0; i-- {
		if strings.HasPrefix(path, sortedPrefixes[i][0]) {
			return sortedPrefixes[i][1] + path[len(sortedPrefixes[i][0]):]
		}
	}
	return ""
}