package mtang

import (
	"regexp"
	"sort"
)

type route struct {
	method    map[string]Handler // maps the http method to the handler
	regPath   *regexp.Regexp     // the regular expression to associate with the path
	pathParam *pathParam         // contains relevent information in regards to the params associated with the request
}

type pathParam struct {
	pathSegmentRgx regexp.Regexp // regular expression for getting the path/param in the url. Is used for breaking up the url into segments of path/param
	positions      []int         // the index in the path for which the param/path is located
	keys           []string      // the keys associated with the param
}

type pathChunk struct {
	position []int  // the position of the path/param in the url
	pathType string //path | param
}

func buildPathChunks(paths, params [][]int) []pathChunk {
	pc := []pathChunk{}
	for _, c := range paths {
		chunk := pathChunk{
			c,
			"path",
		}
		pc = append(pc, chunk)
	}
	for _, val := range params {
		chunk := pathChunk{
			val,
			"param",
		}
		pc = append(pc, chunk)
	}
	sort.Slice(pc, func(i, j int) bool {
		return pc[i].position[0] < pc[j].position[0]
	})
	return pc
}

func buildRoute(path string) *route {
	paramRegex := regexp.MustCompile(`\/:[^/]+([^/])`)
	pathRegex := regexp.MustCompile(`\/[^:][^/]+([^/])`)
	params := paramRegex.FindAllIndex([]byte(path), -1)
	paths := pathRegex.FindAllIndex([]byte(path), -1)

	pathChunks := buildPathChunks(paths, params)

	// the regular expression path to associate with the url
	var regPath string
	if path == "/" {
		regPath = "^/$"
	} else {
		regPath = "^"
	}

	pp := &pathParam{
		*regexp.MustCompile(`[^\/:][^/]+`),
		[]int{},
		[]string{},
	}

	for i, chunk := range pathChunks {
		if chunk.pathType == "path" {
			regPath += path[chunk.position[0]:chunk.position[1]]
		}
		if chunk.pathType == "param" {
			regPath += `/[^/]+`
			pp.positions = append(pp.positions, i)
			key := regexp.MustCompile(`[^\/:][^/]+`).Find([]byte(path[chunk.position[0]:chunk.position[1]]))
			pp.keys = append(pp.keys, string(key))
		}
		if i == len(pathChunks)-1 {
			regPath += `($|[\?=&\w\d_-]+)`
		}
	}

	return &route{
		map[string]Handler{},
		regexp.MustCompile(regPath),
		pp,
	}
}
