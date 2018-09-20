package templates

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func commentifyString(in string) string {
	var newlines []string
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			newlines = append(newlines, line)
		} else {
			if line == "" {
				newlines = append(newlines, "//")
			} else {
				newlines = append(newlines, "// "+line)
			}
		}
	}
	return strings.Join(newlines, "\n")
}

func ExecuteTemplate(name string, data interface{}) (*bytes.Buffer, error) {
	asset, err := Asset(name)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("").Funcs(template.FuncMap{"comment": commentifyString, "basename": filepath.Base}).Parse(string(asset))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func RestoreTemplate(path, name string, data interface{}) error {
	buf, err := ExecuteTemplate(name, data)
	if err != nil {
		return err
	}

	info, err := AssetInfo(name)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf.Bytes(), info.Mode())
	if err != nil {
		return err
	}

	return nil
}

func RestoreTemplates(dir, name, prefix string, data interface{}) error {
	filename := strings.TrimPrefix(name, prefix)
	children, err := AssetDir(name)

	if err != nil {
		if filepath.Ext(filename) == ".tpl" {
			filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		}

		return RestoreTemplate(filepath.Join(dir, filename), name, data)
	}

	for _, child := range children {
		err = RestoreTemplates(dir, filepath.Join(name, child), prefix, data)
		if err != nil {
			return err
		}
	}

	return nil
}
