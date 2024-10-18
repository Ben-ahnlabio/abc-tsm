import argparse
import base64
import hashlib
import json
import logging
import os
import pathlib
import sys

import boto3
from botocore.config import Config

log = logging.getLogger(__name__)
log.setLevel(logging.INFO)
log.addHandler(logging.StreamHandler())

current_path = pathlib.Path(__file__).parent


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--this", help="parameter arn for this node", required=True)
    parser.add_argument("--other", help="parameter arn for other node", required=True)
    parser.add_argument("--db", help="secrets arn for db", required=True)
    args = parser.parse_args()

    log.info("this: %s", args.this)
    log.info("other: %s", args.other)

    log.info("getting env from this param...")
    this_env = get_env_from_this_param(args.this)
    log.info("getting env from other param...")
    other_env = get_env_from_other_param(args.other)
    other_player_address = get_other_player_address()
    env_dict = {
        **this_env,
        **other_env,
        "OTHER_PLAYER_ADDRESS": other_player_address,
        "THIS_PLAYER_ADDRESS": os.getenv("THIS_PLAYER_ADDRESS"),
        "PLAYER_INDEX": os.getenv("NODE_INDEX"),
        "APPLICATION_ID": os.getenv("APPLICATION_ID"),
        "DATABASE_SOURCE_NAME": get_db_string_from_secret_arn(args.db),
    }
    env_dict["NODE_API_KEY"] = api_key_hash(env_dict["NODE_API_KEY"])

    log.info("writing env file...")
    dict_to_env_file(current_path / ".env.production", env_dict)

    return 0


def get_other_player_address():
    """
    node-1.${HostedZoneName}
    node-1.dev-p2.tsm.myabcwallet.com
    """
    
    hosted_zone_name = os.getenv("OTHER_HOSTED_ZONE_NAME")
    return f"node-1.{hosted_zone_name}:9000?connectionPoolSize=6&connectionLifetime=5m"


def get_env_from_this_param(this_param_arn: str) -> dict:
    this_env_keys = {
        "NODE_PRIVATE_KEY_SECRET_ARN": "PLAYER_PRIVATE_KEY",
        "NODE_MASTER_ENCRYPTION_KEY_SECRET_ARN": "MASTER_ENCRYPTION_KEY",
        "NODE_API_KEY_SECRET_ARN": "NODE_API_KEY",
    }
    return get_env_from_param(this_env_keys, this_param_arn)


def get_env_from_other_param(other_param_arn: str) -> dict:
    other_env_keys = {
        "NODE_PUBLIC_KEY_SECERT_ARN": "OTHER_PLAYER_PUBLIC_KEY",
    }
    return get_env_from_param(other_env_keys, other_param_arn)

def get_env_from_param(env_keys: dict, param_arn: str) -> dict:
    param = get_params_from_ssm(param_arn)
    param_value = json.loads(param["Parameter"]["Value"])
    return {
        env_keys[k]: get_secret_value_from_secret_arn(v)
        for k, v in param_value.items()
        if k in env_keys
    }


def dict_to_env_file(env_file_path: pathlib, env_dict: dict):
    with open(env_file_path, "w") as f:
        for k, v in env_dict.items():
            f.write(f"{k}={v}\n")


def get_params_from_ssm(parameter_arn: str) -> dict:
    config = Config(
        region_name=get_region_from_arn(parameter_arn),
    )
    ssm_client = boto3.client("ssm", config=config)
    response = ssm_client.get_parameter(Name=parameter_arn)
    return response


def get_secret_value_from_secret_arn(secret_arn: str) -> str:
    """
    arn:aws:ssm:ap-northeast-2:915486611144:parameter/abcdev/tsm/node-bootstrap-parameter-1
    """
    config = Config(
        region_name=get_region_from_arn(secret_arn),
    )
    secret_client = boto3.client("secretsmanager", config=config)
    response = secret_client.get_secret_value(SecretId=secret_arn)
    return response["SecretString"]

def get_db_string_from_secret_arn(secret_arn: str) -> str:
    secret_string = get_secret_value_from_secret_arn(secret_arn)
    secert_dict = json.loads(secret_string)
    return f"postgres://{secert_dict['username']}:{secert_dict['password']}@{secert_dict['host']}:{secert_dict['port']}/{secert_dict['dbname']}?sslmode=disable"

def get_region_from_arn(arn: str):
    return arn.split(":")[3]


def api_key_hash(api_key: str) -> str:
    # SHA-256 해시 계산
    sha256_hash = hashlib.sha256(api_key.encode()).digest()
    # Base64 인코딩
    return base64.b64encode(sha256_hash).decode()


if __name__ == "__main__":
    exit_code = main()
    sys.exit(exit_code)
