CREATE TABLE IF NOT EXISTS bookmarks
(
    version        VARCHAR,
    url            VARCHAR NOT NULL,
    name           VARCHAR,
    desc           VARCHAR,
    date_added     INT,
    date_modified  INT,
    date_last_used INT,
    guid           VARCHAR PRIMARY KEY,
    icon           VARCHAR
);

CREATE TABLE IF NOT EXISTS folders
(
    version         VARCHAR,
    name            VARCHAR NOT NULL,
    desc            VARCHAR,
    date_added      INT,
    date_modified   INT,
    date_last_used  INT,
    guid            VARCHAR PRIMARY KEY,
    icon            VARCHAR
);

CREATE TABLE IF NOT EXISTS tag_categories
(
    name            VARCHAR PRIMARY KEY,
    desc            VARCHAR,
    icon            VARCHAR,
    date_added      INT,
    date_modified   INT
);

CREATE TABLE IF NOT EXISTS tags
(
    name            VARCHAR PRIMARY KEY,
    desc            VARCHAR,
    icon            VARCHAR,
    date_added      INT,
    date_modified   INT,
    category        REFERENCES tag_categories (name)
);

CREATE TABLE IF NOT EXISTS tag_bookmark_relations
(
    bookmark REFERENCES bookmarks (guid),
    tag REFERENCES tags (name),
    PRIMARY KEY (bookmark, tag)
);

CREATE TABLE IF NOT EXISTS tag_folder_relations
(
    folder REFERENCES folders (guid),
    tag REFERENCES tags (name),
    PRIMARY KEY (folder, tag)
);

CREATE TABLE IF NOT EXISTS subfolder_relations
(
    parent REFERENCES folders (guid),
    child PRIMARY KEY REFERENCES folders (guid)
);

CREATE TABLE IF NOT EXISTS child_relations
(
    parent REFERENCES folders (guid),
    child PRIMARY KEY REFERENCES bookmarks (guid)
);