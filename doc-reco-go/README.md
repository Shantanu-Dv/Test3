#Doc-reco service
Build using Go and gorilla mux. Use to search question on Elasticsearch based on user query or image <br>
BERT query is using Python (CPython v3.7 only) is used to calculate text's vector embeddings before searching on ES. Because of python rich ML ecosystem which outperforms go counterpart in terms of speed. 

## Steps to setup **Doc-reco-go** service
1. Clone this repo
2. Install go(v1.17) on your system ```sudo snap install go --classic``` or [install from binary](https://golang.org/doc/install)
3. Add app secret file here *doc-reco-go/internal/config/app_secret/doc_reco_settings.yml*
4. Execute ```go mod download``` to install all necessary dependencies

   
## TODO: :notebook:
  - Native Setup on M1 Mac. <br>
   Can be run using docker virtualization layer from arm64 to amd64 translation layer at significant low performance <br>
   <b>Build: <b> `docker build --platform linux/x86_64 -t doc-reco-go . -f Dockerfile.dev` <br>
   <b>Run: <b>   `docker run --rm --platform linux/x86_64 -doc-reco-go` <br>
    Hack: remove pyencoder import and corresponding code form `internal/provider/bert_encoder/encoder.go` and then go build

## Steps to run
Locally

    go run main.go
    or build image using Dockerfile.dev (if using docker)

On server

    go build && ./doc-reco-go
    or build image using Dockerfile


### [Preprod server](http://doc-reco-go-preprod.toppr.com/)
  Running server on port 4000  
    /home/apps/doc-reco-go - *doc-reco.service*  

   Docker: `docker run -v ${PWD}/internal/:/app/internal/ -p 4001:4000 --memory=1G doc-reco-go`

Resource for BERT integration
- [DataDog integration of Python inside Go](https://www.datadoghq.com/blog/engineering/cgo-and-python/)
- [How to use DataDog CPython Go wrapper](https://poweruser.blog/embedding-python-in-go-338c0399f3d5)

Issue encountered while setup:
- [Installing Python3.7 from source](https://github.com/DataDog/go-python3/issues/40#issuecomment-893623375)
- [Remove existing python3.7 cache](https://github.com/pyenv/pyenv/issues/1803#issuecomment-1007920388)
- [Update pkgconfig file python-3.7.pc](https://github.com/DataDog/go-python3/issues/19#issuecomment-558277077) (File path: <b>sudo find / -type f -name "python-3.7.pc"</b>)

[Sentence transformers pre-trained models](https://public.ukp.informatik.tu-darmstadt.de/reimers/sentence-transformers/v0.2/)
