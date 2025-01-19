package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
)

type ArticleData struct {
	Title string
	Link  string
}

type IndexPageData struct {
	Articles []ArticleData
}

func main() {
	mainTmpl, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	articlesDir := "articles"
	outputDir := "dist"
	indexData := IndexPageData{}

	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(articlesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			mdContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			htmlContent := markdown.ToHTML(mdContent, nil, nil)

			outputFileName := strings.TrimSuffix(info.Name(), ".md") + ".html"
			outputFilePath := filepath.Join(outputDir, outputFileName)

			err = os.WriteFile(outputFilePath, htmlContent, os.ModePerm)
			if err != nil {
				return err
			}

			fmt.Printf("Created %v\n", outputFileName)

			indexData.Articles = append(indexData.Articles, ArticleData{
				Title: strings.TrimSuffix(info.Name(), ".md"),
				Link:  outputFileName,
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	indexOutputFile := filepath.Join(outputDir, "index.html")
	indexFile, err := os.Create(indexOutputFile)
	if err != nil {
		panic(err)
	}
	defer indexFile.Close()

	fmt.Printf("Index data: %+v\n", indexData)

	err = mainTmpl.Execute(indexFile, indexData)
	if err != nil {
		panic(err)
	}

	fmt.Println("Site generated successfully!")
}
