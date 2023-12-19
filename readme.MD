# video-encoder-golang

## *video* service
HTTP-based service responsible for encoding .mp4 videos.

Depends on:
- Remote Storage (e.g. GCP Cloud Bucket).

Responsibilities:
- Download .mp4 files from a Cloud Bucket in GCP.
- Convert .mp4 files into fragments using mp4fragment.
- Generate video slices using Bento4.

## *upload* service
Background process responsible for uploading recently encoded videos to a remote storage.

Depends on:
- VideoService
- Remote Storage (e.g. GCP Cloud Bucket).

Responsibilities:
- Scan for recently encoded videos.
- Upload video slices in batches to Remote Storage.