package pgutils_test

import (
	pg "nms/lib/pgutils"
	"testing"
)

type Author struct {
	Name  string `pg:",pk"`
	Books []Book `pg:"rel:has-many"`
}

type Book struct {
	Id         int64
	AuthorName string
	Name       string
}

func TestDB(t *testing.T) {
	db, err := pg.NewClient()
	if err != nil {
		// running db
		// docker run --rm -d -p 5432:5432 --env POSTGRES_USER=user --env POSTGRES_PASSWORD=pass --env POSTGRES_DB=nms  postgres
		t.Error("connect db fail")
	}
	defer db.Close()

	a := Author{
		Name: "austin",
		Books: []Book{
			{
				Name:       "english",
				AuthorName: "austin",
			},
			{
				Name:       "chinese",
				AuthorName: "austin",
			},
		},
	}

	d, _ := db.GetDB()
	_, err = d.Exec(`DROP TABLE IF EXISTS authors`)
	if err != nil {
		t.Error(err)
	}
	_, err = d.Exec(`DROP TABLE IF EXISTS books`)
	if err != nil {
		t.Error(err)
	}

	err = db.CreateTable(&Author{}, pg.CreateTableOpt{IfNotExists: true})
	if err != nil {
		t.Error("create config metric state table fail: ", err)
		return
	}
	err = db.CreateTable(&Book{}, pg.CreateTableOpt{IfNotExists: true})
	if err != nil {
		t.Error("create config metric state table fail: ", err)
		return
	}

	_, err = d.Model(&a).Relation("Books").Insert()
	if err != nil {
		t.Error(err)
	}

}
