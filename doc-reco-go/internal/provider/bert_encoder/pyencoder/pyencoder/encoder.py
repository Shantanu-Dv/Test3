from __future__ import annotations

import os
import time
from sentence_transformers import SentenceTransformer

model = None


def init_model(model_name: bytes):
    model_name = bytes.decode(model_name)

    global model
    t = time.time()

    model_path = f'{os.getcwd()}/model/{model_name}'
    model = SentenceTransformer(model_path)
    print(f"{model_path}, loaded in {time.time() - t} sec")


def encode(query: bytes) -> list[float]:
    query = bytes.decode(query)
    a = model.encode(query)
    return a
