import os
from pathlib import Path
from typing import Optional
import jwt
import requests
import typer
from rich import print
from rich.syntax import Syntax
from .orgs import get_org_table
from pureml.components import get_org_id, get_token
from rich.console import Console
from rich.table import Table
from urllib.parse import urljoin
import json
from pureml.schema import BackendSchema, PathSchema


path_schema = PathSchema().get_instance()
backend_schema = BackendSchema().get_instance()
app = typer.Typer()


def save_auth(org_id: str = None, access_token: str = None, email=None):
    token_path = path_schema.PATH_USER_TOKEN

    token_dir = os.path.dirname(token_path)
    os.makedirs(token_dir, exist_ok=True)

    token = {"org_id": org_id, "accessToken": access_token, "email": email}

    token = json.dumps(token)

    with open(token_path, "w") as token_file:
        token_file.write(token)

    print("[green]User token is stored")


@app.command()
def details():
    token = get_token()
    org_id = get_org_id()

    print("Org Id: ", org_id)
    print("Access Token: ", token)


@app.callback()
def callback():
    """
    Authentication user command

    Use with status, signup, login or logout option

    status - Checks current authenticated user
    signup - Creates new user
    login - Logs in user
    logout - Logs out user
    """


@app.command()
def signup():
    print("\n[bold]Create a new account[/bold]\n")
    email: str = typer.prompt("Enter new email")
    handle: str = typer.prompt("Enter new user handle")
    name: str = typer.prompt("Enter new user name")
    password: str = typer.prompt(
        "Enter new password", confirmation_prompt=True, hide_input=True
    )
    # organization_id: str = typer.prompt("Enter Organization id")
    # data = {"email": email, "password": password, "org": organization_id}
    data = {"email": email, "password": password, "handle": handle, "name": name}

    url_path_1 = "user/signup"
    url = urljoin(backend_schema.BASE_URL, url_path_1)

    response = requests.post(url, json=data)

    if not response.ok:
        print(f"[bold red]Could not create account! Please try again later")
        return
    print(
        f"[bold green]Successfully created your account! You can now login using ```pure auth login```"
    )


def list_org(access_token):

    url_path = "org"
    url = urljoin(backend_schema.BASE_URL, url_path)

    headers = {
        "accept": "application/json",
        "Authorization": "Bearer {}".format(access_token),
    }

    response = requests.get(url, headers=headers)

    if response.ok:
        print()
        print("[bold green]Select the Organization from the list below!")
        org_all = response.json()["data"]

        console = Console()

        table = Table("User Handle", "Organization Id")
        for org in org_all:
            table.add_row(org["org"]["handle"], org["org"]["uuid"])

        console.print(table)
        print()

    else:
        print("[bold red]Unable to fetch existing Organizations!")


def check_org_status(access_token):

    org_id: str = typer.prompt("Enter your Org Id")

    url_path = "org/id/{}".format(org_id)
    url = urljoin(backend_schema.BASE_URL, url_path)

    headers = {
        "accept": "application/json",
        "Authorization": "Bearer {}".format(access_token),
    }

    response = requests.get(url, headers=headers)

    if response.ok:
        print("[bold green]Organization Exists!")
        return org_id
    else:
        print("[bold red]Organization doesnot Exists!")
        return None


@app.command()
def login():

    print(f"\n[bold]Enter your credentials to login[/bold]\n")
    email: str = typer.prompt("Enter your email")
    password: str = typer.prompt("Enter your password", hide_input=True)
    data = {"email": email, "password": password}

    url_path = "user/login"
    url = urljoin(backend_schema.BASE_URL, url_path)

    response = requests.post(url, json=data)

    if response.ok:
        token = response.text
        token = json.loads(token)["data"][0]

        access_token = token["accessToken"]
        email = token["email"]

        
        org_id = get_org_table(access_token)
        # list_org(access_token=access_token)

        # org_id = check_org_status(access_token=access_token)

        if org_id is not None:

            save_auth(org_id=org_id, access_token=access_token, email=email)
            print(f"[bold green]Successfully logged in to your account!")

            print("org_id:", org_id)
            print("accessToken:", access_token)
        else:
            print(f"[bold red]Unable to login to your account!")

    else:
        print(f"[bold red]Unable to login to your account!")


# @app.command()
# def status():
#     print()
#     path = PATH_USER_TOKEN

#     curr_path = Path(__file__).parent.resolve()
#     if os.path.exists(path):
#         token = open(path, "r").read()
#         public_key = open(f"{curr_path}/public.pem", "rb").read()
#         payload = jwt.decode(token, public_key, algorithms=["RS256"])
#         print(f"[bold green]You are currently logged in as {payload['email']}")
#     else:
#         print("[bold red]You are not logged in!")


def statusHelper():
    path = path_schema.PATH_USER_TOKEN

    if os.path.exists(path):
        return open(path, "r").read()
    else:
        return None


def auth_validate():
    token = statusHelper()
    if not token:
        print("[bold red]You are not logged in!")
        raise typer.Exit()
    return token


# @app.command()
# def logout():
#     print()
#     path = PATH_USER_TOKEN

#     if os.path.exists(path):
#         os.remove(path)
#         print(f"[bold yellow]Successfully logged out!")
#     else:
#         print(f"[bold red]You are not logged in!")

if __name__ == "__main__":
    app()
