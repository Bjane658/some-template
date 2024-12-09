package main

import (
	"bufio"
	"fmt"
	"github.com/alexflint/go-arg"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var args struct {
	Add      bool     `arg:"-a,--add" help:"add template"`
	View     bool     `arg:"-v,--view" help:"view template"`
	Edit     bool     `arg:"-e,--edit" help:"edit template"`
	Keywords []string `arg:"positional"`
}
var Config struct {
	TemplateRoot string
}

func init() {
	homeDir := os.Getenv("HOME")
	config, err := loadConfig(filepath.Join(homeDir, ".somerc"))
	if err != nil {
		fmt.Println("Could not load config: ", err)
		return
	}
	templateDir := config["template_root"]
	if len(templateDir) > 0 {
		Config.TemplateRoot = templateDir
	} else {
		Config.TemplateRoot = "."
	}
}

func main() {
	arg.MustParse(&args)
	if args.Add {
		addTemplate(args.Keywords)
		return
	}
	if args.View {
		viewTemplate(args.Keywords)
		return
	}
	if args.Edit {
		editTemplate(args.Keywords)
		return
	}
	applyTemplate(args.Keywords)
}

func loadConfig(filename string) (map[string]string, error) {
	// Open the .somerc file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	// Create a map to hold the configuration
	config := make(map[string]string)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Split the line at the '=' symbol
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		} else {
			log.Printf("Skipping invalid line: %s", line)
		}
	}

	// Check for errors while reading
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return config, nil
}

func applyTemplate(keywords []string) {
	templateFilepath := getFilepathFrom(keywords)
	sourceFile, err := os.Open(templateFilepath)
	if err != nil {
		fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	filename := getFilenameFrom(keywords)
	destinationPath := filepath.Join(".", filename)
	// Create the destination file
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the content from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		fmt.Errorf("failed to copy file content: %w", err)
	}

	// Sync to ensure all data is written to disk
	err = destinationFile.Sync()
	if err != nil {
		fmt.Errorf("failed to sync destination file: %w", err)
	}
}
func addTemplate(keywords []string) {
	fmt.Println("Adding template: " + keywords[len(keywords)-1])
	filePath := writeFile(keywords[:])
	openTerminal(filePath, "vim")
}

func getFolderpathFromKeywords(keywords []string) string {
	folders := keywords[:len(keywords)-1]
	folderPathElements := append([]string{Config.TemplateRoot, "templates"}, folders...)
	return filepath.Join(folderPathElements...)
}

func getFilepathFrom(keywords []string) string {
	folders := keywords[:len(keywords)-1]
	folderPathElements := append([]string{Config.TemplateRoot, "templates"}, folders...)
	folderPath := filepath.Join(folderPathElements...)
	filename := getFilenameFrom(keywords)
	return filepath.Join(folderPath, filename)
}

func editTemplate(keywords []string) {
	filePath := getFilepathFrom(keywords)
	openTerminal(filePath, "vim")
}

func viewTemplate(keywords []string) {
	filePath := getFilepathFrom(keywords)
	openTerminal(filePath, "less")
}

func printFile(filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("failed to read file: ", err)
		return
	}
	fileContent := string(file)
	if len(fileContent) == 0 {
		fmt.Println("Template is empty")
	}
	_, err = io.Copy(os.Stdout, strings.NewReader(fileContent))
	if err != nil {
		fmt.Println("failed to read file content: ", err)
	}
}

func createFolderStructureFrom(keywords []string) {
	folderPath := getFolderpathFromKeywords(keywords)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

}

func writeFile(keywords []string) string {
	createFolderStructureFrom(keywords)
	filePath := getFilepathFrom(keywords)
	f, err := os.Create(filePath)
	if err != nil {
		log.Printf("1")
		log.Fatal(err)
	}
	f.Close()
	return filePath
}

func openTerminal(filePath string, editor string) {
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Printf("2")
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
	}
}

func getFilenameFrom(keywords []string) string {
	return keywords[len(keywords)-1] + languageToFileExtension(keywords[0])
}

func languageToFileExtension(language string) string {
	switch language {
	case "java":
		return ".java"
	case "python":
		return ".py"
	case "go":
		return ".go"
	}
	return ""
}
