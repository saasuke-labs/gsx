# GSX

**GSX** is a JSX-inspired templating syntax and parser for Go.

It lets you write HTML with component-like tags in `.gsx` files and compile them to Go templates.

---

## ✨ Example

Input GSX:

```html
<Tag name="Go" />
```

Compiles to:

```gotemplate
{{ template "Tag" (dict "name" "Go") }}
```

---

## ✅ Features

- JSX-style self-closing components (capitalized tags)
- Standard HTML syntax for lowercase elements
- Go template injection: `name={ .TagName }` → `{{ .TagName }}`
- Lint warnings for unquoted attributes
- Easy to integrate with static site generators, CMSs, and other Go tooling

---

## 🚫 Rules

- **GSX components must be self-closing**  
  Tags like `<PostCard />` are valid.  
  Tags like `<PostCard>...</PostCard>` are **not supported** in v1.

- **Only HTML elements may contain children**  
  You can write `<div><span>ok</span></div>`  
  But not `<Tag><span>bad</span></Tag>`

---

## 🔧 Usage

```go
import "github.com/saasuke-labs/gsx"

output, warnings, err := gsx.ParseString(`<Tag name="Go" />`, nil)
if err != nil {
    // handle error
}

fmt.Println(output) // → {{ template "Tag" (dict "name" "Go") }}
```

---

## 🧪 Lint Warnings

If you write something like:

```html
<Tag name=Go />
```

GSX will still parse it but emit a warning:

```
⚠️  unquoted value for attribute "name"
```

---

## 🛠 Integrations

You can use GSX inside:

- Static site generators (like [Gengo](https://github.com/saasuke-labs/gengo))
- Comment systems (like Kotomi)
- Any Go-based site renderer

---

## 📦 Installation

```bash
go get github.com/saasuke-labs/gsx
```

---

## 📄 License

MIT — © Saasuke Labs
