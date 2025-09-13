package view

import (
    "html/template"
    "io"
    "bytes"
    "os"
    "path/filepath"
)

// Manager loads and executes file-based templates.
type Manager struct {
    dir   string
    dev   bool
    tmpl  *template.Template
}

// NewManager creates a new template manager rooted at dir.
func NewManager(dir string, dev bool) *Manager {
    return &Manager{dir: dir, dev: dev}
}

// funcs returns the template function map.
func (m *Manager) funcs() template.FuncMap {
    return template.FuncMap{
        "safeHTML": func(s string) template.HTML { return template.HTML(s) },
    }
}

// parse (re)parses all templates under the manager directory.
func (m *Manager) parse() error {
    // Resolve dir for different working directories (tests vs. prod)
    candidates := []string{
        m.dir,
        filepath.Join("..", "..", m.dir),
        filepath.Join(".", m.dir),
    }
    var root string
    for _, c := range candidates {
        if st, err := os.Stat(c); err == nil && st.IsDir() {
            root = c
            break
        }
    }
    if root == "" {
        // Fall back to configured dir; will error below when walking
        root = m.dir
    }

    var files []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if filepath.Ext(path) == ".html" {
            files = append(files, path)
        }
        return nil
    })
    if err != nil {
        return err
    }

    base := template.New("base").Funcs(m.funcs())
    tmpl, err := base.ParseFiles(files...)
    if err != nil {
        return err
    }
    m.tmpl = tmpl
    return nil
}

// Execute renders a named template. In dev, templates are reparsed each time.
func (m *Manager) Execute(w io.Writer, name string, data interface{}) error {
    if m.dev || m.tmpl == nil {
        if err := m.parse(); err != nil {
            return err
        }
    }
    return m.tmpl.ExecuteTemplate(w, name, data)
}

// RenderString renders a named template to a string (useful for partials/pages).
func (m *Manager) RenderString(name string, data interface{}) (string, error) {
    if m.dev || m.tmpl == nil {
        if err := m.parse(); err != nil {
            return "", err
        }
    }
    var buf bytes.Buffer
    if err := m.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
        return "", err
    }
    return buf.String(), nil
}
