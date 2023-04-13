import typer
from pureml.cli.puremlconfig import PureMLConfigYML
from pureml.config.parser import Config
# from pureml.trainer.train import Trainer
import os

app = typer.Typer()


import typer
from rich import print
# from puretrainer.train import Trainer
import os
from dotenv import load_dotenv

load_dotenv()


import pureml.cli.auth as auth
import pureml.cli.secrets as secrets
import pureml.cli.orgs as orgs

app = typer.Typer()
app.add_typer(auth.app, name="auth")
app.add_typer(secrets.app, name="secrets")
app.add_typer(orgs.app, name="orgs")


@app.callback(no_args_is_help=True)
def validate_user_authentication(ctx: typer.Context):
    # print(ctx.invoked_subcommand)
    if ctx.invoked_subcommand in ['auth']:
        return
    # user_token = auth.auth_validate()
    return


# init the config file
@app.command()
def init(
    silent: bool = typer.Option(
        False, "--silent", "-s", help="Run in silent mode"
    ),
    # repository: str = typer.Option(
    #     None, "--repository", "-r", help="The repository to use"
    # ),
    backend_url: str = typer.Option(
        None, "--backend-url", "-b", help="The backend url to use"
    ),
    frontend_url: str = typer.Option(
        None, "--frontend-url", "-f", help="The frontend url to use"
    ),
):
    """
    This command will initialize the puremlconfig.yml file for your project
    """
    from pathlib import Path

    project_path = Path.cwd()

    puremlconfig = PureMLConfigYML(project_path / "puremlconfig.yml")
    
    # check if the config file exists
    if puremlconfig.file.exists():
        print("A puremlconfig.yml file already exists in this directory - aborting")
        return
    
    # if not silent, ask for the config values
    if not silent:
        print("Let's get started with the configuration of your project")

        # if not repository:
        #     repository = typer.prompt("Enter your preferred storage repository path [eg. file://pureml_data/]")
        # else:
        #     print(f"Using given repository: {repository}")

        isSelfHosted = typer.confirm("Are you using a self-hosted PureML instance?")

        if isSelfHosted:
            if not backend_url:
                backend_url = typer.prompt("Enter the backend url", default="http://localhost:8080/api/")
            else:
                print(f"Using given backend url: {backend_url}")
            
            if not frontend_url:
                frontend_url = typer.prompt("Enter the frontend url", default="http://localhost:3000/")
            else:
                print(f"Using given frontend url: {frontend_url}")
    else:
        # if not repository:
        #     print("Missing required arguments")
        #     return
        # print(f"Using given repository: {repository}")
        if backend_url:
            print(f"Using given backend url: {backend_url}")
        if frontend_url:
            print(f"Using given frontend url: {frontend_url}")
    
    # create the config file
    data = {
        # "repository": repository,
    }
    if backend_url:
        data["backend_url"] = backend_url
    if frontend_url:
        data["frontend_url"] = frontend_url
    puremlconfig.save(data)