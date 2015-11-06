package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"strings"
	"time"
)

var rootNodeAttr = "meta clear-block"

type Article struct {
	Content  []string
	metadata ArticleMetadata
}

type ArticleMetadata struct {
	Author string
	Title  string
	Tags   []string
	link   string
	Date   time.Time
}

type unsafeAccess func(node *html.Node) ([]string, error)

func formatDate(dateString string) time.Time {
	parts := strings.Split(dateString, " ")
	parts[0] = parts[0][:3]
	parts[1] = parts[1][:len(parts[1])-3]
	newDateString := strings.Join(parts, " ")
	parsedDate, err := time.Parse("Jan 2 2006 at 3:04 PM", newDateString)
	if err != nil {
		fmt.Println(err)
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}
	return parsedDate
}

func parseMetadataForArticles(article io.Reader) (metadata []ArticleMetadata, err error) {
	roots := getAllMatchingNodeAttrs(article, rootNodeAttr)
	for _, root := range roots {
		parsedMetadata, err := parseMetadata(root)
		if err == nil {
			metadata = append(metadata, parsedMetadata)
		} else {
			fmt.Println(err)
		}
	}
	return metadata, err
}

func parseMetadata(root *html.Node) (metadataFromPage ArticleMetadata, err error) {
	Link, err := getLink(root)
	if err != nil {
		fmt.Println(err)
		return metadataFromPage, err
	}
	author, err := getAuthor(root)
	if err != nil {
		fmt.Println(err)
		return metadataFromPage, err
	}
	title, err := getTitle(root)
	if err != nil {
		fmt.Println(err)
		return metadataFromPage, err
	}
	date, err := getDate(root)
	if err != nil {
		fmt.Println(err)
		return metadataFromPage, err
	}
	tags, err := getTags(root)
	if err != nil {
		fmt.Println(err)
		return metadataFromPage, err
	}
	//TODO: find better way to do this
	metadataFromPage.Author = author
	metadataFromPage.Title = title
	metadataFromPage.Date = date
	metadataFromPage.Tags = tags
	metadataFromPage.link = Link
	return metadataFromPage, err
}

func ProcessArticle(unparsedArticle io.Reader) (Article, error) {
	root := getFirstNode(unparsedArticle, rootNodeAttr)
	articleMetadata, _ := parseMetadata(root)
	article := Article{Content: getContent(root),
		metadata: articleMetadata}
	return article, nil
}

func getContent(node *html.Node) []string {
	contentNode := node.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild
	content := []string{}
	content = getAllDataFromChildren(contentNode, content, atom.P)
	return content
}

func getTags(node *html.Node) ([]string, error) {
	tags, err := getListItems(node.NextSibling.NextSibling.FirstChild.NextSibling)
	return tags, err
}

func getLink(root *html.Node) (string, error) {
	var articleLink string
	for _, attr := range root.FirstChild.NextSibling.NextSibling.NextSibling.NextSibling.NextSibling.Attr {
		if attr.Key == "href" {
			articleLink = attr.Val
		}
	}
	if articleLink == "" {
		return "", errors.New("Could not find article link")
	}
	stringsToJoin := make([]string, 2)
	stringsToJoin[0] = baseLink
	stringsToJoin[1] = strings.Split(articleLink, "#")[0]
	return strings.Join(stringsToJoin, ""), nil
}

func getDate(root *html.Node) (time.Time, error) {
	return formatDate(root.FirstChild.NextSibling.NextSibling.NextSibling.FirstChild.Data), nil
}

func getAuthor(root *html.Node) (string, error) {
	if root.FirstChild.NextSibling.FirstChild.FirstChild == nil {
		return "", errors.New("Could not find author")
	}
	return root.FirstChild.NextSibling.FirstChild.FirstChild.Data, nil
}

func getTitle(root *html.Node) (string, error) {
	var path *html.Node
	if root.PrevSibling.PrevSibling == nil {
		path = root.Parent.PrevSibling.PrevSibling.PrevSibling.PrevSibling.FirstChild
	} else {
		path = root.PrevSibling.PrevSibling.FirstChild.FirstChild
	}

	if path.Data == "" {
		return "", errors.New("No title found")
	} else {
		return path.Data, nil
	}
}

func getAllMatchingNodeAttrs(article io.Reader, nodeAttr string) []*html.Node {
	doc, err := html.Parse(article)
	if err != nil {
		fmt.Println(err)
	}
	roots := []*html.Node{}
	roots = traverseDOMTreeByAttr(doc, roots, nodeAttr)
	return roots
}

func getFirstNode(article io.Reader, nodeAttr string) *html.Node {
	return getAllMatchingNodeAttrs(article, nodeAttr)[0]
}

func getFirstMatchingChild(startNode *html.Node, nodeAttr string) *html.Node {
	nodes := []*html.Node{}
	if startNode.FirstChild == nil {
		return nil
	} else {
		returnVals := traverseDOMTreeByAttr(startNode.FirstChild, nodes, "content")
		if len(returnVals) >= 1 {
			return returnVals[0]
		} else {
			return nil
		}
	}
}

func traverseDOMTreeByAttr(node *html.Node, roots []*html.Node, attrToFind string) []*html.Node {
	// Depth first search to find the node in chronological order
	if node == nil {
		return roots
	}
	for _, attr := range node.Attr {
		if attr.Val == attrToFind {
			roots = append(roots, node)
		}
	}
	roots = traverseDOMTreeByAttr(node.FirstChild, roots, attrToFind)
	roots = traverseDOMTreeByAttr(node.NextSibling, roots, attrToFind)
	return roots
}

func getAllDataFromChildren(node *html.Node, data []string, tagTypeToFind atom.Atom) []string {
	if node == nil {
		return data
	}

	if node.Type == html.TextNode {
		text := strings.TrimSpace(node.Data)
		if len(text) != 0 {
			data = append(data, text)
		}
	}

	data = getAllDataFromChildren(node.FirstChild, data, tagTypeToFind)
	data = getAllDataFromChildren(node.NextSibling, data, tagTypeToFind)
	return data
}

func getListItems(node *html.Node) (items []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	listItem := node.FirstChild
	for listItem != nil {
		items = append(items, listItem.FirstChild.FirstChild.Data)
		listItem = listItem.NextSibling.NextSibling
	}
	return items, err
}
