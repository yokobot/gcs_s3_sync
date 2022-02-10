# coding: utf-8
import json
import os
import boto3
from google.cloud import storage
from google.oauth2 import service_account

REGION = os.environ.get('REGION')
PARAM_KEY = os.environ.get('PARAM_KEY')
SCOPES = ['https://www.googleapis.com/auth/devstorage.read_write']


def get_parameters(param_key):
    ssm = boto3.client('ssm', region_name=REGION)
    response = ssm.get_parameters(
        Names=[
            param_key,
        ],
        WithDecryption=True
    )
    return response['Parameters'][0]['Value']


def lambda_handler(event, context):
    param_value = get_parameters(PARAM_KEY)

    service_account_info = json.loads(param_value)
    credentials = service_account.Credentials.from_service_account_info(
        service_account_info,
        scopes=SCOPES
    )

    client = storage.Client(credentials=credentials)
    #TODO ここからgcsにファイルを置く処理と削除する処理を書く