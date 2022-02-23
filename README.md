# gcs_s3_sync

Google Cloud Platform のストレージ (GCS) と AWS のストレージ (S3) を同期するアプリケーションです。

### 仕様

- 同期する2つのストレージの名前は同一とする
- 同期方式
    - s3
        - 作成、変更
            - 同名ファイルがgcsに存在しているか確認して、存在しなければファイルをgcsにコピーする
        - 削除
            - gcsに同名ファイルが存在しているか確認して、存在していればファイルを削除する
    - gcs
        - 作成、変更
            - 同名ファイルがs3に存在しているか確認して、存在していれば何もしない、存在しなければファイルをs3にコピーする
        - 削除
            - s3に同名ファイルが存在しているか確認して、存在していればファイルを削除する

### テスト用データ

- gcp
    - ローカルテスト用に `functions-framework-go` を使用しています

    ```
    [テスト用サーバ起動]

    cd terraform/gcp/src
    go run cmd/main.go
    ```

    ```
    [テスト実行]

    curl localhost:8080 \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{
          "context": {
            "eventId": "1147091835525187",
            "timestamp": "2020-04-23T07:38:57.772Z",
            "eventType": "google.storage.object.finalize",
            "resource": {
               "service": "storage.googleapis.com",
               "name": "projects/_/buckets/MY_BUCKET/MY_FILE.txt",
               "type": "storage#object"
            }
          },
          "data": {
            "bucket": "MY_BUCKET",
            "contentType": "text/plain",
            "kind": "storage#object",
            "md5Hash": "...",
            "metageneration": "1",
            "name": "MY_FILE.txt",
            "size": "352",
            "storageClass": "MULTI_REGIONAL",
            "timeCreated": "2020-04-23T07:38:57.230Z",
            "timeStorageClassUpdated": "2020-04-23T07:38:57.230Z",
            "updated": "2020-04-23T07:38:57.230Z"
          }
        }'
    ```
