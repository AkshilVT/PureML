import { CodeBlock, sunburst } from "react-code-blocks";
import Tag from "../ui/Tag";

export default function TestingSection() {
  return (
    <div className="h-fit flex flex-col gap-y-6 pt-16 md:py-16">
      <h1 className="flex items-center text-3xl md:text-4xl lg:text-5xl !text-slate-400">
        02
      </h1>
      <div className="flex flex-col gap-y-12 text-slate-600">
        <div className="flex flex-col gap-y-12">
          <div className="flex flex-col gap-y-6 md:w-3/4">
            <h1 className="font-medium text-3xl md:text-4xl lg:text-5xl pb-2">
              PureML-eval : Testing & Quality Control
            </h1>
          </div>
          <div className="flex flex-col md:flex-row gap-x-12 gap-y-6">
            <div className="md:w-1/2 flex flex-col gap-y-12">
              <div className="flex flex-col gap-y-6">
                <h2 className="text-lg md:text-xl lg:text-2xl">
                  The terms L0, L1, L2, L3, and L4 refer to teams involved in
                  developing machine learning systems. Our research study which
                  surveyed 100 companies found that the majority of teams are
                  currently at the L0 level. However, some teams are striving to
                  build highly dependable systems and are at L3 or L4 levels. We
                  are developing tools for the teams to facilitate their journey
                  to attain the L4 level.
                </h2>
                <Tag intent="landingpg" />
              </div>
            </div>
            <div className="md:w-1/2">
              <img
                src="/imgs/landingPage/MobTestingTable.svg"
                alt="Evaluate"
                className="pt-6 md:pt-0 flex md:hidden w-full"
              />
              <img
                src="/imgs/landingPage/TestingTable.svg"
                alt="Evaluate"
                className="pt-6 md:pt-0 hidden md:flex"
              />
            </div>
          </div>
          <div className="flex flex-col md:flex-row gap-x-12 gap-y-6">
            <div className="md:w-1/2 flex flex-col gap-y-4">
              <h2 className="font-medium text-lg md:text-xl lg:text-3xl text-slate-950">
                Step a.1: Use an existing model for validation
              </h2>
              <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                <CodeBlock
                  text={`import pureml

pureml.dataset.validation(“petdata:dev:v1”)`}
                  language="python"
                  theme={sunburst}
                  showLineNumbers={false}
                  wrapLines
                />
              </div>
              <h2 className="text-lg md:text-xl lg:text-2xl">
                Use{" "}
                <span className="bg-slate-200 px-2 py-1 text-lg md:text-xl lg:text-2xl">
                  @validation
                </span>
                for adding a dataset as validation while saving it.
              </h2>
              <h2 className="text-lg md:text-xl lg:text-2xl">
                This helps us capture not just one instance of this dataset, but
                all the future variations without any intervention.
              </h2>
            </div>
            <div className="md:w-1/2 flex flex-col gap-y-4">
              <h2 className="font-medium text-lg md:text-xl lg:text-3xl text-slate-950">
                Step a.2: Register validation dataset
              </h2>
              <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                {/* to highlight specific code */}
                <div className="relative">
                  <div className="overflow-hidden">
                    <div className="bg-yellow-400 opacity-30 h-[3.5rem] z-30 w-full absolute mt-[9.6rem]"></div>
                  </div>
                </div>
                <CodeBlock
                  text={`import tensorflow as tf
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
  return train_ds, val_ds`}
                  language="python"
                  theme={sunburst}
                  showLineNumbers={false}
                  wrapLines
                />
              </div>
            </div>
          </div>
          <div className="flex flex-col md:flex-row gap-x-12 gap-y-6">
            <div className="md:w-1/2 flex flex-col gap-y-4">
              <h2 className="font-medium text-lg md:text-xl lg:text-3xl text-slate-950">
                Step b: Predictor for model
              </h2>
              <h2 className="text-lg md:text-xl lg:text-2xl">
                We recommend utilizing our base predictor class when developing
                your model. By doing so, you can leverage the predict function
                in this class as your model's prediction function, which can be
                used in various stages such as testing, inference, and
                dockerization.
              </h2>
            </div>
            <div className="md:w-1/2">
              <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                {/* to highlight specific code */}
                <div className="relative">
                  <div className="overflow-hidden">
                    <div className="bg-yellow-400 opacity-30 h-8 z-30 w-full absolute mt-[9.6rem]"></div>
                  </div>
                </div>
                <CodeBlock
                  text={`from pureml import BasePredictor
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

    return predictions`}
                  language="python"
                  theme={sunburst}
                  showLineNumbers={false}
                  wrapLines
                />
              </div>
            </div>
          </div>
          <div className="flex flex-col md:flex-row gap-x-12 gap-y-6">
            <div className="md:w-1/2 flex flex-col gap-y-4">
              <h2 className="font-medium text-lg md:text-xl lg:text-3xl text-slate-950">
                Step c: Evaluating your model is done as follows
              </h2>
            </div>
            <div className="md:w-1/2">
              <div className="codeblock w-[92vw] md:w-[43vw] lg:w-full overflow-hidden md:overflow-visible">
                <CodeBlock
                  text={`import pureml

pureml.model.evaluate("pet_classifier:dev:v1", "petdata:dev:v1")`}
                  language="python"
                  theme={sunburst}
                  showLineNumbers={false}
                  wrapLines
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
