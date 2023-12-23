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


## How to test the app?

1. Start all applications in background with `docker compose up -d`
1. Navigate to your local instance of RabbitMQ Manager at localhost:15
2. Create a new exchange called "dlx" type "fanout" and durability "durable".
3. Create a new queue named "processed-videos" and bind it to "amq.direct" with routing key "jobs"
4. Create a new queue named "rejected-videos" and bind it to "dlx" with no specific routing key.
5. Publish a new message to "videos" with the following format:

```json
{
  "resource_id": "my-resource-id-can-be-a-uuid-type",
  "file_path": "convite.mp4"
}
```
6. Follow the logs to see the application running.
7. The queue `processed-videos` should have a new message by now. Otherwise, something should be forwarded to `rejected-videos` in case of failure.