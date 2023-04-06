import { CodeBlock, sunburst } from "react-code-blocks";
import LandingPgTab from "./Tabs";

export default function VersionSection() {
  return (
    <div className="h-fit flex flex-col gap-y-6 pt-16 md:py-16">
      <h1 className="flex font-medium items-center text-3xl md:text-4xl lg:text-5xl !text-slate-400">
        01
      </h1>
      <div className="flex flex-col gap-y-6 md:gap-y-12 text-slate-600">
        <div className="flex flex-col gap-y-8">
          <div className="md:w-3/4">
            <h1 className="font-medium text-3xl md:text-4xl lg:text-5xl pb-2">
              PureML-version
            </h1>
          </div>
          <div className="flex flex-col md:flex-row gap-x-12 gap-y-6">
            <div className="md:w-1/2 flex flex-col gap-y-12">
              <h2 className="text-lg md:text-xl lg:text-2xl">
                Manage versioning of datasets and models with our python SDK.
                Versioning is semantic and automatically managed.
              </h2>
              <div className="flex flex-col gap-y-4">
                <div className="flex flex-col gap-y-2">
                  <div className="text-slate-950 text-lg md:text-xl lg:text-2xl font-medium">
                    Install
                  </div>
                  <h2 className="text-lg md:text-xl lg:text-2xl">
                    Getting started is simple
                  </h2>
                </div>
                <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                  <CodeBlock
                    text={`$ pip install pureml`}
                    language="bash"
                    theme={sunburst}
                    showLineNumbers={false}
                    wrapLines
                  />
                </div>
              </div>
              <div className="flex flex-col gap-y-2">
                <div className="text-slate-950 text-lg md:text-xl lg:text-2xl font-medium">
                  Dataset
                </div>
                <h2 className="text-lg md:text-xl lg:text-2xl">
                  Simply use our decorator{" "}
                  <span className="bg-slate-200 px-2 py-1 text-lg md:text-xl lg:text-2xl">
                    @dataset
                  </span>{" "}
                  for managing the versions of your dataset.
                </h2>
              </div>
              <div className="flex flex-col gap-y-6">
                <div className="flex flex-col gap-y-2">
                  <div className="text-slate-950 text-lg md:text-xl lg:text-2xl font-medium">
                    Model
                  </div>
                  <h2 className="text-lg md:text-xl lg:text-2xl">
                    Use{" "}
                    <span className="bg-slate-200 px-2 py-1 text-lg md:text-xl lg:text-2xl">
                      @model
                    </span>{" "}
                    decorator for managing models. Check out the docs for other
                    built in features such as data lineage and branching.
                  </h2>
                </div>
                <a
                  href="https://pureml.mintlify.app"
                  target="_blank"
                  rel="noreferrer"
                >
                  <button className="btn btn-primary btn-sm font-normal text-white w-full md:w-fit hover-effect px-4 rounded-lg text-lg">
                    READ DOCS
                  </button>
                </a>
              </div>
            </div>
            <div className="md:w-1/2">
              <LandingPgTab
                tab1="DATASET"
                tab2="MODEL"
                tab3=""
                tab1Content={
                  <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                    {/* to highlight specific code */}
                    <div className="relative">
                      <div className="overflow-hidden">
                        <div className="bg-yellow-400 opacity-30 h-8 z-30 w-full absolute mt-[9.6rem]"></div>
                      </div>
                    </div>
                    <CodeBlock
                      text={`import tensorflow as tf
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
  return train_ds, val_ds`}
                      language="python"
                      theme={sunburst}
                      showLineNumbers={false}
                      wrapLines
                    />
                  </div>
                }
                tab2Content={
                  <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                    {/* to highlight specific code */}
                    <div className="relative">
                      <div className="overflow-hidden">
                        <div className="bg-yellow-400 opacity-30 h-8 z-30 w-full absolute mt-[14.6rem]"></div>
                      </div>
                    </div>
                    <CodeBlock
                      text={`from tensorflow.keras.applications.inception_v3
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
  return model_inc`}
                      language="python"
                      theme={sunburst}
                      showLineNumbers={false}
                      wrapLines
                    />
                  </div>
                }
                tab3Content=""
              />
            </div>
          </div>
        </div>
        <div className="flex flex-col md:flex-row justify-between text-xl gap-y-6 gap-x-12">
          <div className="w-full">
            <h1 className="text-slate-950 text-3xl pb-1 font-medium flex gap-x-3">
              <img
                src="/imgs/landingPage/icons/FlashIcon.svg"
                alt="FlashIcon"
                className="w-6"
              />
              Key-value pairs
            </h1>
            <h2 className="text-lg md:text-xl lg:text-2xl">
              Our system captures key-value metadata such as metrics and
              associates it with the version of the model.
            </h2>
          </div>
          <div className="w-full">
            <h1 className="text-slate-950 text-3xl pb-1 font-medium flex gap-x-3">
              <img
                src="/imgs/landingPage/icons/ScalableIcon.svg"
                alt="ScalableIcon"
                className="w-6"
              />
              Large files
            </h1>
            <h2 className="text-lg md:text-xl lg:text-2xl">
              Our ML versioning system is built to natively support large files,
              unlike Git.
            </h2>
          </div>
          <div className="w-full">
            <h1 className="text-slate-950 text-3xl pb-1 font-medium flex gap-x-3">
              <img
                src="/imgs/landingPage/icons/FlexibleIcon.svg"
                alt="FlexibleIcon"
                className="w-6"
              />
              Powerful yet flexible
            </h1>
            <h2 className="text-lg md:text-xl lg:text-2xl">
              Our SDK is designed to be user-friendly, yet robust enough to meet
              a wide range of use cases.
            </h2>
          </div>
        </div>
      </div>
    </div>
  );
}
