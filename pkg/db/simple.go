package db

// SimpleFolder is a data transfer object used in front-end application
// where only a small subset of directory information is required.
type SimpleFolder struct {
	Name       string         `sql:"name" json:"name,omitempty"`
	Subfolders []SimpleFolder `sql:"-" json:"subfolders,omitempty"`
}

func GetFolders() ([]*SimpleFolder, error) {
	db, err := openDB()
	if err != nil {
		return []*SimpleFolder{}, err
	}
	defer closeDB(db)
	rows, err := db.Query(`SELECT folders.name, folders.guid
						   FROM folders 
                           EXCEPT SELECT DISTINCT (child) 
                           FROM subfolder_relations
                           JOIN folders ON subfolder_relations.child = folders.guid`)
	if err != nil {
		return []*SimpleFolder{}, err
	}
	defer closeResult(rows)
	folders := make([]*SimpleFolder, 0, 25)
	for rows.Next() {
		folder := SimpleFolder{}
		var guid string
		err = rows.Scan(&folder.Name, &guid)
		if err != nil {
			return folders, err
		}
		folder.Subfolders, err = GetSubfolders(guid)
		folders = append(folders, &folder)
		if err != nil {
			return folders, err
		}
	}
	err = rows.Err()
	return folders, err
}

func GetSubfolders(guid string) ([]SimpleFolder, error) {
	db, err := openDB()
	if err != nil {
		return []SimpleFolder{}, err
	}
	defer closeDB(db)
	rows, err := db.Query(
		`SELECT folders.name, folders.guid
		FROM folders
    	JOIN subfolder_relations ON folders.guid = subfolder_relations.child
		WHERE subfolder_relations.parent = $1`, guid)
	if err != nil {
		return []SimpleFolder{}, err
	}
	defer closeResult(rows)
	folders := make([]SimpleFolder, 0, 25)
	for rows.Next() {
		folder := SimpleFolder{}
		var subGuid string
		err = rows.Scan(&folder.Name, &subGuid)
		if err != nil {
			return folders, err
		}
		folder.Subfolders, err = GetSubfolders(subGuid)
		folders = append(folders, folder)
		if err != nil {
			return folders, err
		}
	}
	err = rows.Err()
	return folders, err
}
