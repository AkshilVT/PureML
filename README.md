[![PureML](/assets/PureMLCoverImg.png)](https://pureml.com)

<h align="center">

# The next-gen developer platform for Production ML.

</h>

<div align="center">
  <a
    href="https://pypi.org/project/pureml/"
  >
    <img alt="PyPi" src="https://img.shields.io/pypi/v/pureml?color=green&logo=pureml" />
  </a>
  &nbsp;
  <a
    href="https://pypi.org/project/pureml/"
  >
    <img alt="^3.8" src="https://img.shields.io/pypi/pyversions/pureml">
  </a>
  &nbsp;
  <a
    href="https://pypi.org/project/pureml/"
  >
    <img alt="Coverage" src="https://img.shields.io/codecov/c/github/PureMLHQ/PureML">
  </a>
  &nbsp;
  <a
    href="https://opensource.org/licenses/Apache-2.0"
  >
    <img alt="License" src="https://img.shields.io/pypi/l/pureml?color=red&logo=Apache&logoColor=red" />
  </a>
  &nbsp;
  <a
    href="https://pepy.tech/project/pureml"
  >
    <img alt="Downloads" src="https://static.pepy.tech/badge/pureml">
  </a>
</div>

<br/>
<br/>

## ⏳ Status

This is an early alpha. The implementation might change between versions without warning. Please use at your own risk and pin to a specific version if you're relying on this for anything important!

## ⏱ Getting started

### 1. Installation

Manage versioning of datasets and models with our python SDK. Versioning is semantic and managed automatically. You can install and run PureML using `pip`.

Getting started is simple:

```bash
pip install pureml
```

<br />

If you are trying to manage versions of dataset all you have to do is use our decorator `@dataset`.

```python
import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers
from pureml.decorators import dataset

@dataset("petdata:dev")
def load_data(img_folder = "PetImages"):
  image_size = (180, 180)
  batch_size = 16
  train_ds,
  val_ds = tf.keras.utils.img_dataset_from_directory(
      img_folder,
      validation_split=0.2,
      subset="both",
      seed=1337,
      image_size=image_size,
      batch_size=batch_size,
  )
  data_augmentation = keras.Sequential(
   [layers.RandomFlip("horizontal"),
   layers.RandomRotation(0.1),]
  )
  train_ds = train_ds.map(
    lambda img, label: (data_augmentation(img), label),
    num_parallel_calls=tf.data.AUTOTUNE,
  )
  train_ds = train_ds.prefetch(tf.data.AUTOTUNE)
  val_ds = val_ds.prefetch(tf.data.AUTOTUNE)
  return train_ds, val_ds
```

<br/>

For managing models we have to use `@model` decorator. We have some other features built in such as data lineage and branching. For more information refer [docs](https://pureml.mintlify.app/core-concepts/model).

```python
from tensorflow.keras.applications.inception_v3
import InceptionV3
from tensorflow.keras.preprocessing import image
from tensorflow.keras.models import Model
from tensorflow.keras.layers import Dense,
GlobalAveragePooling2D, Input
from pureml.decorators import model

@model("pet_classifier:dev")
def train_model(train_ds, val_ds):
  input_tensor = Input(shape=(180, 180, 3))
  base_model = InceptionV3(
   input_tensor=input_tensor,
   weights='imagenet',
   include_top=False
  )
  x = base_model.output
  x = GlobalAveragePooling2D()(x)
  x = Dense(1024, activation='relu')(x)
  predictions = Dense(1, activation='softmax')(x)
  model_inc = Model(
   inputs=base_model.input,
   outputs=predictions
  )
  model_inc.compile(
   optimizer='rmsprop',
   loss='binary_crossentropy',
   metrics=["accuracy"]
  )
  model_inc.fit(
    train_ds,
    epochs=2,
    validation_data=val_ds,
    )
  return model_inc
```

<br/>

### 2. PureML-eval : Testing & Quality Control

#### Step A: Use an existing model for validation

```python
import pureml

pureml.dataset.validation(“petdata:dev:v1”)
```

If you want to add a dataset as validation while saving it, you can use our `@validation`. This helps us capture not just one instance of this dataset but all the future variations without any intervention.

<br/>

#### Step B: Register validation dataset

```python
import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers
from pureml.decorators import dataset, validation

@validation
@dataset("petdata:dev")
def load_data(img_folder = "PetImages"):
  image_size = (180, 180)
  batch_size = 16
  train_ds,
  val_ds = tf.keras.utils.img_dataset_from_directory(
    img_folder,
    validation_split=0.2,
    subset="both",
    seed=1337,
    image_size=image_size,
    batch_size=batch_size,
  )
  data_augmentation = keras.Sequential(
   [
     layers.RandomFlip("horizontal"),
     layers.RandomRotation(0.1),
   ]
  )
  train_ds = train_ds.map(
    lambda img, label: (data_augmentation(img), label),
    num_parallel_calls=tf.data.AUTOTUNE,
  )
  train_ds = train_ds.prefetch(tf.data.AUTOTUNE)
  val_ds = val_ds.prefetch(tf.data.AUTOTUNE)
  return train_ds, val_ds
```

<br/>

#### Step C: Predictor for model

We recommend utilizing our base predictor class when developing your model. By doing so, you can leverage the predict function in this class as your model's prediction function, which can be used in various stages such as testing, inference, and dockerization.

```python
from pureml import BasePredictor
import pureml
import tensorflow as tf
from tensorflow import keras

class Predictor(BasePredictor):
  model_details = ['pet_classifier:dev:latest']
  input={'type': 'image'},
  output={'type': 'numpy ndarray' }

  def load_models(self):
    self.model = pureml.model.fetch(self.model_details)

  def predict(self, pred_img):
    pred_img = keras.preprocessing.image.img_to_array(
      pred_img
    )
    pred_img = tf.expand_dims(pred_img, 0)
    predictions = self.model.predict(pred_img)
    predictions = float(predictions[0])

    return predictions
```

<br/>

#### Step D: Evaluating your model is done as follows

Lets see how PureML makes it easier to identify and correct any issues with its review feature and allows you to evaluate the quality of their data and the accuracy of their model.

```python
import pureml

pureml.model.evaluate("pet_classifier:dev:v1", "petdata:dev:v1")
```

![Review](/assets/ReviewModel.png)

For more detailed explanation, please visit our [Documentation](https://pureml.mintlify.app) for more reference.

### 3. PureML-package

PureML is a versatile tool that allows you to package your machine learning models into a standard, production-ready container. Additionally, you can utilize a user-friendly web interface to demonstrate your machine learning model, making it easily accessible to anyone, from anywhere.

Docker

```python
pureml.docker.create(“pet_classifier:dev:v1”)
```

FastAPI

```python
pureml.fastapi.create(“pet_classifier:dev:v1”)
```

<br/>

### 4. PureML-deploy

PureML gives you the ability to deploy machine learning models without the need for managing infrastructure or servers.

```bash
pureml deploy pet_classifier:dev:v1
```

<br/>

## 💻 Self-Host on Local Machine

### Requirements

A subsystem with docker-engine installed. [Installation guidelines](https://docs.docker.com/engine/install/).

> _Currently docker-compose is the best way to self-host the official images [puremlhq/pureml_backend](https://hub.docker.com/r/puremlhq/pureml_backend) and [puremlhq/pureml_frontend](https://hub.docker.com/r/puremlhq/pureml_frontend)_

### Using [Docker Compose](https://github.com/PureMLHQ/PureML/blob/main/packages/pureml_docker/docker-compose.yml)

1. Create a new file `docker-compose.yml`

2. Add the following content from our [official docker-compose example file](https://raw.githubusercontent.com/PuremlHQ/PureML/main/packages/pureml_docker/docker-compose.yml).

```yml
version: "3"

services:
  backend:
    image: puremlhq/pureml_backend:local-base
    environment:
      - PURE_SITE_BASE_URL=http://localhost:3000
    ports:
      - 8080:8080
    volumes:
      - pureml-data:/pureml_backend/data

  frontend:
    image: puremlhq/pureml_frontend:dev
    environment:
      - NEXT_PUBLIC_BACKEND_URL=http://backend:8080/api/
    ports:
      - 3000:3000
    links:
      - backend

volumes:
  pureml-data:
```

### Running the containers

Check out the [official documentation](https://docs.docker.com/compose/gettingstarted/) for more information on running docker compose

```bash
docker compose up
```

> Additionally, to run the containers in background you can use the command `docker compose up -d` or `docker compose up --detach`.

The PureML UI will be available at http://localhost:3000. For more information about supported environment variables, please consult the documentation for [Environment Variables](https://github.com/PureMLHQ/PureML/blob/main/packages/pureml_docker/README.md).

For further information, go through this [self hosting PureML docs](https://pureml.mintlify.app/self-hosting).

<br />

## 🏁 Demo

PureML quick start demo in just 2 mins.

#### For models:

[![PureML Demo Video](http://i3.ytimg.com/vi/gqHjN4kPywo/hqdefault.jpg)](https://www.youtube.com/watch?v=gqHjN4kPywo "PureML Demo Video")
<br/>

#### For datasets:

[![PureML Demo Video](http://i3.ytimg.com/vi/gqHjN4kPywo/hqdefault.jpg)](https://www.youtube.com/watch?v=HSUh66nZbyo "PureML Demo Video")

<sub><i>Click the image to play video</i></sub>
<br/>

### Live demo

Build and run a PureML project to create data lineage and a model with our <b>[demo colab link](https://colab.research.google.com/drive/1LlrpaKiREwgesaRcnwkJP-w2MPesXf1t?usp=sharing)</b>.

<br />

## 📍 [Main Features](https://pureml.mintlify.app)

|                 |                                                                                                                                                                                                               |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Version Control | Easily version your models and datasets, ensuring that you have a clear record of all changes made throughout the development process.                                                                        |
| Commit Process  | Use review commit process that helps you ensure that only the best and most reliable models and datasets are shipped to your customers.                                                                       |
| Packaging       | Whether you need to package your models into Docker, Gradio, or FastAPI, Pureml has you covered. Our platform allows you to easily package your models in the format that works best for your specific needs. |
| Data Lineage    | Get full visibility into your datasets' lineage, making it easy to trace back any issues that may arise during development or deployment.                                                                     |
| Branches        | Enabling teams to work on multiple versions of a model or dataset simultaneously by customising in different branches for both models and datasets.                                                           |
| Testing         | Pureml includes testing for ML, making it easy to ship models reliably to your customers. With Pureml, you can be confident that your models will work as expected, every time.                               |

<br />

## ⚙ Core concepts

These are the fundamental concepts that PureML uses to operate.

|         |                                                                                                                                                                                                                                                                               |
| ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Model   | Centralized location for users to store their models and manage their lifecycle collaboratively. This makes it easier for stakeholders to manage models and promotes transparency in accessing models for tests, deployment, audit, and other purposes.                       |
| Dataset | Serves as an empty container for storing the elements of the datasets and contains lineage, dataset-related graphs, and dataset files.                                                                                                                                        |
| Log     | Provides the ability to log a range of metadata related to models and datasets. The specific types of metadata that can be logged via the PureML package and how they are represented in the PureML app will depend on the data type and format being used.                   |
| Lineage | Enables the tracking of the data flow from its origin to the end goal, which includes all the intermediate processes and transformations. In the context of PureML, lineage involves capturing the provenance of data and transformations applied to produce a final dataset. |

<br />

## 🤝 Why to get involved

Version control is much more common in software than in machine learning. So why isn’t everyone using Git? Git doesn’t work well with machine learning. It can’t handle large files, it can’t handle key/value metadata like metrics, and it can’t record information automatically from inside a training script.

GitHub wasn’t designed with data as a core project component. This along with a number of other differences between AI and more traditional software projects makes GitHub a bad fit for artificial intelligence, contributing to the reproducibility crisis in machine learning.

From manually tracking models to git based versioning systems that do not follow an intuitive versioning mechanism, there is no standardized way to track objects. Using these mechanisms, it is hard enough to track or get your model from a month ago running, let alone of a teammates!

We are trying to build a version control system for machine learning objects. A mechanism that is object dependant and intuitive for users.

Lets build this together. If you have faced this issue or have worked out a similar solution for yourself, please join us to help build a better system for everyone.

<br />

## 🧮 Tutorials

- [Registering Data lineage](https://pureml.mintlify.app/core-concepts/lineage)
- [Registering models](https://pureml.mintlify.app/version/version-model#register-models)
- [Quick Start: Tabular](https://docs.pureml.com/docs/get-started/quickstart_tabular)
- [Quick Start: Computer Vision](https://docs.pureml.com/docs/get-started/quickstart_cv)
- [Quick Start: NLP](https://docs.pureml.com/docs/get-started/quickstart_nlp)
- [Logging](https://pureml.mintlify.app/core-concepts/log)

<br />

## 🐞 Reporting Bugs

To report any bugs you have faced while using PureML package, please

1. Report it in [Discord](https://discord.gg/xNUHt9yguJ) channel
2. Open an [issue](https://github.com/PureMLHQ/PureML/issues)

<br />

## ⌨ Contributing and Developing

Lets work together to improve the features for everyone. Here's step one for you to go through our [Contributing Guide](./CONTRIBUTING.md). We are already waiting for amazing ideas and features which you all have got.

Work with mutual respect. Please take a look at our public [Roadmap here](https://pureml.notion.site/7de13568835a4cf18913307503a2cdd4?v=82199f96833a48e5907023c8a8d565c6).

<br />

## 👨‍👩‍👧‍👦 Community

To get quick updates of feature releases of PureML, follow us on:

[<img alt="Twitter" height="20" src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" />](https://twitter.com/getPureML) [<img alt="LinkedIn" height="20" src="https://img.shields.io/badge/LinkedIn-0077B5?style=for-the-badge&logo=linkedin&logoColor=white" />](https://www.linkedin.com/company/PuremlHQ/) [<img alt="GitHub" height="20" src="https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white" />](https://github.com/PureMLHQ/PureML) [<img alt="GitHub" height="20" src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" />](https://discord.gg/DBvedzGu)

<br/>

## 📄 License

See the [Apache-2.0](./License) file for licensing information.
