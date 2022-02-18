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
