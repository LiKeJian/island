package persist

import (
	"context"
	"errors"
	"log"
	pb "reptiles/crawler_distributed/proto"

	"github.com/olivere/elastic/v7"
)

func ItemSaver(
	index string) (chan pb.Item, error) {
	client, err := elastic.NewClient(
		// Must turn off sniff in docker
		elastic.SetSniff(false))

	if err != nil {
		return nil, err
	}

	out := make(chan pb.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: got item "+
				"#%d: %v", itemCount, item)
			itemCount++

			err := Save(client, index, &item)
			if err != nil {
				log.Printf("Item Saver: error "+
					"saving item %v: %v",
					item, err)
			}
		}
	}()

	return out, nil
}

func Save(
	client *elastic.Client, index string,
	item *pb.Item) error {

	if item.Type == "" {
		return errors.New("must supply Type")
	}

	indexService := client.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}

	_, err := indexService.
		Do(context.Background())

	return err
}
