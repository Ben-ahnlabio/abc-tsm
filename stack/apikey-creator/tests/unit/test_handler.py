from key_creator import app


def test_create_key():
    private_key, public_key = app.create_key()
    assert private_key
    assert public_key
