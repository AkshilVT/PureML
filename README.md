<a href="https://pureml.com" align="center">
  <img
    align="center"
    alt="PureML"
    src="/assets/PureMLCoverImg.png"
    style="width:100%;"
  />
</a>

<div align="center">

# The next-gen developer platform for Production ML.

</div>

## ⛳ Quick Access

<p align="center">
  <a
    href="https://docs.pureml.com"
  ><b>Documentation</b></a>
  &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;
  <a
    href="https://www.youtube.com/watch?v=HdzLFEWS4s8&t=1s"
  ><b>Watch Demo</b></a>
  &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;
  <a
    href="https://docs.pureml.com/docs/get-started/quickstart_tabular"
  ><b>Quick example</b></a>
  &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;
  <a
    href="#"
  ><b>Get Instant Help</b></a>
  &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;
  <a
    href="https://app.pureml.com/auth/signup"
  ><b>Sign Up for free</b></a>
    &nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;
  <a
    href="https://pureml.notion.site/7de13568835a4cf18913307503a2cdd4?v=82199f96833a48e5907023c8a8d565c6"
  ><b>Roadmap</b></a>
</p>
</br>
</br>

<div align="center">
  <a
    href="https://pypi.org/project/pureml/"
  >
    <img alt="PyPi" src="https://img.shields.io/pypi/v/pureml?color=green&logo=pureml" />
  </a>
  &nbsp;
  <a
    href="https://python-poetry.org/"
  >
    <img src="https://img.shields.io/badge/poetry-1.1.14-blue?style=flat&logo=poetry&logoColor=white" />
  </a>
  &nbsp;
  <a
    href="https://opensource.org/licenses/Apache-2.0"
  >
    <img alt="License" src="https://img.shields.io/pypi/l/pureml?color=red&logo=Apache&logoColor=red" />
  </a>
  &nbsp;
  <a
    href="https://discord.gg/xNUHt9yguJ"
  >
    <img alt="Discord" src="https://img.shields.io/badge/Discord-Join%20Discord-blueviolet?style=flat&logo=discord&logoColor=white" />
  </a>
  &nbsp;
  <a
    href="https://pepy.tech/project/pureml"
  >
    <img alt="Downloads" src="https://static.pepy.tech/badge/pureml">
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
    href="https://pypi.org/project/pureml/"
  >
    <img alt="Coverage" src="https://img.shields.io/github/stars/PureMLHQ/PureML?style=social">
  </a>

</div>

## 💎 Intro

PureML is an open-source version control for machine learning.

