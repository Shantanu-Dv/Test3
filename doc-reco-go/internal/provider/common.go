package provider

import (
	"log"
	"sync"

	"doc-reco-go/internal/provider/elasticsearch"
	"doc-reco-go/internal/provider/monitoring"
)

func InitializeProvider() error {
	var err error

	wg := sync.WaitGroup{}
	wg.Add(3)

	//go func() {
	//	defer wg.Done()
	//
	//	if err = bertEncoder.InitBertModel(true); err != nil {
	//		log.Fatalln(err)
	//	}
	//}()

	go func() {
		defer wg.Done()
		if err = elasticsearch.InitializeEsClients(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err = monitoring.InitializeSentry(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err = monitoring.InitializeDatadog(); err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()
	return nil
}

func ReleaseProviderResources() {
	monitoring.StopDdTracer()
	//bertEncoder.Encoder.DeferEncoder()
}
