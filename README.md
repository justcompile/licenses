# licenses

A tool for inspecting application dependencies requirements/packages/modules in order to extract their license

## Usage
```
license examine [SOURCE] [flags]

[SOURCE] can be a file path or '-' to read from stdin 
```

The following languages are supported:

### Python 
`--lang=py or --lang=python`

Supporting Pip formatting only currently
```bash
# To examine the packages from a requirements file:
license examine requirements.txt --lang=py 

# to examine the packages directly from pip
pip freeze | license examine - --lang=py
```

### Golang
`--lang=go or --lang=golang`

Supporting go modules 

```bash
license examine go.mod --lang=go
```