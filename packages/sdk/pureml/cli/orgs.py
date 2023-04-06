import requests
import typer
from rich import print
from rich.console import Console
from rich.table import Table
from urllib.parse import urljoin
from pureml.schema import BackendSchema, PathSchema


path_schema = PathSchema().get_instance()
backend_schema = BackendSchema().get_instance()
app = typer.Typer()

def printecho():
    return 'echo'

def get_org_table(access_token):
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
        count = 0
        table = Table("Sr.No.","User Handle", "Name", "Description", "Role", "Organization Id")
        for org in org_all:
            count += 1
            table.add_row(str(count), org["org"]["handle"], org["org"]["name"], org["org"]["description"], org["role"], org["org"]["uuid"])

        console.print(table)
        print()
        sr_no = -1
        while int(sr_no) not in range(1, count + 1):
            sr_no: str = typer.prompt("Enter your Sr.No. of Organization (1 .... " + str(count) + ")")
            if int(sr_no) not in range(1, count + 1):
                print("[bold red]Invalid Sr.No. of Organization!")
                print("Try Again!")
                print()
        org_id = org_all[int(sr_no) - 1]["org"]["uuid"]
        url_path = "org/id/{}".format(org_id)
        url = urljoin(backend_schema.BASE_URL, url_path)

        headers = {
            "accept": "application/json",
            "Authorization": "Bearer {}".format(access_token),
        }

        response = requests.get(url, headers=headers)

        if response.ok:
            print("[bold green]Organization Selected!")
            return org_id
        else:
            print("[bold red]Organization doesn't Exists!")
            return None
    else:
        print("[bold red]Invalid Credentials!")
        return None