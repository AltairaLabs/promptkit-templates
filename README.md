# PromptKit Community Templates

Index of community templates for PromptArena. The index is K8s-style:

```yaml
apiVersion: promptkit.altairalabs.ai/v1
kind: TemplateIndex
spec:
  entries:
    - name: basic-chatbot
      version: "1.0.0"
      description: Minimal chatbot template (mock provider)
      source: basic-chatbot/template.yaml
      tags: [chatbot, mock]
```

Templates are package files (YAML) with embedded files under `files:`. Example `basic-chatbot/template.yaml` includes Arena config, providers, prompts, scenarios, and README.

## Usage

List from this repo:
```bash
promptarena templates list --index https://raw.githubusercontent.com/AltairaLabs/promptkit-templates/main/index.yaml
```

Fetch/render:
```bash
promptarena templates fetch --index https://raw.githubusercontent.com/AltairaLabs/promptkit-templates/main/index.yaml --template basic-chatbot --version 1.0.0
promptarena templates render --template basic-chatbot --version 1.0.0 --values values.yaml --out ./out
```

Init with remote template:
```bash
promptarena init --template basic-chatbot --template-index https://raw.githubusercontent.com/AltairaLabs/promptkit-templates/main/index.yaml
```

## Contributing
- Add template packages under a directory matching the template name.
- Update `index.yaml` with the new entry (apiVersion `promptkit.altairalabs.ai/v1`, kind `TemplateIndex`).
- Keep versions semver-like (major.minor.patch).
