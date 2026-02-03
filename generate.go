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

	fmt.Println("Parsing modules...")
	for module := range modulesCfg.Modules {
		//fmt.Println(module, modulesCfg.Modules[module], filesCfg.Modules[module])
		fmt.Printf("\nModule %s: ", module)
		if !modulesCfg.Modules[module] {
			fmt.Println("skipped")
			continue
		}
		fmt.Println("generating")

		for _, filename := range filesCfg.Modules[module] {
			// can't use path.Join because we need trailing slash for directories
			newPath := *outputPath + "/" + filename
			fmt.Printf("Generating file %s: ", newPath)

			/* don't need to check existence right now
			if _, err := os.Stat(newPath); err == nil {
				// continue
			} else if !errors.Is(err, os.ErrNotExist) {
				log.Fatal("os.Stat:", err)
			}
			*/

			if err := os.MkdirAll(filepath.Dir(newPath), 0775); err != nil {
				log.Fatal(err)
			}

			if info, err := os.Stat(filename); err == nil && info.IsDir() {
				continue
			}

			generateFile(filename, newPath, modulesCfg)
			fmt.Println("done")
		}
	}

	generateFile(gomodPath, *outputPath+"/go.mod", modulesCfg)
	generateFile(gosumPath, *outputPath+"/go.sum", modulesCfg)

	cmdGet := exec.Command("go", "get", "go.uber.org/mock/mockgen")
	cmdGet.Dir = *outputPath
	out, err := cmdGet.CombinedOutput()
	if err != nil {
		log.Fatalf("go get failed: %v\nOutput:\n%s", err, out)
	}

	cmdGen := exec.Command("go", "generate", "./...")
	cmdGen.Dir = *outputPath
	out, err = cmdGen.CombinedOutput()
	if err != nil {
		log.Fatalf("go generate failed: %v\nOutput:\n%s", err, out)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = *outputPath
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("go mod tidy failed: %v\nOutput:\n%s", err, out)
	}
	fmt.Println(string(out))
	fmt.Println("Successfully generated.")
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
		log.Fatalf("ParseFiles: error in %s: %s", templatePath, err.Error())
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//fmt.Println(data)
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatal(err)
	}
}
