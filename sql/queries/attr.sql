--
-- Global
--

-- name: CreateGlobalAttribute :one
INSERT INTO global_attributes
       (key, value, is_protected)
VALUES
       (?, ?, ?)
RETURNING *;

-- name: GetGlobalAttributes :many
SELECT
	id, key, value, is_protected
FROM
	global_attributes
WHERE
	key GLOB ?
ORDER BY key ASC;

-- name: DeleteGlobalAttribute :exec
DELETE FROM
	global_attributes
WHERE
	id = ?;

-- name: UpdateGlobalAttribute :one
UPDATE
	global_attributes
SET
	key = ?,
	value = ?,
	is_protected = ?
WHERE
	id = ?
RETURNING *;

--
-- Zone
--

-- name: CreateZoneAttribute :one
INSERT INTO zone_attributes
       (zone_id, key, value, is_protected)
VALUES
       ((SELECT id FROM zones WHERE zones.zone = ?), ?, ?, ?)
RETURNING *;

-- name: GetZoneAttributes :many
SELECT
	a.id, z.zone, a.key, a.value, a.is_protected
FROM
	zone_attributes a, zones z
WHERE
	z.zone GLOB ?
AND
	a.key GLOB ?
AND
	a.zone_id=z.id
ORDER BY z.zone, a.key ASC;

-- name: DeleteZoneAttribute :exec
DELETE FROM
	zone_attributes
WHERE
	id = ?;

-- name: UpdateZoneAttribute :one
UPDATE
	zone_attributes
SET
	key = ?,
	value = ?,
	is_protected = ?
WHERE
	id = ?
RETURNING *;

--
-- Cluster
--

-- name: CreateClusterAttribute :one
INSERT INTO cluster_attributes
       (cluster_id, key, value, is_protected)
VALUES
       ((SELECT clusters.id
         FROM clusters, zones
	 WHERE zones.zone = ?
	 AND clusters.cluster = ?
	 AND cluster.zone_id = zones.id), ?, ?, ?)
RETURNING *;
       
-- name: GetClusterAttributes :many
SELECT
	a.id, z.zone, c.cluster, a.key, a.value, a.is_protected
FROM
	cluster_attributes a, zones z, clusters c
WHERE
	z.zone GLOB ?
AND
	c.cluster GLOB ?
AND
	a.key GLOB ?
AND
	c.zone_id=z.id
AND	
	a.cluster_id=c.id
ORDER BY z.zone, c.cluster, a.key ASC;

-- name: DeleteClusterAttribute :exec
DELETE FROM
	cluster_attributes
WHERE
	id = ?;

-- name: UpdateClusterAttribute :one
UPDATE
	cluster_attributes
SET
	cluster = (SELECT id FROM clusters c WHERE c.cluster = ?),
	key = ?,
	value = ?,
	is_protected = ?
WHERE
	cluster_attributes.id = ?
RETURNING *;

--
-- Model
--

-- name: CreateModelAttribute :one
INSERT INTO model_attributes
       (model_id, key, value, is_protected)
VALUES
       ((SELECT id FROM models WHERE models.model = ?), ?, ?, ?)
RETURNING *;
       
-- name: GetModelAttributes :many
SELECT
	m.model, a.key, a.value, a.is_protected
FROM
	model_attributes a, models m
WHERE
	m.model GLOB ?
AND
	a.key GLOB ?
AND
	a.model_id=m.id
ORDER BY m.model, a.key ASC;

-- name: DeleteModelAttribute :exec
DELETE FROM
	model_attributes
WHERE
	id = ?;

-- name: UpdateModelAttribute :one
UPDATE
	model_attributes
SET
	model = (SELECT id FROM models m WHERE m.model = ?),
	key = ?,
	value = ?,
	is_protected = ?
WHERE
	model_attributes.id = ?
RETURNING *;

--
-- Environment
--

-- name: CreateEnvironmentAttribute :one
INSERT INTO environment_attributes
       (environment_id, key, value, is_protected)
VALUES
       ((SELECT id FROM environments WHERE environments.environment = ?), ?, ?, ?)
RETURNING *;

-- name: GetEnvironmentAttributes :many
SELECT
	a.id, e.environment, a.key, a.value, a.is_protected
FROM
	environment_attributes a, environments e
WHERE
	e.environment GLOB ?
AND
	a.key GLOB ?
AND
	a.environment_id=e.id
ORDER BY e.environment, a.key ASC;

-- name: DeleteEnvironmentAttribute :exec
DELETE FROM
	environment_attributes
WHERE
	id = ?;

-- name: UpdateEnvironmentAttribute :one
UPDATE
	environment_attributes
SET
	environment = (SELECT id FROM environments e WHERE e.environment = ?),
	key = ?,
	value = ?,
	is_protected = ?
WHERE
	environment_attributes.id = ?
RETURNING *;

--
-- Device
--

-- name: CreateDeviceAttribute :one
INSERT INTO device_attributes
       (device_id, key, value, is_protected)
VALUES
       ((SELECT devices.id
         FROM devices, zones, clusters
	 WHERE devices.device = ?
	 AND zones.zone = ?
	 AND clusters.cluster = ?
	 AND devices.cluster_id = cluster.id
	 AND cluster.zone_id = zones.id
	 ), ?, ?, ?)
RETURNING *;
       
-- name: GetDeviceAttributes :many
SELECT
	a.id, d.device, a.key, a.value, a.is_protected
FROM
	device_attributes a, devices d, zones z, clusters c
WHERE
	z.zone GLOB ?
AND
	c.cluster GLOB ?
AND
	d.device GLOB ?
AND
	a.key GLOB ?
AND
	a.device_id=d.id
AND
	d.cluster_id = clusters.cluster
AND
	c.zone_id = zones.zone
ORDER BY d.device ASC;

-- name: DeleteDeviceAttribute :exec
DELETE FROM
	device_attributes
WHERE
	id = ?;

-- name: UpdateDeviceAttribute :one
UPDATE
	device_attributes
SET
	device = (SELECT id FROM devices d WHERE d.device = ?),
	key = ?,
	value = ?,
	is_protected = ?
WHERE
	device_attributes.id = ?
RETURNING *;


