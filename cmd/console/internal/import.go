/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package internal

import (
	"encoding/csv"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"shopingList/pkg/models"
	"shopingList/pkg/repositories"
	"shopingList/pkg/services"
	"strings"
	"unicode"
)

var (
	csvFile string
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import of product catalog from csv",
	Long: `import of product catalog from csv-file. 
  - Fields delimiter: ;'
  - Count fields in the file: 2
	
  Fields:
  - First column: Title of category of a product (string)
  - Second column: Title of a product (string)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("import of product catalog from the file: " + csvFile)
		importFile(csvFile)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&csvFile, "file", "f", "", "csv-file (required)")
	importCmd.MarkFlagRequired("file")
}

func importFile(filename string) {
	log.Infoln("Starting import refbook")

	db, err := openDb(appConfig.Database)
	if err != nil {
		log.Fatal("Error open database")
	}

	tx, err := db.Begin()
	if err != nil || tx == nil {
		log.Fatalln("Error open transaction; " + err.Error())
	}

	defer tx.Rollback()

	categoriesRepository := repositories.NewRefbookCategoriesRepository(tx)
	productsRepository := repositories.NewRefbookProductsRepository(tx)
	service := services.NewRefbookImportService(categoriesRepository, productsRepository)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Error open file", err)
	}
	defer file.Close()

	var productsForms []services.ProductImportForm
	categories := make(map[string]models.RefbookCategory)
	categoriesIds := make(map[string]int64)

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = 2
	record, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Error reading file", err)
	}

	for row := range record {
		categoryName := capitalize(strings.ToLower(record[row][0]))
		productName := capitalize(strings.ToLower(record[row][1]))

		_, ok := categoriesIds[categoryName]
		if !ok {
			categories[categoryName] = models.RefbookCategory{Title: categoryName}
		}

		productsForms = append(productsForms, ProductForm{Title: productName, CategoryTitle: categoryName})
	}

	if len(productsForms) == 0 {
		log.Fatalln("array of productsForms is empty")
	}

	err = service.ImportCategories(categories)
	if err != nil {
		log.Fatalln("Error import categories", err)
	}

	for _, category := range categories {
		categoriesIds[category.Title] = category.ID
	}

	err = service.ImportProducts(productsForms, categoriesIds)
	if err != nil {
		log.Fatalln("Error import categories", err)
	}

	tx.Commit()

	log.Infoln("Import completed successfully")
	log.Infof("imported categories: %v", len(categories))
	log.Infof("imported products: %v", len(productsForms))
}

func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}

	tmp := []rune(str)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}

type ProductForm struct {
	Title         string
	CategoryTitle string
}

func (s ProductForm) GetTitle() string {
	return s.Title
}

func (s ProductForm) GetCategoryTitle() string {
	return s.CategoryTitle
}
