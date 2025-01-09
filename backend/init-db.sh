#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS artists (
        id SERIAL PRIMARY KEY,
        name TEXT UNIQUE NOT NULL
    );

    CREATE TABLE IF NOT EXISTS formats (
        id SERIAL PRIMARY KEY,
        name TEXT UNIQUE NOT NULL
    );

    CREATE TABLE IF NOT EXISTS releases (
        id INT PRIMARY KEY,
        title TEXT NOT NULL,
        year INT,
        catalog_no TEXT,
        thumb TEXT,
        resource_url TEXT
    );

    CREATE TABLE IF NOT EXISTS releases_artists (
        release_id INT NOT NULL,
        artist_id INT NOT NULL,
        PRIMARY KEY (release_id, artist_id),
        FOREIGN KEY (release_id) REFERENCES releases (id) ON DELETE CASCADE,
        FOREIGN KEY (artist_id) REFERENCES artists (id) ON DELETE CASCADE
    );

    CREATE TABLE IF NOT EXISTS releases_formats (
        release_id INT NOT NULL,
        format_id INT NOT NULL,
        PRIMARY KEY (release_id, format_id),
        FOREIGN KEY (release_id) REFERENCES releases (id) ON DELETE CASCADE,
        FOREIGN KEY (format_id) REFERENCES formats (id) ON DELETE CASCADE
    );
EOSQL