1. [Quick start](#quick-start)
2. [How it works](#how-it-works)
3. [Demo](#demo)
4. [Main Features](#main-features)
5. [Core design principles](#core-design-principles)
6. [Core abstractions](#core-abstractions)
7. [Why to get involved](#why-to-get-involved)
8. [Tutorials](#tutorials)

<br />

## ⏱ Quick start

Manage versioning of datasets and models with our python SDK. Versioning is semantic and managed automatically. You can install and run PureML using `pip`.

Getting started is simple:

```bash
pip install pureml
```

<br />

## 📋 How it works

Just add a few lines of code. You don't need to change the way you work. PureML is a Python library that uploads metadata to S3.

If you are trying to manage versions of dataset all you have to do is use our decorator `@dataset`.

For managing models we have to use `@model` decorator. We have some other features built in such as data lineage and branching. For more information refer docs.

<br/>

#### Step 1: Use an existing model for validation

```python
import pureml

pureml.dataset.validation(“petdata:dev:v1”)
```

If you want to add a dataset as validation while saving it, you can use our `@validation`. This helps us capture not just one instance of this dataset but all the future variations without any intervention.

#### Step 2: Register validation dataset

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

#### Step 3: Predictor for model

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

#### Step 4: Evaluating your model is done as follows

```python
import pureml

pureml.model.evaluate("pet_classifier:dev:v1", "petdata:dev:v1")
```

<br/>

Lets see how PureML makes it easier to identify and correct any issues with its review feature and allows you to evaluate the quality of their data and the accuracy of their model.

<h1 align="center">
<img
   align="center"
   src="/assets/ReviewDataset.svg"
    />
</h1>

For more detailed explanation, please visit our [Documentation](https://docs.pureml.com) for more reference.

## 💻 Demo

PureML quick start demo in just 2 mins.

[![PureML Demo Video](https://img.youtube.com/vi/HdzLFEWS4s8/0.jpg)](https://www.youtube.com/watch?v=HdzLFEWS4s8 "PureML Demo Video")
<br/>
<sub><i>Click the image to play video</i></sub>

### Live demo

Build and run a PureML project to create data lineage and a model with our <b>[demo colab link](https://colab.research.google.com/drive/1LlrpaKiREwgesaRcnwkJP-w2MPesXf1t?usp=sharing)</b>.

<br />

## 📍 [Main Features](https://docs.pureml.com/)

|                          |                                                                               |
| ------------------------ | ----------------------------------------------------------------------------- |
| Data Lineage             | Automatic generation of data lineage                                          |
| Dataset Versioning       | Object-based Automatic Semantic Versioning of datasets                        |
| Model Versioning         | Object-based Automatic Semantic Versioning of models                          |
| Comparision              | Comparing different versions of models or datasets                            |
| Branches (_Coming Soon_) | Separation between experimentation and production ready models using branches |
| Review (_Coming Soon_)   | Review and approve models, and datasets to production ready branch            |

<br />

## 🔮 Core design principles

|                                     |                                                                                                                 |
| ----------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| Easy developer experience           | An intuitive open source package aimed to bridge the gaps in data science teams                                 |
| Engineering best practices built-in | Integrating PureML functionalities in your code doesnot disrupt your workflow                                   |
| Object Versioning                   | A reliable object versioning mechanism to track changes to your datasets, and models                            |
| Data is a first-class citizen       | Your data is secure. It will never leave your system.                                                           |
| Reduce Friction                     | Have access to operations performed on data using data lineage without having to spend time on lengthy meetings |

<br />

## ⚙ Core abstractions

These are the fundamental concepts that PureML uses to operate.

|                                                                     |                                                                                                                                                         |
| ------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Project](https://docs.pureml.com/docs/projects/about_projects)     | A data science project. This is where you store datasets, models, and their related objects. It is similar to a github repository with object storage.  |
| [Lineage](https://docs.pureml.com/docs/data/register_data_pipeline) | Contains a series of transformations performed on data to generate a dataset.                                                                           |
| Data Versioning                                                     | Versioning of the data should be comprehensible to the user and should encapsulate the changes in the data, its creation mechanism, among others.       |
| Model Versioning                                                    | Versioning of the model should be comprehensible to the user and should encapuslate the changes in training data, model architecture, hyper parameters. |
| Fetch                                                               | This functionality is used to fetch registered Models, and Datasets.                                                                                    |

<br />

## 🤝 Why to get involved

Version control is much more common in software than in machine learning. So why isn’t everyone using Git? Git doesn’t work well with machine learning. It can’t handle large files, it can’t handle key/value metadata like metrics, and it can’t record information automatically from inside a training script.

GitHub wasn’t designed with data as a core project component. This along with a number of other differences between AI and more traditional software projects makes GitHub a bad fit for artificial intelligence, contributing to the reproducibility crisis in machine learning.

From manually tracking models to git based versioning systems that do not follow an intuitive versioning mechanism, there is no standardized way to track objects. Using these mechanisms, it is hard enough to track or get your model from a month ago running, let alone of a teammates!

We are trying to build a version control system for machine learning objects. A mechanism that is object dependant and intuitive for users.

Lets build this together. If you have faced this issue or have worked out a similar solution for yourself, please join us to help build a better system for everyone.

<br />

## 🧮 Tutorials

- [Registering Data lineage](https://docs.pureml.com/docs/data/register_data_pipeline)
- [Registering models](https://docs.pureml.com/docs/models/register_models)
- [Quick Start: Tabular](https://docs.pureml.com/docs/get-started/quickstart_tabular)
- [Quick Start: Computer Vision](https://docs.pureml.com/docs/get-started/quickstart_cv)
- [Quick Start: NLP](https://docs.pureml.com/docs/get-started/quickstart_nlp)
- [Logging](https://docs.pureml.com/docs/log/overview)

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
