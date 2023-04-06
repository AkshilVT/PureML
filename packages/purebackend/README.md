[![PureML](/assets/BackendCoverImg.png)](https://pureml.com)

<br/>
<br/>

<div align="center">
  <a
    href="https://github.com/golang/go"
  >
    <img alt="Remix" src="https://img.shields.io/badge/golang-1.19-blue?style=flat&logo=go&logoColor=white" />
  </a>
  &nbsp;
  <a
    href="https://opensource.org/licenses/Apache-2.0"
  >
    <img alt="License" src="https://img.shields.io/badge/License-Apache2.0-green?style=flat&logo=apache&logoColor=white" />
  </a>
  &nbsp;
  <a
    href="https://discord.gg/xNUHt9yguJ"
  >
    <img alt="Discord" src="https://img.shields.io/badge/Discord-Join%20Discord-blueviolet?style=flat&logo=discord&logoColor=white" />
  </a>
  &nbsp;
</div>
<br/>

## Quick start

PureML Backend enables you to setup your own API endpoint server to manage all the cloud data pushed by the client SDKs. It takes no time to run on your local system. Follow below steps to run:

Install dependencies:
```bash
go get .
```

Start the server:
> _Assuming system supports make command. Check out [how to install make](https://www.gnu.org/software/make/#download)_
```bash
make run
```

After making any API changes to automatically build swagger documentation run:
```bash
make docs run
```

Server will run at http://localhost:8080/api/ by default
Open [/api/swagger/index.html](http://localhost:8080/api/swagger/index.html) to see the Open API documentation.

<br/>

## Directory Structure

```
purebackend (will be soon renamed to pureml_backend)
├── auth                                  # Authentication module
│   ├── api                               # Auth API files
│   │   ├── service                       # Auth APIs
│   │   │   ├── admin.go
│   │   │   ├── user.go
│   │   │   └── ...
│   │   └── tests                         # Auth API tests
│   │       ├── tests.go
│   │       └── user_test.go
│   └── middlewares
│       ├── jwt_authenticate.go           # Jwt Authentication Middleware
│       └── require_auth_context.go       # Middleware to Require Auth in any API
│
├── core                                  # Core module
│   ├── app.go                            # Public interface for pureml backend App. App instance is defined in this
│   ├── base.go                           # Defines all the functionality of app interface
│   │
│   ├── apis                              # Core API setup files
│   │   ├── app.go                        # Defines serve function
│   │   ├── base.go                       # Defines InitAPI function that binds all modules and serves them
│   │   ├── service                       # Core API util functionalities
│   │   │   ├── api.go                    # API interface                    ──┐
│   │   │   ├── handler.go                # Common Request response handler  ──┴─ (Need to be present in all modules to prevent cyclic imports)
│   │   │   ├── health_check.go           # Health API
│   │   │   ├── helpers.go
│   │   │   └── utils.go
│   │   └── tests
│   │       └── health_test.go            # Health API test
│   │
│   ├── common                            # Requirements that are common for models and datasets
│   │   ├── dbmodels                      # Common database models
│   │   │   ├── base.go
│   │   │   └── dbmodels.go
│   │   └── models                        # Common schema models
│   │       ├── models.go
│   │       └── sources.go
│   │
│   ├── config                            # Config interface for pureml backend instance manual configuration
│   │   └── config.go
│   │
│   ├── daos                              # Data access objects
│   │   ├── daos.go                       # Public interface for DAOs (All public database operations go here)
│   │   ├── datastore                     # Datastore object for actual database operations
│   │   │   └── datastore.go
│   │   ├── seeder                        # Seeder package
│   │   │   ├── seeder.exe                # Seeder build to seed local or prod db
│   │   │   └── seeder.go
│   │   └── seeds                         # Seed interface for defining seeds
│   │       ├── seed.go
│   │       └── seeds.go
│   │
│   ├── dbmodels                          # Core Database models (Base Model)
│   │   └── dbmodels.go
│   │
│   ├── middlewares                       # Core middleware utilities
│   │   └── utils.go
│   │
│   ├── models                            # Core request response schema models
│   │   ├── request.go
│   │   └── response.go
│   │
│   ├── settings                          # Configuration regarding all the secrets and environment
│   │   └── settings.go
│   │
│   └── tools                             # Tools used in application functionality
│       ├── filesystem                    # Managing filesystem (upload storage & retrieve data)
│       │   └── ...
│       ├── inflector                     # String manipulation
│       │   └── ...
│       ├── list                          # List manipulation
│       │   └── ...
│       ├── mailer                        # Mailing service
│       │   └── ...
│       ├── search                        # Search integration
│       │   └── ...
│       ├── security                      # Security functions
│       │   └── ...
│       └── types                         # Type handling
│           └── ...
│
├── dataset                               # Dataset module
│   ├── api                               # Dataset API files
│   │   ├── service                       # Dataset APIs
│   │   │   └── ...
│   │   └── tests                         # Dataset API tests
│   │       └── ...
│   ├── dbmodels                          # Dataset database models
│   │   └── dbmodels.go
│   ├── middlewares                       # Dataset validation middlewares
│   │   └── ...
│   └── models                            # Dataset schema models
│       └── models.go
│
├── docs                                  # Auto generated Open API Swagger Docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── examples                              # Example usage of Pureml backend package
│   ├── base                              # Base/Default coonfiguration
│   │   └── main.go
│   └── prod                              # Production configuration
│       └── main.go
│
├── model                                 # Model module
│   ├── api                               # Model API files
│   │   ├── service                       # Model APIs
│   │   │   └── ...
│   │   └── tests                         # Model API tests
│   │       └── ...
│   ├── dbmodels                          # Model database models
│   │   └── dbmodels.go
│   ├── middlewares                       # Model validation middlewares
│   │   └── ...
│   └── models                            # Model schema models
│       └── models.go
│
├── scripts                               # Helper scripts package
│   ├── run.sh
│   ├── scripts.exe
│   └── scripts.go
│
├── test                                  # tests
│   ├── api.go                            # TestScenario definition
│   ├── app.go                            # Public test interface
│   ├── base.go                           # Test execution logic
│   ├── data
│   │   └── pureml.db                     # Default test seed database (sqlite3)
│   └── request.go                        # Test API request helpers
│
├── user_org                              # User and Organization module
│   ├── api                               # User and Organization API files
│   │   ├── service                       # User and Organization APIs
│   │   │   └── ...
│   │   └── tests                         # User and Organization API tests
│   │       └── ...
│   ├── dbmodels                          # User and Organization database models
│   │   ├── org.go
│   │   └── user.go
│   ├── middlewares                       # User and Organization validation middlewares
│   │   └── validate_org.go
│   └── models                            # User and Organization schema models
│       ├── org.go
│       └── user.go
│
├── purebackend.go                        # Pureml backend package public interface. All package methods defined in this
├── go.mod                                # Go modules
├── go.sum                                # Go modules checksum
├── golangci.yml                          # Go formating and linting config
├── .dockerignore                         # List of files and folders not tracked by Docker
├─ .gitignore                             # List of files and folders not tracked by Git
├── Dockerfile                            # Dockerfile for production docker image configuration
├── Dockerfile.base                       # Dockerfile for official base docker image configuration
├── Makefile                              # Makefile for script automation
└── README.md                             # This file
```

## Technology used

1. [Go Lang](https://go.dev/)
2. [Echo Framework](https://echo.labstack.com/)
3. [Swaggo](https://github.com/swaggo/swag)

## Reporting Bugs

To report any bugs you have faced while using PureML package, please

1. Report it in [Discord](https://discord.gg/xNUHt9yguJ) channel
2. Open an [issue](https://github.com/PureMLHQ/PureML/issues)

<br />

## Contributing and Developing

Lets work together to improve the features for everyone. Here's step one for you to go through our [Contributing Guide](./CONTRIBUTING.md). We are already waiting for amazing ideas and features which you all have got.

Work with mutual respect. Please take a look at our public [Roadmap here](https://pureml.notion.site/7de13568835a4cf18913307503a2cdd4?v=82199f96833a48e5907023c8a8d565c6).

<br />

## Community

To get quick updates of feature releases of PureML, follow us on:

[<img alt="Twitter" height="20" src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" />](https://twitter.com/getPureML) [<img alt="LinkedIn" height="20" src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" />](https://www.linkedin.com/company/PuremlHQ/) [<img alt="GitHub" height="20" src="https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white" />](https://github.com/PureMLHQ/PureML) [<img alt="GitHub" height="20" src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" />](https://discord.gg/DBvedzGu)

<br/>

## 📄 License

See the [Apache-2.0](./License) file for licensing information.
