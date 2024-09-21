SELECT assert_latest_migration(0);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE profile
(
    id         uuid UNIQUE PRIMARY KEY           DEFAULT uuid_generate_v4(),
    email      varchar(254) UNIQUE      NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP,
);

SELECT log_migration(1);