package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <repo-url>")
        os.Exit(1)
    }

    repoURL := os.Args[1]
    dir, err := os.MkdirTemp("", "repo")
    if err != nil {
        panic(err)
    }
    defer os.RemoveAll(dir) // clean up

    // Clone the given repo
    fmt.Println("Cloning repository...")
    _, err = git.PlainClone(dir, false, &git.CloneOptions{
        URL:           repoURL,
        ReferenceName: plumbing.ReferenceName("refs/heads/main"), // Adjust this as necessary
        SingleBranch:  true,
    })
    if err != nil {
        panic(err)
    }

    outputFile, err := os.Create("all_files_content.txt")
    if err != nil {
        panic(err)
    }
    defer outputFile.Close()
    writer := bufio.NewWriter(outputFile)

    // Walk through all files in the cloned directory
   err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
        return err
    }
    if !info.IsDir() {
        relativePath, err := filepath.Rel(dir, path)
        if err != nil {
            return err
        }

        // Read file content to determine its MIME type
        fileContent, err := os.ReadFile(path)
        if err != nil {
            return err
        }

        mimeType := http.DetectContentType(fileContent)
        isText := strings.HasPrefix(mimeType, "text/") || strings.Contains(mimeType, "javascript") || strings.Contains(mimeType, "json") || strings.Contains(mimeType, "xml")

        fmt.Println("Processing file:", relativePath)
        var contentToWrite string
        header := fmt.Sprintf("#########################################################\n# FILE: /%s\n#########################################################\n", relativePath)
        if isText {
            // It's a text or code file, write its content
            contentToWrite = header + string(fileContent) + "\n\n"
        } else {
            // It's not a text file, write a message instead
            contentToWrite = "This file is not a text or code file and was skipped.\n"
        }

        _, err = writer.WriteString(contentToWrite)
        if err != nil {
            return err
        }
    }
    return nil
})

if err != nil {
    panic(err)
}

err = writer.Flush()
if err != nil {
    panic(err)
}

fmt.Println("Completed. Check all_files_content.txt for the combined content.")

}
