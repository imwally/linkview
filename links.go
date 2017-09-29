package main

import (
	"errors"
	"io"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Link struct {
	NodeName string
	Text     string
	URL      string
}

func NormalizeString(s string) string {
	return strings.Join(strings.Fields(s), " ")
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
		return "NO TEXT"
	}

	imageURL, err := url.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	return path.Base(imageURL.Path)
}

func FindLinks(file io.Reader) ([]Link, error) {
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
		text := NormalizeString(s.Text())

		s.Children().Each(func(i int, s *goquery.Selection) {
			node = goquery.NodeName(s)
			if node == "img" {
				text = ImageText(s)
				return
			}
		})

		if text == "" {
			text = "NO TEXT"
		}

		links = append(links, Link{node, text, url})
	})

	if len(links) == 0 {
		return nil, errors.New("no links found")
	}

	return links, nil
}
