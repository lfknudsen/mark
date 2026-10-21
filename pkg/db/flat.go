package db

import (
	"fmt"
	"log"
)

type FlatFolder struct {
	Name  string `sql:"name" json:"name,omitempty"`
	Guid  string `sql:"guid" json:"guid,omitempty"`
	Child string `sql:"child" json:"child,omitempty"`
}

func (ff *FlatFolder) String() string {
	return fmt.Sprintf("%s\t%s\t%s", ff.Name, ff.Guid, ff.Child)
}

func GetFoldersFlat() ([]*FlatFolder, error) {
	return GetSubfoldersFlat(ROOT_UUID)
}

func GetSubfoldersFlat(guid string) ([]*FlatFolder, error) {
	db, err := openDB()
	if err != nil {
		return []*FlatFolder{}, err
	}
	defer closeDB(db)
	rows, err := db.Query(
		`SELECT folders.name, subfolder_relations.parent, subfolder_relations.child
		FROM folders
    	JOIN subfolder_relations ON folders.guid = subfolder_relations.parent 
		WHERE subfolder_relations.parent = $1`, guid)
	if err != nil {
		log.Fatalf("Error getting subfolders: %v", err)
		return []*FlatFolder{}, err
	}
	defer closeResult(rows)
	folders := make([]*FlatFolder, 0, 25)
	for rows.Next() {
		folder := FlatFolder{}
		err = rows.Scan(&folder.Name, &folder.Guid, &folder.Child)
		if err != nil {
			return folders, err
		}
		folders = append(folders, &folder)
	}
	err = rows.Err()
	if err != nil {
		return folders, err
	}
	for _, folder := range folders {
		recursive, err := GetSubfoldersFlat(folder.Child)
		if err != nil {
			return folders, err
		}
		folders = append(folders, recursive...)
	}
	if len(folders) == 0 {
		rows, err = db.Query(`SELECT folders.name FROM folders WHERE folders.guid = $1`, guid)
		if !rows.Next() {
			return folders, rows.Err()
		}
		folder := FlatFolder{
			Guid:  guid,
			Child: "",
		}
		err = rows.Scan(&folder.Name)
		if err != nil {
			return folders, err
		}
		folders = append(folders, &folder)
	}
	return folders, err
}
