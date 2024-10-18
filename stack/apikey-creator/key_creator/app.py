import os
import logging
import base64
import uuid

import boto3
import cfnresponse
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)


def lambda_handler(event, context):
    LOGGER.info("Received api key generation event: %s", event)

    request_type = event.get("RequestType")
    LOGGER.info("Request Type: %s", request_type)
    if request_type is None:
        key_create()
    elif request_type == "Create":
        key_create()
        cfnresponse.send(
            event, context, "SUCCESS", {"Message": "Resource create successful!"}
        )
    elif request_type == "Update":
        cfnresponse.send(
            event, context, "SUCCESS", {"Message": "Resource update successful!"}
        )
    elif request_type == "Delete":
        cfnresponse.send(
            event, context, "SUCCESS", {"Message": "Resource deletion successful!"}
        )


def key_create():
    secret_client = boto3.client("secretsmanager")

    api_key_secret_arn = os.getenv("NODE_API_KEY_SECRET_ARN")
    private_key_secret_arn = os.getenv("NODE_PRIVATE_KEY_SECRET_ARN")
    public_key_secret_arn = os.getenv("NODE_PUBLIC_KEY_SECRET_ARN")

    secret_client = boto3.client("secretsmanager")
    private_key, public_key = create_key()
    api_key = str(uuid.uuid4())

    secret_client.put_secret_value(
        SecretId=api_key_secret_arn,
        SecretString=api_key,
    )

    secret_client.put_secret_value(
        SecretId=private_key_secret_arn,
        SecretString=private_key,
    )

    secret_client.put_secret_value(
        SecretId=public_key_secret_arn,
        SecretString=public_key,
    )


def create_key():
    # P-256 타원 곡선 개인 키 생성
    private_key = ec.generate_private_key(ec.SECP256R1())

    # 개인 키를 DER 형식으로 인코딩
    der_private_key = private_key.private_bytes(
        encoding=serialization.Encoding.DER,
        format=serialization.PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.NoEncryption(),
    )

    # DER 형식의 개인 키를 Base64로 인코딩
    base64_private_key = base64.b64encode(der_private_key).decode("utf-8")

    # 공개 키 추출
    public_key = private_key.public_key()

    # 공개 키를 DER 형식으로 인코딩
    der_public_key = public_key.public_bytes(
        encoding=serialization.Encoding.DER,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )

    # DER 형식의 공개 키를 Base64로 인코딩
    base64_public_key = base64.b64encode(der_public_key).decode("utf-8")
    return base64_private_key, base64_public_key
