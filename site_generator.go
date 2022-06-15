package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

func main() {
	// data, err := os.ReadFile("site/blog/chitri_2022_kickoff.md")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// html := md_to_html(data)

	// layout, _ := os.ReadFile("site/templates/blog.html")

	// blog_post, _ := create_blog_post_html(string(html), string(layout))

	// os.WriteFile("test.html", []byte(blog_post), os.ModePerm)

	// fm := "---\ntitle: test title\ntags: triathlon, sports\n---"
	// var fm_struct frontmatter

	// for _, line := range strings.Split(fm, "\n") {
	// 	if strings.Contains(line, "title:") {
	// 		fm_struct.title = strings.Trim(strings.Split(line, ":")[1], " ")
	// 	} else if strings.Contains(line, "tags:") {
	// 		for _, tag := range strings.Split(strings.Split(line, ":")[1], ",") {
	// 			fm_struct.tags = append(fm_struct.tags, strings.Trim(tag, " "))
	// 		}
	// 	}

	// }

	// fmt.Println(fm_struct)

	deploy_path := "deploy"
	site_path := "site"

	if err := recreate_deploy(deploy_path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := generate_deploy(site_path, deploy_path); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type frontmatter struct {
	title string
	tags  []string
	uri   string
}

func generate_deploy(site_path string, deploy_path string) error {
	walk_func := func(path string, info os.FileInfo, err error) error {
		if path == site_path {
			return nil
		}

		new_deploy_path := deploy_path + "/" + strings.Replace(path, site_path+"/", "", 1)

		if info.IsDir() {
			if err := os.MkdirAll(new_deploy_path, os.ModePerm); err != nil {
				return err
			}

		} else {
			if strings.Contains(path, "blog/") && strings.Contains(path, ".md") {
				if err := generate_blog_post_file(path, deploy_path); err != nil {
					return err
				}
			} else {
				if err := copy_file(path, new_deploy_path); err != nil {
					return err
				}
			}
		}

		return nil
	}

	if err := filepath.Walk(site_path, walk_func); err != nil {
		return err
	}

	return nil
}

func generate_blog_post_file(path string, deploy_path string) error {
	data, _ := os.ReadFile(path)
	md, front_matter, _ := remove_frontmatter(string(data))
	html := md_to_html([]byte(md))
	fm_struct, _ := string_to_frontmatter(front_matter)
	layout, _ := os.ReadFile("site/templates/blog.html")
	blog_post, _ := create_blog_post_html(string(html), string(layout))
	// take path and get filename without extension
	// file_name := strings.Replace(filepath.Base(path), ".md", ".html", 1)
	file_name := fm_struct.uri + ".html"
	if err := os.WriteFile(deploy_path+"/blog/"+file_name, []byte(blog_post), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func copy_file(oldpath string, newpath string) error {
	src, err := os.Open(oldpath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(newpath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func recreate_deploy(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func md_to_html(data []byte) []byte {
	html := markdown.ToHTML(data, nil, nil)

	return html
}

func remove_frontmatter(content string) (string, string, error) {
	var new_lines []string
	var frontmatter []string
	lines := strings.Split(content, "\n")
	start_append := false

	for i, line := range lines {
		if start_append == true {
			new_lines = append(new_lines, line)
		} else {
			frontmatter = append(frontmatter, line)
		}

		if i > 0 && strings.Contains(line, "---") {
			start_append = true
		}
	}

	return strings.Join(new_lines, "\n"), strings.Join(frontmatter, "\n"), nil
}

func create_blog_post_html(blog_post string, template string) (string, error) {
	var blog_lines []string
	lines := strings.Split(template, "\n")
	for _, line := range lines {
		if strings.Contains(line, "{{.BlogPost}}") {
			new := strings.Replace(line, "{{.BlogPost}}", blog_post, 1)
			blog_lines = append(blog_lines, new)
		} else {
			blog_lines = append(blog_lines, line)
		}
	}

	return strings.Join(blog_lines, "\n"), nil
}

func string_to_frontmatter(front_matter string) (frontmatter, error) {
	var fm_struct frontmatter

	for _, line := range strings.Split(front_matter, "\n") {
		if strings.Contains(line, "title:") {
			fm_struct.title = strings.Trim(strings.Split(line, ":")[1], " ")
		} else if strings.Contains(line, "tags:") {
			for _, tag := range strings.Split(strings.Split(line, ":")[1], ",") {
				fm_struct.tags = append(fm_struct.tags, strings.Trim(tag, " "))
			}
		} else if strings.Contains(line, "uri:") {
			fm_struct.uri = strings.Trim(strings.Split(line, ":")[1], " ")
		}
	}

	return fm_struct, nil
}
