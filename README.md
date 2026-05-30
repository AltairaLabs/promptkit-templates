# PromptKit Community Templates

Community templates for PromptArena. The index follows the Kubernetes-style resource shape:

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

Each entry points to a template package (`template.yaml`) plus any referenced files. We also include a `chatbot-source` example that uses the new `files[].source` field to load external content instead of embedding everything inline.

## Quick start

Install the CLI (requires Node):

```bash
npm install -g @altairalabs/promptarena
```

Generate a project from a remote template:

```bash
# See what’s available (uses the default community repo; no index URL needed)
promptarena templates list

# List with repo prefix (when multiple repos are configured)
promptarena templates list --index community

# Render the basic chatbot template (from the community repo)
# Fetch it into cache (one time)
promptarena templates fetch --template basic-chatbot --version 1.0.0

# Render from cache
promptarena templates render \
  --template basic-chatbot \
  --version 1.0.0 \
  --values values.example.yaml \
  --out ./out

# Or initialize a new project directly
promptarena init --template basic-chatbot

# Add another repo (optional) and list from it
promptarena templates repo add --name internal --url https://example.com/index.yaml
promptarena templates list --index internal
```

## Templates

- `basic-chatbot` — minimal mock-provider chatbot with inline file content.
- `chatbot-source` — similar chatbot, but uses `files[].source` to pull file bodies from separate files in the template package.
- `iot-maintenance-demo` — IoT predictive maintenance demo with red-team self-play, assertions, and SDK starter.
- `release-notes-agent` — release-notes agent packaged with the full engineering loop: eval-gated PRs against real models, a sticky failure comment, and a content-addressed GitHub Release on merge. The pack used by the [promptpack.org](https://promptpack.org) hero demo.

## Contributing
- Add template packages under a directory matching the template name.
- Update `index.yaml` with the new entry (`apiVersion: promptkit.altairalabs.ai/v1`, `kind: TemplateIndex`).
- Keep versions semver-like (major.minor.patch).
