package sqlite

import "log"

// function to add a category
func AddCategory(name string) {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()

	//insert category into table
	statement, err := db.Prepare("INSERT INTO categories (category) VALUES (?)")
	if err != nil {
		log.Println(err)
		return
	}
	statement.Exec(name)
}

// function to get all categories to category struct
func GetCategories() []string {
	db, err := OpenDb()
	if err != nil {
		log.Println(err)
		return nil
	}

	defer db.Close()

	//get all categories
	rows, err := db.Query("SELECT * FROM categories")
	if err != nil {
		log.Println(err)
		return nil
	}
	var id int
	var category string
	var categories []Category
	for rows.Next() {
		rows.Scan(&id, &category)
		categories = append(categories, Category{Id: id, Category: category})
	}

	var categoryNames []string
	for _, category := range categories {
		categoryNames = append(categoryNames, category.Category)
	}

	return categoryNames
}

type Category struct {
	Id       int
	Category string
}
