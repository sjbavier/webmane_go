# GraphQL API Documentation

This document describes the GraphQL API endpoints exposed by the Webmane Go Music API, their purpose, and how they relate to the underlying PostgreSQL database schema managed via the Ent framework.

---

## Overview

The GraphQL API provides queries and mutations to interact with two main entities:

- **Song** (represented in the database by the `music_ent` table)
- **Playlist** (represented in the database by the `playlists` table)

The API supports fetching, creating, updating, and managing songs and playlists, including the relationships between them.

---

## GraphQL Queries

### `music`

- **Description:** Fetches a paginated list of songs with optional search filtering.
- **Arguments:**
  - `pageNumber` (Int, optional): The page number for pagination.
  - `pageSize` (Int, optional): The number of items per page.
  - `searchText` (String, optional): Text to filter songs by title, artist, album, or genre.
- **Returns:** `MusicResponse` object containing:
  - `songs`: List of `Song` objects.
  - `totalItemsCount`: Total number of songs matching the criteria.

### `song`

- **Description:** Fetch a single song by its unique ID.
- **Arguments:**
  - `id` (ID!): The unique identifier of the song.
- **Returns:** A `Song` object.
- **Note:** This query is defined in the schema but currently not implemented.

### `playlists`

- **Description:** Fetches a paginated list of playlists with optional search filtering.
- **Arguments:**
  - `pageNumber` (Int, optional): The page number for pagination.
  - `pageSize` (Int, optional): The number of items per page.
  - `searchText` (String, optional): Text to filter playlists by name.
- **Returns:** `PlaylistResponse` object containing:
  - `playlists`: List of `Playlist` objects.
  - `totalItemsCount`: Total number of playlists matching the criteria.

### `playlist`

- **Description:** Fetch a single playlist by its unique ID.
- **Arguments:**
  - `id` (ID!): The unique identifier of the playlist.
- **Returns:** A `Playlist` object.

---

## GraphQL Mutations

### `upsertSong`

- **Description:** Creates a new song or updates an existing one.
- **Arguments:**
  - `input` (SongInput!): The song data to create or update.
- **Returns:** The created or updated `Song` object.

### `additivePathUpsertSong`

- **Description:** Similar to `upsertSong`, but may have additive behavior (implementation-specific).
- **Arguments:**
  - `input` (SongInput!): The song data.
- **Returns:** The created or updated `Song` object.

### Playlist Mutations

- **`upsertPlaylist`**: Create or update a playlist.
- **`deletePlaylist`**: Delete a playlist by ID.
- **`addSongToPlaylist`**: Add a song to a playlist by their IDs.
- **`removeSongFromPlaylist`**: Remove a song from a playlist by their IDs.

---

## GraphQL Types and Database Schema Mapping

### Song

- **GraphQL Type:** `Song`
- **Database Table:** `music_ent`
- **Fields:**

| GraphQL Field  | Database Field | Description                        |
| -------------- | -------------- | ---------------------------------- |
| `id`           | `id` (auto)    | Unique identifier (auto-generated) |
| `path`         | `path`         | File path of the song (unique)     |
| `lastUpdate`   | `last_update`  | Timestamp of last update           |
| `title`        | `title`        | Song title                         |
| `artist`       | `artist`       | Artist name                        |
| `album`        | `album`        | Album name                         |
| `genre`        | `genre`        | Genre of the song                  |
| `release_year` | `release_year` | Release year (stored as string)    |
| `cover_art`    | `cover_art`    | Cover art image path or URL        |

- **Notes:**
  - The `path` field is unique and not empty.
  - The `last_update` field is automatically set to the current time on creation and update.
  - Songs have a many-to-many relationship with playlists.

### Playlist

- **GraphQL Type:** `Playlist`
- **Database Table:** `playlists`
- **Fields:**

| GraphQL Field  | Database Field          | Description                                  |
| -------------- | ----------------------- | -------------------------------------------- |
| `id`           | `id` (auto)             | Unique identifier (auto-generated)           |
| `name`         | `name`                  | Playlist name                                |
| `lastUpdate`   | `last_update`           | Timestamp of last update                     |
| `lastAccessed` | `last_accessed`         | Timestamp of last access (optional)          |
| `cover_art`    | `cover_art`             | Cover art image path or URL (optional)       |
| `songs`        | Relation to `music_ent` | List of songs in the playlist (many-to-many) |

- **Notes:**
  - The `name` field is required and not empty.
  - The `last_update` field is automatically set on creation and update.
  - The `last_accessed` and `cover_art` fields are optional.
  - Playlists have a many-to-many relationship with songs.

---

## Relationships

- Songs and Playlists have a many-to-many relationship.
- In the database schema, this is represented by edges in Ent:
  - `Music` schema has an edge from `playlists` referencing `Playlist`.
  - `Playlist` schema has an edge to `songs` referencing `Music`.
- GraphQL mutations `addSongToPlaylist` and `removeSongFromPlaylist` manage this relationship.

---

## Summary

The GraphQL API provides a type-safe, flexible interface to manage music data stored in a PostgreSQL database via the Ent ORM. The schema definitions and resolvers map closely to the database schema, ensuring consistency and ease of maintenance.

- Queries allow fetching songs and playlists with pagination and search.
- Mutations allow creating, updating, and managing songs and playlists, including their relationships.
- The database schema enforces constraints such as unique paths for songs and required fields for playlists.

This design supports efficient music library management with a clear separation of concerns between API, business logic, and data persistence.
