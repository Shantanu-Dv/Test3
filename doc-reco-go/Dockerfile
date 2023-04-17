FROM python:3.7.12-slim as build
WORKDIR /app

RUN apt-get update
RUN apt-get install -y --no-install-recommends build-essential wget pkg-config gcc musl-dev unzip

RUN pip3 install torch==1.10.1+cpu torchvision==0.11.2+cpu -f https://download.pytorch.org/whl/torch_stable.html sentence-transformers==2.1.0

RUN wget https://go.dev/dl/go1.17.6.linux-amd64.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz

COPY go.mod go.sum ./
RUN /usr/local/go/bin/go mod download
COPY . .

RUN /usr/local/go/bin/go build -o doc-reco-go main.go

RUN wget https://public.ukp.informatik.tu-darmstadt.de/reimers/sentence-transformers/v0.2/multi-qa-MiniLM-L6-cos-v1.zip \
    && mkdir -p model/sentence-transformers/multi-qa-MiniLM-L6-cos-v1 \
    && unzip -o -d model/sentence-transformers/multi-qa-MiniLM-L6-cos-v1 multi-qa-MiniLM-L6-cos-v1.zip \
    && rm -rf model/sentence-transformers/multi-qa-MiniLM-L6-cos-v1/.git


FROM python:3.7.12-slim
WORKDIR /app
COPY --from=build /usr/local/lib /usr/local/lib
COPY --from=build /app/doc-reco-go .
COPY --from=build /app/model/ model
COPY internal/provider/bert_encoder/pyencoder internal/provider/bert_encoder/pyencoder
ENTRYPOINT ["./doc-reco-go"]