package main

import "fmt"

type python struct {
	language *language
	venvName string
}

const (
	// DisplayVirtualEnv shows or hides the virtual env
	DisplayVirtualEnv Property = "display_virtual_env"
)

func (p *python) string() string {
	if p.venvName == "" || !p.language.props.getBool(DisplayVirtualEnv, true) {
		return p.language.string()
	}
	version := p.language.string()
	if version == "" {
		return p.venvName
	}
	return fmt.Sprintf("%s %s", p.venvName, version)
}

func (p *python) init(props *properties, env environmentInfo) {
	p.language = &language{
		env:          env,
		props:        props,
		commands:     []string{"python", "python3"},
		versionParam: "--version",
		extensions:   []string{"*.py", "*.ipynb", "pyproject.toml", "venv.bak", "venv", ".venv"},
		versionRegex: `Python (?P<version>[0-9]+.[0-9]+.[0-9]+)`,
		loadContext:  p.loadContext,
		inContext:    p.inContext,
	}
}

func (p *python) enabled() bool {
	return p.language.enabled()
}

func (p *python) loadContext() {
	venvVars := []string{
		"VIRTUAL_ENV",
		"CONDA_ENV_PATH",
		"PYENV_VERSION",
	}
	var venv string
	for _, venvVar := range venvVars {
		venv = p.language.env.getenv(venvVar)
		if venv != "" {
			p.venvName = base(venv, p.language.env)
			break
		}
	}
}

func (p *python) inContext() bool {
	return p.venvName != ""
}
