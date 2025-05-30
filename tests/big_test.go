package tests

import (
	"bytes"
	"testing"

	"github.com/saasuke-labs/gsx"
	"github.com/stretchr/testify/assert"
)

const TagComponent = "<div class=\"tag\"><span>{{ .name }}</span></div>"
const LandingPageResult = `
<!DOCTYPE html>
<html>
<head>
	<title>Complex Page</title>
	<link rel="stylesheet" href="/static/style.css">
	<script src="/static/script.js"></script>
</head>
<body>
	<header>
		<h1>Welcome to the Complex Page</h1>
		<nav>
			<ul>
				<li><a href="/">Home</a></li>
				<li><a href="/about">About</a></li>
				<li><a href="/contact">Contact</a></li>
			</ul>
		</nav>
	</header>
	<main>
		<section>
			<h2>Section Title</h2>
			<p>This is a paragraph in the section.</p>
			<div class="tag"><span>Example Tag</span></div>
		</section>
		<section>
			<h2>Another Section</h2>
			<p>More content here.</p>
			<div class="tag"><span>Another Tag</span></div>
		</section>
		<section>
			<h2>Final Section</h2>
			<p>Last bit of content.</p>
			<div class="tag"><span>Final Tag</span></div>
		</section>
	</main>
	<footer>
		<p>&copy; 2023 Complex Page</p>
	</footer>
</body>
</html>
`

const LandingPage = `
<!DOCTYPE html>
<html>
<head>
	<title>Complex Page</title>
	<link rel="stylesheet" href="/static/style.css">
	<script src="/static/script.js"></script>
</head>
<body>
	<header>
		<h1>Welcome to the Complex Page</h1>
		<nav>
			<ul>
				<li><a href="/">Home</a></li>
				<li><a href="/about">About</a></li>
				<li><a href="/contact">Contact</a></li>
			</ul>
		</nav>
	</header>
	<main>
		<section>
			<h2>Section Title</h2>
			<p>This is a paragraph in the section.</p>
			<Tag name="Example Tag" />
		</section>
		<section>
			<h2>Another Section</h2>
			<p>More content here.</p>
			<Tag name="Another Tag" />
		</section>
		<section>
			<h2>Final Section</h2>
			<p>Last bit of content.</p>
			<Tag name="Final Tag" />
		</section>
	</main>
	<footer>
		<p>&copy; 2023 Complex Page</p>
	</footer>
</body>
</html>
`

func TestGenerate_ComplexHtmlPage(t *testing.T) {
	templ, _, err := gsx.ParseString("Tag", TagComponent)
	if err != nil {
		t.Fatalf("Failed to parse tag component: %v", err)
	}
	templ, _, err = gsx.ParseStringInto("LandingPage", LandingPage, templ)
	if err != nil {
		t.Fatalf("Failed to parse complex HTML page: %v", err)
	}

	buffer := &bytes.Buffer{}
	err = templ.ExecuteTemplate(buffer, "Tag", map[string]interface{}{
		"name": "tag name",
	})
	if err != nil {
		t.Fatalf("Failed to execute simple tag: %v", err)
	}
	assert.Equal(t, `<div class="tag"><span>tag name</span></div>`, buffer.String(), "The simple tag component did not render correctly")

	buffer.Reset()
	err = templ.ExecuteTemplate(buffer, "LandingPage", nil)
	if err != nil {
		t.Fatalf("Failed to execute complex HTML template: %v", err)
	}
	result := buffer.String()

	assert.Equal(t, LandingPageResult, result, "Generated HTML does not match expected output")

}
