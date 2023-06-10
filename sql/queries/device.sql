-- name: CreateDevice :one
INSERT INTO devices
	(device, cluster_id, appliance_id, model_id, environment_id)
VALUES
	(?,
	(select id FROM clusters c where c.cluster = ?),
	(select id FROM appliances a where a.appliance = ?),
	(select id FROM models m where m.model = ?),
	(select id FROM environments e where e.environment = ?))
RETURNING *;

-- name: GetDevices :many
SELECT
	d.id,
	d.device,
	z.zone,
	c.cluster,
	a.appliance,
	m.model,
	e.environment
FROM
	devices d
	LEFT JOIN zones z ON d.zone_id = z.id
	LEFT JOIN clusters c ON d.cluster_id = c.id
	LEFT JOIN appliances a ON d.appliance_id = a.id
	LEFT JOIN models m ON d.model_id = m.id
	LEFT JOIN environments e ON d.environment_id = e.id
WHERE
	z.zone GLOB ?
AND
	c.cluster GLOB ?
AND
	d.device GLOB ?
AND
	c.zone_id = z.id
AND
	d.cluster_id = c.id; 

-- name: UpdateDevice :one
UPDATE
       devices
SET
       device = ?,
       cluster_id = (select id FROM clusters c where c.cluster = ?),
       appliance_id = (select id FROM appliances a where a.appliance = ?),
       model_id = (select id FROM models m where m.model = ?),
       environment_id = (select id FROM environments e where e.environment = ?)
WHERE
       devices.id = ?
RETURNING *;

-- name: DeleteDevice :exec
DELETE FROM
	devices
WHERE
	id = ?;
