-- name: v1-pad-from-id
SELECT id, content, proof FROM pad WHERE BINARY id=?;

-- name: v1-insert-pad
INSERT INTO pad (id, content, proof) VALUES (?, ?, ?);

-- name: v1-update-pad
UPDATE pad SET content=?, proof=? WHERE id=?;

-- name: v1-remove-pad
DELETE FROM pad WHERE id=?;
