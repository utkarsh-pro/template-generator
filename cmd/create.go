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

	"github.com/spf13/viper"
	"github.com/utkarsh-pro/tempgen/helper"

	"github.com/spf13/cobra"
)

var template, name *string
var language string
var dir, languages *bool
var supportedLanguages []string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [language]",
	Short: "Creates either a file or directory by using default or a custom template",
	Long:  `Creates either a file or directory by using default or a custom template`,
	Run: func(cmd *cobra.Command, args []string) {
		supportedLanguages = viper.GetStringSlice("supportedLanguages")
		if len(args) != 1 {
			cmd.Help()
			os.Exit(1)
		}

		language = args[0]
		supported := isPresent(supportedLanguages, language)

		// Don't respect the command line argument
		if *template != "" {
			*dir, _ = helper.IsDirectory(*template)
		}

		if supported == false || *languages == true {
			if supported == false {
				fmt.Printf("Unsupported language: '%v'\n", language)
			}
			fmt.Println("Supported languages are:")
			for i, lang := range supportedLanguages {
				fmt.Printf(" %v) %v\n", i+1, lang)
			}
			fmt.Println("Consider adding the language of your choice by using:\n  tempgen add [language]")
			os.Exit(1)
		}

		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Couldn't read current working directory!")
			os.Exit(1)
		}

		if *dir == false {
			err = helper.CopyFile(getTemplatePath(), path.Join(wd, getFileName()))
		} else {
			err = helper.CopyDir(getTemplatePath(), path.Join(wd, getFileName()))
		}

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Successfuly created " + *name)
	},
}

func getTemplatePath() string {

	if *template == "" {
		if *dir == false {
			return path.Join(currentPath, "templates", language, "file", "main."+language)
		}

		return path.Join(currentPath, "templates", language, "dir")
	}

	return *template
}

func getFileName() string {
	if *dir == false {
		return *name + "." + language
	}

	return *name
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
	template = createCmd.Flags().StringP("template", "t", "", "Set custom template path")
	name = createCmd.Flags().StringP("name", "n", "main", "Specify name of the file")
	languages = createCmd.Flags().Bool("languages", false, "Set this flag to see the supported languages")
}
