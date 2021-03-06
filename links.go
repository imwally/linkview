package main

import (
	"bytes"
	"errors"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"mvdan.cc/xurls/v2"
)

type Link struct {
	NodeName string
	Text     string
	URL      string
}

func NormalizeString(s string) string {
	unescaped := html.UnescapeString(s)
	return strings.Join(strings.Fields(unescaped), " ")
}

func ImageText(s *goquery.Selection) string {
	alt, _ := s.Attr("alt")
	if alt != "" {
		return alt
	}

	title, _ := s.Attr("title")
	if title != "" {
		return title
	}

	src, ok := s.Attr("src")
	if !ok || src == "" {
		return ""
	}

	imageURL, err := url.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	return path.Base(imageURL.Path)
}

func FindLinksHTML(file io.Reader) ([]Link, error) {
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	var links []Link
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			return
		}

		node := goquery.NodeName(s)
		text := s.Text()

		s.Children().Each(func(i int, s *goquery.Selection) {
			node = goquery.NodeName(s)
			if node == "img" {
				text = ImageText(s)
				return
			}
		})

		text = NormalizeString(text)
		if text == "" {
			text = url
		}

		links = append(links, Link{node, text, url})
	})

	if len(links) == 0 {
		return nil, errors.New("no links found")
	}

	return links, nil
}

func FindLinksRegEx(file []byte) ([]Link, error) {
	var links []Link
	rxStrict := xurls.Strict()
	regexLinks := rxStrict.FindAllString(string(file), -1)
	for _, link := range regexLinks {
		links = append(links, Link{"", link, link})
	}

	if len(links) == 0 {
		return nil, errors.New("no links found")
	}

	return links, nil
}

func FindLinks(file io.Reader) ([]Link, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)

	links, err := FindLinksHTML(tee)
	if err != nil {
		b, err := ioutil.ReadAll(&buf)
		if err != nil {
			return nil, err
		}

		return FindLinksRegEx(b)
	}

	return links, err
}
