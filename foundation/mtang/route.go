package mtang

import (
	"fmt"
	"regexp"
	"sort"
)

func buildRoute(path string) *Route {
	fmt.Println(len(path))
	paramRegex := regexp.MustCompile(`\/:[\d\w]+([^/])`)
	pathRegex := regexp.MustCompile(`\/[^:][\d\w]+([^/])`)
	params := paramRegex.FindAllIndex([]byte(path), -1)
	paths := pathRegex.FindAllIndex([]byte(path), -1)
	fmt.Println(paths)
	fmt.Println(params)

	pathChunks := []pathChunk{}
	for _, val := range paths {
		chunk := pathChunk{
			val,
			"path",
		}
		pathChunks = append(pathChunks, chunk)
		// fmt.Println(absolutePath)
	}
	for _, val := range params {
		chunk := pathChunk{
			val,
			"param",
		}
		pathChunks = append(pathChunks, chunk)
		// fmt.Println(absolutePath)
	}
	sort.Slice(pathChunks, func(i, j int) bool {
		return pathChunks[i].position[0] < pathChunks[j].position[0]
	})

	regPath := "^"
	pp := &pathParam{
		*regexp.MustCompile(`[^\/:][\w\d-_]+`),
		[]int{},
		[]string{},
	}
	for i, val := range pathChunks {
		if val.pathType == "path" {
			regPath += path[val.position[0]:val.position[1]]
		} else if val.pathType == "param" {
			regPath += `/[\d\w_-]+`
			pp.positions = append(pp.positions, i)
			key := regexp.MustCompile(`[^\/:][\w\d-_]+`).Find([]byte(path[val.position[0]:val.position[1]]))
			pp.keys = append(pp.keys, string(key))
		}
		if i == len(pathChunks)-1 {
			regPath += `($|[\?=&\w\d_-]+)`
		}
	}
	fmt.Println(regPath)

	return &Route{
		map[string]Handler{},
		regexp.MustCompile(regPath),
		pp,
	}
}
