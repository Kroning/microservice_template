package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	modulesToUseFileName = "config.yaml"
	filesToCopyFileName  = "files.yaml"
	defaultOutputPath    = "output"
	gomodPath            = "mod/go.mod"
	gosumPath            = "mod/go.sum"
)

var (
	outputPath = flag.String("output", defaultOutputPath, "Output directory of project")
)

type Config struct {
	Modules map[string][]string `yaml:"modules"`
}

type ModulesConfig struct {
	App     map[string]string `yaml:"app"`
	Modules map[string]bool   `yaml:"modules"`
}

func main() {
	flag.Parse()

	modulesCfg := loadConfig2[ModulesConfig](modulesToUseFileName)

	filesCfg := loadConfig2[Config](filesToCopyFileName)

	for module := range modulesCfg.Modules {
		fmt.Println(module, modulesCfg.Modules[module], filesCfg.Modules[module])
		if !modulesCfg.Modules[module] {
			continue
		}

		for _, filename := range filesCfg.Modules[module] {
			// can't use path.Join because we need trailing slash for directories
			newPath := *outputPath + "/" + filename

			/* don't need to check existence right now
			if _, err := os.Stat(newPath); err == nil {
				// continue
			} else if !errors.Is(err, os.ErrNotExist) {
				log.Fatal("os.Stat:", err)
			}
			*/

			if err := os.MkdirAll(filepath.Dir(newPath), 0770); err != nil {
				log.Fatal(err)
			}

			if info, err := os.Stat(filename); err == nil && info.IsDir() {
				continue
			}

			generateFile(filename, newPath, modulesCfg)
		}
	}
	generateFile(gomodPath, *outputPath+"/go.mod", modulesCfg)
	generateFile(gosumPath, *outputPath+"/go.sum", modulesCfg)
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = *outputPath
	err := cmd.Run()
	if err != nil {
		log.Fatal("fatal Run:", err.Error())
	}
}

func loadConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func loadConfig2[T any](path string) T {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg T
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func generateFile(templatePath, outputPath string, data interface{}) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("ParseFiles:", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println(data)
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatal(err)
	}
}
