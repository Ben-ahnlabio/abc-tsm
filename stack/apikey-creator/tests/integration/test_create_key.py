import os
import json

import pytest
import dotenv

from key_creator import app


@pytest.fixture()
def env():
    """Generates API GW Event"""
    dotenv.load_dotenv()


@pytest.fixture()
def cloudformation_event():
    return {}


def test_create_key(cloudformation_event, env):
    app.lambda_handler(cloudformation_event, None)
