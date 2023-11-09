# go-blog
A personal go script that turns Markdown posts and HTML templates into a static blog.

## Usage
`./go-blog`

### Building from Source
1. Clone the repo with `git clone https://www.github.com/nateranda/go-blog`
2. Build with `go build go-blog.go`

### Extra Note
In order to protect against accidential data deletion when re-building, go-blog doesn't remove the contents of the `build` directory and instead just overwrites the HTML files. This means that any posts that are published and subsequently un-published will still be in the `build` directory, even though the post won't appear on the site. Keep this in mind when pushing to a public repository.

## Layout
go-blog looks for Markdown posts in the `posts` directory and templates in the `templates` directory at the root of the program.

### HTML Templating
go-blog uses the `text/template` package for templating. It looks for three templates: `menu.html`, `post.html`, and `tag.html`.

`menu.html` is passed an array of post structs, `post.html` is passed a post struct, and `tag.html` is passed a tag struct. More info about these structs is available in the documentation, and an example template is included in the source code.

### Markdown Metadata
Each Markdown post should contain YAML frontmatter:
```yaml
---
name: "Test 1" # Post name
slug: "test-1" # Post URL
published: "Jan 1, 2000" # Date published
tags: ["tag1", "tag2"] # Tags in YAML list format
---
# MARKDOWN GOES HERE
```
The `name` and `slug` are required, but you are free to leave out the other two. If `published` is omitted, the post won't be published on the site. If `tags` is omitted, the post won't have any tags.