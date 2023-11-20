CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS frog (
    id uuid default uuid_generate_v4() primary key,
    photoId varchar unique not null
);

CREATE TABLE IF NOT EXISTS turbo (
    id uuid default uuid_generate_v4() primary key,
    userId int not null,
    linerId varchar(3) not null,
    createdAt timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS donate (
    id uuid default uuid_generate_v4() primary key,
    username varchar(16) not null,
    sum float not null,
    createdAt timestamp default current_timestamp
);

CREATE TABLE IF NOT EXISTS anime (
    id uuid default uuid_generate_v4() primary key,
    photoId varchar unique not null
);
