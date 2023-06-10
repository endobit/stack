-- attributes for all levels
--
-- global
-- zone
-- cluster
-- model
-- environment
-- device

CREATE TABLE IF NOT EXISTS global_attributes (
       id    	     	   INTEGER PRIMARY KEY,
       key   	    	   TEXT    NOT NULL,
       value 	    	   TEXT    NOT NULL,
       is_protected 	   INTEGER NOT NULL CHECK (is_protected IN (0,1))
);

CREATE TABLE IF NOT EXISTS zone_attributes (
       id    	    	   INTEGER PRIMARY KEY,
       zone_id	    	   INTEGER NOT NULL,
       key 	    	   TEXT    NOT NULL,
       value 	    	   TEXT    NOT NULL,
       is_protected 	   INTEGER NOT NULL CHECK (is_protected IN (0,1)),
       FOREIGN KEY (zone_id) REFERENCES zones(id)
);

CREATE TABLE IF NOT EXISTS cluster_attributes (
       id    	    	   INTEGER PRIMARY KEY,
       cluster_id   	   INTEGER NOT NULL,
       key 	    	   TEXT    NOT NULL,
       value 	    	   TEXT    NOT NULL,
       is_protected 	   INTEGER NOT NULL CHECK (is_protected IN (0,1)),
       FOREIGN KEY (cluster_id) REFERENCES clusters(id)
);

CREATE TABLE IF NOT EXISTS model_attributes (
       id    	    	   INTEGER PRIMARY KEY,
       model_id     	   INTEGER NOT NULL,
       key 	    	   TEXT    NOT NULL,
       value 	    	   TEXT    NOT NULL,
       is_protected 	   INTEGER NOT NULL CHECK (is_protected IN (0,1)),
       FOREIGN KEY (model_id) REFERENCES models(id)
);

CREATE TABLE IF NOT EXISTS environment_attributes (
       id    	      	   INTEGER PRIMARY KEY,
       environment_id 	   INTEGER NOT NULL,
       key 	      	   TEXT    NOT NULL,
       value 	      	   TEXT    NOT NULL,
       is_protected   	   INTEGER NOT NULL CHECK (is_protected IN (0,1)),
       FOREIGN KEY (environment_id) REFERENCES environments(id)
);

CREATE TABLE IF NOT EXISTS device_attributes (
       id    	    	   INTEGER PRIMARY KEY,
       device_id	   INTEGER NOT NULL,
       key 	    	   TEXT    NOT NULL,
       value 	    	   TEXT    NOT NULL,
       is_protected 	   INTEGER NOT NULL CHECK (is_protected IN (0,1)),
       FOREIGN KEY (device_id) REFERENCES devices(id)
);

-- zone
-- cluster
-- model
-- environment

CREATE TABLE IF NOT EXISTS zones (
       id           	   INTEGER PRIMARY KEY,
       zone 	    	   TEXT    UNIQUE NOT NULL,
       time_zone    	   TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS clusters (
       id    	    	   INTEGER PRIMARY KEY,
       zone_id 	    	   INTEGER NOT NULL,
       cluster 	    	   TEXT    NOT NULL,
       FOREIGN KEY  (zone_id) REFERENCES zones(id)
);

CREATE TABLE IF NOT EXISTS appliances (
       id    	    	   INTEGER PRIMARY KEY,
       appliance    	   TEXT    UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS models (
       id    	    	   INTEGER PRIMARY KEY,
       model		   TEXT    UNIQUE NOT NULL,
       arch 	    	   TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS environments (
       id    	    INTEGER PRIMARY KEY,
       environment  TEXT    UNIQUE NOT NULL
);

-- networks

CREATE TABLE IF NOT EXISTS networks (
       id    	    	   INTEGER PRIMARY KEY,
       network	    	   TEXT    NOT NULL,
       address 	    	   TEXT    NOT NULL,
       gateway 	    	   TEXT,
       mtu 	    	   INTEGER NOT NULL,
       is_pxe 	    	   INTEGER NOT NULL CHECK (is_pxe IN (0,1))
);

CREATE TABLE IF NOT EXISTS zone_networks (
       network_id    	   INTEGER NOT NULL,
       zone_id 	    	   INTEGER NOT NULL,
       FOREIGN KEY (network_id) REFERENCES networks(id),
       FOREIGN KEY (zone_id) 	REFERENCES zones(id)
);

CREATE TABLE IF NOT EXISTS cluster_networks (
       network_id   	   INTEGER NOT NULL,
       cluster_id 	    INTEGER NOT NULL,
       FOREIGN KEY (network_id) REFERENCES networks(id),
       FOREIGN KEY (cluster_id) REFERENCES clusters(id)
);

-- devices
--
-- hosts
-- switches

CREATE TABLE IF NOT EXISTS devices (
       id            	   INTEGER PRIMARY KEY,
       device 	    	   TEXT    NOT NULL,
       cluster_id	   INTEGER,
       appliance_id	   INTEGER,
       model_id		   INTEGER,
       environment_id	   INTEGER,
       FOREIGN KEY (zone_id)	    REFERENCES zones(id),
       FOREIGN KEY (cluster_id)     REFERENCES clusters(id),
       FOREIGN KEY (appliance_id)   REFERENCES appliances(id),
       FOREIGN KEY (model_id) 	    REFERENCES models(id),
       FOREIGN KEY (environment_id) REFERENCES environments(id)
);

CREATE TABLE IF NOT EXISTS hosts (
       id     	     	   INTEGER PRIMARY KEY,
       device_id	   INTEGER,
       rack   	    	   TEXT,
       rank   	    	   INTEGER,
       type   	    	   TEXT,
       FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE TABLE IF NOT EXISTS switches (
       id     	     	   INTEGER PRIMARY KEY,
       device_id	   INTEGER,
       rack   	    	   TEXT,
       rank   	    	   INTEGER,
       FOREIGN KEY (device_id) REFERENCES devices(id)
);

-- NICs

CREATE TABLE IF NOT EXISTS network_interfacs (
       id     	     	   INTEGER PRIMARY KEY,
       device_id	   INTEGER,
       network_id	   INTEGER,
       ip   	     	   TEXT,
       mac   	     	   TEXT,
       is_dhcp 	     	   INTEGER NOT NULL CHECK (is_dhcp IN (0,1)),
       is_management 	   INTEGER NOT NULL CHECK (is_management IN (0,1)),
       FOREIGN KEY (device_id) 	REFERENCES devices(id),
       FOREIGN KEY (network_id) REFERENCES networks(id)
);
