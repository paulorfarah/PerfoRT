package main

import (
	"archive/zip"
	"io"
	"path/filepath"
	"strings"
)

func GetClasspathFromJar(jarPath string) []string {
	return getDependenciesFromManifestFile(jarPath)

}

// TODO
func getDependenciesFromManifestFile(jarPath string) []string {
	// Open jar as zip archive
	// log.Println(jarPath)

	r, err := zip.OpenReader(jarPath)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	var manifestContent string

	// Find and read manifest
	for _, f := range r.File {
		if strings.EqualFold(f.Name, "META-INF/MANIFEST.MF") {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}
			defer rc.Close()

			bytes, err := io.ReadAll(rc)
			if err != nil {
				panic(err)
			}

			manifestContent = string(bytes)
			break
		}
	}

	if manifestContent == "" {
		panic("Manifest not found")
	}

	// Resolve full classpath entries
	jarDir := filepath.Dir(jarPath)
	var classpathEntries []string

	lines := strings.Split(manifestContent, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "Class-Path:") {
			cpLine := strings.TrimSpace(line[len("Class-Path:"):])
			jars := strings.Fields(cpLine) // split by space
			for _, jar := range jars {
				fullPath := filepath.Join(jarDir, jar)
				classpathEntries = append(classpathEntries, fullPath)
			}
		}
	}

	// Print resolved paths
	// for _, cp := range classpathEntries {
	// 	fmt.Println("Resolved classpath jar:", cp)
	// }
	return classpathEntries
}
