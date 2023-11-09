package main

import (
	"log"
	"os"
	"slices"
	"sort"
	"time"

	"path/filepath"
	"text/template"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
)

type metadata struct {
	Name      string   `yaml:"name"`
	Slug      string   `yaml:"slug"`
	Published string   `yaml:"published"`
	Tags      []string `yaml:"tags"`
}

type Post struct {
	Metadata metadata
	Date     time.Time
	Content  string
}

type Tag struct {
	Name  string
	Posts []Post
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func extractPost(path string) Post {
	file, err := os.Open(path)
	logError(err)

	var matter metadata
	content, err := frontmatter.Parse(file, &matter)
	logError(err)
	var date time.Time
	if matter.Published != "" {
		date, err = time.Parse("Jan 2, 2006", matter.Published)
		logError(err)
	}
	html := markdown.ToHTML(content, nil, nil)
	Post := Post{Metadata: matter, Date: date, Content: string(html)}

	return Post
}

func getPosts() []string {
	f, err := os.Open("posts")
	logError(err)

	files, err := f.ReadDir(0)
	logError(err)

	var paths []string
	for _, f := range files {
		filename := f.Name()
		if filename[len(filename)-3:] == ".md" {
			path := filepath.Join("posts", filename)
			paths = append(paths, path)
		}
	}
	return paths
}

func combinePosts(posts []string) []Post {
	var post_list []Post
	for _, post := range posts {
		post_content := extractPost(post)
		if post_content.Metadata.Published != "" {
			post_list = append(post_list, post_content)
		}
	}
	return post_list
}

func sortPostsByDate(posts []Post) []Post {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.Before(posts[j].Date)
	})
	slices.Reverse(posts)
	return posts
}

func buildSite(posts []Post) {
	// Menu
	err := os.Mkdir("build", os.ModePerm)
	if !os.IsExist(err) {
		logError(err)
	}
	t, err := template.ParseFiles(filepath.Join("templates", "menu.html"))
	logError(err)
	f, err := os.Create(filepath.Join("build", "index.html"))
	logError(err)
	err = t.Execute(f, posts)
	logError(err)
	err = os.Mkdir(filepath.Join("build", "post"), os.ModePerm)
	if !os.IsExist(err) {
		logError(err)
	}
	// Posts
	for _, post := range posts {
		err = os.Mkdir(filepath.Join("build", "post", post.Metadata.Slug), os.ModePerm)
		if !os.IsExist(err) {
			logError(err)
		}
		t, err := template.ParseFiles(filepath.Join("templates", "post.html"))
		logError(err)
		f, err := os.Create(filepath.Join("build", "post", post.Metadata.Slug, "index.html"))
		logError(err)
		err = t.Execute(f, post)
		logError(err)
	}
	// Tags
	tag_list := getTagList(posts)
	err = os.Mkdir(filepath.Join("build", "tag"), os.ModePerm)
	if !os.IsExist(err) {
		logError(err)
	}
	for _, tag := range tag_list {
		err = os.Mkdir(filepath.Join("build", "tag", tag.Name), os.ModePerm)
		if !os.IsExist(err) {
			logError(err)
		}
		t, err := template.ParseFiles(filepath.Join("templates", "tag.html"))
		logError(err)
		f, err := os.Create(filepath.Join("build", "tag", tag.Name, "index.html"))
		logError(err)
		err = t.Execute(f, tag)
		logError(err)
	}
}

func getTagList(posts []Post) []Tag {
	var tag_list []string
	for _, post := range posts {
		tags := post.Metadata.Tags
		for _, tag := range tags {
			if !slices.Contains(tag_list, tag) {
				tag_list = append(tag_list, tag)
			}
		}
	}
	var tags []Tag
	for _, tag_list_entry := range tag_list {
		var tag_entry Tag
		tag_entry.Name = tag_list_entry
		for _, post := range posts {
			if slices.Contains(post.Metadata.Tags, tag_list_entry) {
				tag_entry.Posts = append(tag_entry.Posts, post)
			}
		}
		tags = append(tags, tag_entry)
	}
	return tags
}

func main() {
	posts := getPosts()
	post_list := combinePosts(posts)
	post_list = sortPostsByDate(post_list)
	buildSite(post_list)
}
