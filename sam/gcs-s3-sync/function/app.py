# coding: utf-8
import json
import os
import boto3
from google.cloud import storage
from google.oauth2 import service_account

REGION = os.environ.get('REGION')
PARAM_KEY = os.environ.get('PARAM_KEY')
SCOPES = ['https://www.googleapis.com/auth/devstorage.read_write']
BUCKET_NAME = 'yokobot-dev'


def get_parameters(param_key):
    ssm = boto3.client('ssm', region_name=REGION)
    response = ssm.get_parameters(
        Names=[
            param_key,
        ],
        WithDecryption=True
    )
    return response['Parameters'][0]['Value']


def get_google_credentials(google_api_key):
    service_account_info = json.loads(google_api_key)
    credentials = service_account.Credentials.from_service_account_info(
        service_account_info,
        scopes=SCOPES
    )
    return credentials


def make_s3_client():
    s3_client = boto3.client('s3')
    return s3_client


def make_gcs_client(credentials):
    gcs_client = storage.Client(credentials=credentials)
    return gcs_client


def copy_object(s3_client, gcs_client, bucket_name, obj_name):
    with open('/tmp/' + str(obj_name), 'wb') as data:
        s3_client.download_fileobj(bucket_name, obj_name, data)

    gcs_bucket = gcs_client.get_bucket(bucket_name)
    blob = gcs_bucket.blob(obj_name)
    blob.upload_from_filename('/tmp/' + str(obj_name))


def lambda_handler(event, context):
    google_api_key = get_parameters(PARAM_KEY)
    credentials = get_google_credentials(google_api_key)
    s3_client = make_s3_client()
    gcs_client = make_gcs_client(credentials)

    if event['Records'][0]['eventName'].startswith('ObjectCreated'):
        copy_object(
            s3_client,
            gcs_client,
            BUCKET_NAME,
            event['Records'][0]['s3']['object']['key']
        )
    elif event['Records'][0]['eventName'].startswith('ObjectRemoved'):
        #ここにgcs削除を書く
        None
    else:
        None
