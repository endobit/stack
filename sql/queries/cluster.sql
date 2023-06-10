-- name: CreateCluster :one
INSERT INTO clusters
	(cluster, zone)
VALUES
	(?, (select id from zones z where z.zone = ?))
RETURNING *;

-- name: GetClusters :many
SELECT
	c.id,
	z.zone,
	c.cluster
FROM
	zones z, clusters c
WHERE
	z.zone GLOB ?
AND
	c.cluster GLOB ?
AND
	c.zone_id = z.id;


-- name: UpdateCluster :one
UPDATE clusters
SET
	cluster = ?,
	zones = (select id from zones z where z.zone = ?)
WHERE
	clusters.id = ?
RETURNING *;

-- name: DeleteCluster :exec
DELETE FROM clusters
WHERE
	id = ?;


