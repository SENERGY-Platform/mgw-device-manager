CREATE TABLE IF NOT EXISTS devices
(
    id          TEXT NOT NULL,
    ref         TEXT NOT NULL,
    name        TEXT,
    state       TEXT,
    type        TEXT NOT NULL,
    created     TEXT NOT NULL,
    updated     TEXT,
    usr_name    TEXT,
    usr_updated TEXT,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS device_attributes
(
    dev_id   TEXT    NOT NULL,
    is_usr   INTEGER NOT NULL,
    key_name TEXT    NOT NULL,
    value    TEXT,
    FOREIGN KEY (dev_id) REFERENCES devices (id) ON DELETE CASCADE ON UPDATE RESTRICT
);