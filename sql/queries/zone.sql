-- name: CreateZone :one
INSERT INTO zones
	(zone, time_zone)
VALUES
	(?, ?)
RETURNING *;

-- name: GetZones :many
SELECT
	id, zone, time_zone
FROM
	zones
WHERE
	zone GLOB ?
ORDER by zone;

-- name: UpdateZone :one
UPDATE zones
SET
	zone = ?,
	time_zone = ?
WHERE
	id = ?
RETURNING *;

-- name: DeleteZone :exec
DELETE FROM zones
WHERE
	id = ?;


