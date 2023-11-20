CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS frog (
    id uuid DEFAULT uuid_generate_v4(),
    photoId VARCHAR UNIQUE NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS turbo (
    id uuid DEFAULT uuid_generate_v4(),
    userId INT NOT NULL,
    linerId VARCHAR(3) NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
