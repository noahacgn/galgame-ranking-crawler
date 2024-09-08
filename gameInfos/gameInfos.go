package gameInfos

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type GameInfo struct {
	Title   string
	Rank    int
	Date    time.Time
	Chinese bool
	Point   float64
}

func Extract(url string) ([]GameInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var gameInfos []GameInfo
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key != "class" || a.Val != "game-info" {
					continue
				}

				s := strings.Split(getTextContent(n), "|")
				var gameInfo GameInfo
				if len(s) == 35 {
					gameInfo.Title = s[2]
					gameInfo.Rank, _ = strconv.Atoi(strings.TrimPrefix(s[12], "Rank: "))
					gameInfo.Date, _ = time.Parse("2006-01-02", strings.TrimSuffix(s[15], "发行"))
					if strings.HasPrefix(s[18], "有") {
						gameInfo.Chinese = true
					}
					gameInfo.Point, _ = strconv.ParseFloat(s[25], 64)
				} else {
					gameInfo.Title = s[2]
					gameInfo.Rank, _ = strconv.Atoi(strings.TrimPrefix(s[9], "Rank: "))
					gameInfo.Date, _ = time.Parse("2006-01-02", strings.TrimSuffix(s[12], "发行"))
					if strings.HasPrefix(s[15], "有") {
						gameInfo.Chinese = true
					}
					gameInfo.Point, _ = strconv.ParseFloat(s[22], 64)
				}

				gameInfos = append(gameInfos, gameInfo)
			}

		}
	}
	forEachNode(doc, visitNode, nil)
	return gameInfos, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += getTextContent(c) + "|"
	}
	return strings.TrimSpace(result)
}
