CREATE TRIGGER IF NOT EXISTS auto_root_child AFTER INSERT ON bookmarks
BEGIN
    INSERT INTO child_relations (parent, child) VALUES ('0', NEW.guid);
END;

CREATE TRIGGER IF NOT EXISTS auto_root_subfolder AFTER INSERT ON folders
BEGIN
    INSERT INTO subfolder_relations (parent, child) VALUES ('0', NEW.guid);
END;