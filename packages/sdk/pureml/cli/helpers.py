import json
import os
from pureml.schema import PathSchema

path_schema = PathSchema().get_instance()

def save_auth(org_id: str = None, access_token: str = None, email: str = None):
    token_path = path_schema.PATH_USER_TOKEN

    token_dir = os.path.dirname(token_path)
    os.makedirs(token_dir, exist_ok=True)

    # Read existing token
    if os.path.exists(token_path):
        with open(token_path, "r") as token_file:
            token = json.load(token_file)

        if org_id is not None:
            token["org_id"] = org_id
        if access_token is not None:
            token["accessToken"] = access_token
        if email is not None:
            if "email" in token and token["email"] != email:
                token["org_id"] = ""
            token["email"] = email
    else:
        token = {"org_id": org_id, "accessToken": access_token, "email": email}
        if org_id is None:
            token["org_id"] = ""

    token = json.dumps(token)

    with open(token_path, "w") as token_file:
        token_file.write(token)
