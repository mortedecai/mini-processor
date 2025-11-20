CREATE TABLE IF NOT EXISTS scan_data(
    ip varchar(128) NOT NULL,
    port int NOT NULL,
    service varchar(256) NOT NULL,
    scan_date int NOT NULL,
    response text NOT NULL,
    PRIMARY KEY (ip, port, service)
)
