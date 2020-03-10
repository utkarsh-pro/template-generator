/*
Copyright © 2020 Utkarsh Srivastava <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/utkarsh-pro/tempgen/copy"
)

var supportedLanguages = []string{"cpp", "js", "go", "py"}
var currentPath = getCurrentPath()
var template, language, name *string
var dir *bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create command is used to create either a file or directory by using default or a custom template",
	Long:  `create command is used to create either a file or directory by using default or a custom template`,
	Run: func(cmd *cobra.Command, args []string) {
		supported := isPresent(supportedLanguages, *language)

		if supported == false {
			fmt.Printf("Unsupported language: '%v'\nSupported languages are:\n", *language)
			for i, lang := range supportedLanguages {
				fmt.Printf(" %v) %v\n", i+1, lang)
			}
			os.Exit(1)
		}

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Couldn't read current working directory!")
			os.Exit(1)
		}

		if *dir == false {
			err := copy.File(path.Join(currentPath, "templates", *language, "file", "main"), path.Join(wd, *name))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := copy.Dir(path.Join(currentPath, "templates", *language, "dir"), path.Join(wd, *name))
			if err != nil {
				fmt.Println(err)
			}
		}

		fmt.Println("Successfuly created " + *name)
	},
}

func isPresent(arr []string, val string) bool {
	flag := false
	for _, el := range arr {
		if val == el {
			flag = true
			break
		}
	}

	return flag
}

func getCurrentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	return exPath
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	dir = createCmd.Flags().BoolP("directory", "d", false, "Default is false, set to true to specify if a directory is to be created")
	template = createCmd.Flags().StringP("template", "t", "", "Set custom template, accepts local location or URL")
	language = createCmd.Flags().StringP("language", "l", "cpp", "Set programming language")
	name = createCmd.Flags().StringP("name", "n", "main.cpp", "Specify name of the file")
}
