# Music Streaming Endpoint Documentation

This document describes the REST endpoint for streaming music files in the Webmane Go Music API. The endpoint serves audio files stored on the server, supports HTTP range requests for efficient streaming, and integrates with the PostgreSQL database via the Ent ORM to locate music files.

---

## Endpoint Overview

- **Endpoint URL:** `/music`
- **HTTP Method:** `GET`
- **Query Parameters:**
  - `id` (required): The unique integer ID of the music record in the database.

---

## Functionality

The `/music` endpoint streams audio files to clients based on the requested music ID. It performs the following steps:

1. **ID Validation:**

   - Extracts the `id` query parameter from the request URL.
   - Validates that the `id` is present and is a valid integer.
   - Returns HTTP 400 Bad Request if validation fails.

2. **Database Lookup:**

   - Uses the Ent client to query the `music_ent` table for the record with the specified ID.
   - Retrieves the file path (`path` field) of the music file.
   - Returns HTTP 500 Internal Server Error if the database query fails.

3. **File Handling:**

   - Opens the audio file at the retrieved path.
   - Returns HTTP 404 Not Found if the file does not exist.
   - Returns HTTP 500 Internal Server Error for other file access errors.

4. **Content-Type Determination:**

   - Determines the MIME type based on the file extension (`.mp3`, `.flac`, `.m4a`, `.mp4`).
   - Sets the appropriate `Content-Type` header for the response.
   - Defaults to `application/octet-stream` for unknown extensions, logging a warning.

5. **Range Request Support:**

   - Supports HTTP `Range` headers to allow partial content delivery.
   - Parses the `Range` header to determine the byte range requested.
   - Validates the range and returns HTTP 416 Range Not Satisfiable if invalid.
   - Sets `Content-Range` and `Content-Length` headers accordingly.
   - Returns HTTP 206 Partial Content status for range responses.
   - If no `Range` header is present, serves the entire file with HTTP 200 OK.

6. **Streaming:**
   - Streams the requested byte range (or entire file) to the client efficiently.
   - Handles errors during streaming by logging them.

---

## Supported Audio Formats

The endpoint supports streaming of the following audio file formats, identified by their extensions:

- `.mp3` → `audio/mp3`
- `.flac` → `audio/flac`
- `.m4a` → `audio/m4a`
- `.mp4` → `audio/mp4`

Files with unsupported extensions are served with a generic `application/octet-stream` content type.

---

## Error Handling

- **400 Bad Request:** Missing or invalid `id` parameter, or malformed `Range` header.
- **404 Not Found:** Music file not found on the server.
- **416 Range Not Satisfiable:** Requested byte range is outside the file size.
- **500 Internal Server Error:** Database errors, file access errors, or other unexpected failures.

---

## Integration with Database Schema

- The endpoint relies on the `music_ent` table managed by the Ent ORM.
- Each music record has a unique integer ID and a `path` field storing the file system path to the audio file.
- The endpoint queries the database to retrieve the file path for the requested music ID.
- This design decouples the API from the file storage details, allowing flexible file management.

---

## Example Request

```
GET /music?id=123 HTTP/1.1
Host: example.com
Range: bytes=0-1023
```

- Streams the first 1024 bytes of the music file with ID 123.

---

## Summary

The `/music` REST endpoint provides efficient, range-supported streaming of audio files stored on the server. It integrates tightly with the PostgreSQL database via Ent to locate files by ID, ensuring a robust and scalable music streaming service within the Webmane Go Music API.
