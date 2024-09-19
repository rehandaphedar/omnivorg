# Introduction

`omnivorg` is a program to sync [Omnivore](https://omnivore.app) highlights to an Org Mode file.

# Installation

```sh
go install git.sr.ht/~rehandaphedar/omnivorg@latest
```

# Configuration

## CLI Program

Create/Edit `$XDG_CONFIG_DIR/omnivorg/config.yaml`.
- `api_key` is your [Omnivore API Key](https://docs.omnivore.app/integrations/api.html#getting-an-api-token).
- `template_key` is the key to be passed to `org-capture`.
- `timestamp` is the last time the program was run. Do not edit this unless you know what you are doing.

Example configuration:
```yaml
api_key: dea6d3bb-c3bd-4aa2-9e1e-32fe07a0e520
template_key: o
timestamp: "2024-07-07T15:00:37+05:30"
```

## Emacs Setup

All the program does is fetches updates and runs:
```sh
xdg-open "org-protocol://capture?template=template&url=url&title=title&body=body
```

You need to have [Org Protocol](https://orgmode.org/worg/org-contrib/org-protocol.html) setup, along with an appropriate template key.

For example:
```elisp
(require 'org-protocol)
(setq org-capture-templates
	  `(("o" "Omnivore" entry
		 (file+headline
		  ,(expand-file-name
			(file-name-concat org-roam-directory "inbox.org")) "Omnivore")
		 (file "~/.config/emacs/roam-templates/omnivore.org")
		 :immediate-finish t)))
```

Where `~/.config/emacs/roam-templates/omnivore.org` is tangled from:
```org
,* [[%:link][%:description]]

%i
%?
```

Note that I also use [Org Roam](https://www.orgroam.com) (with custom file name templates). If you do not use it, you will not have access to the `org-roam-directory` variable.

# Usage

Make sure the [Emacs Server](https://www.gnu.org/software/emacs/manual/html_node/emacs/Emacs-Server.html) is running.
```sh
omnivorg
```

You can use utilities like `cron`, startup scripts in window managers, etc. to automate fetching updates.

# Limitations

I made this mainly for personal use. I only need the url, title, highlights, and annotations. I also do not need a full fledged template configuration. As such, this program only supports a basic template for now.

# Development

Clone the repository:
```sh
git clone https://git.sr.ht/~rehandaphedar/omnivorg
cd omnivorg
```

Fetch the schema (If updated):
```sh
./fetch-schema.sh
```

Generate queries:
```sh
go run github.com/Khan/genqlient
```

Run in development environment:
```sh
go run .
```

Build:
```sh
go build .
```
