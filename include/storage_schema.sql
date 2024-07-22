CREATE TABLE IF NOT EXISTS devices
(
    id          TEXT NOT NULL,
    ref         TEXT NOT NULL,
    name        TEXT DEFAULT '',
    type        TEXT NOT NULL,
    created     TEXT NOT NULL,
    updated     TEXT DEFAULT '',
    usr_name    TEXT DEFAULT '',
    usr_updated TEXT DEFAULT '',
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS device_attributes
(
    dev_id   TEXT    NOT NULL,
    is_usr   INTEGER NOT NULL,
    key_name TEXT    NOT NULL,
    value    TEXT DEFAULT '',
    FOREIGN KEY (dev_id) REFERENCES devices (id) ON DELETE CASCADE ON UPDATE RESTRICT
);