// package s3

// import (
// 	"fmt"

// 	"github.com/chuxorg/chux-models/models/products"
// )
// func Parse() {

// 	//item := "202320/product/items_zzsounds.com-2023-03-20T18_20_42.720000.json"
// 	bucketName := "chux-crawler"
// 	//fileName := "items_zzsounds.com-2023-03-20T18_20_42.720000.json"

// 	bucket := NewBucket(bucketName)
// 	files := bucket.DownloadAll()
	
// 	product := products.New(
// 		products.WithDatabase
// 	)
	
// 	for _, f := range files {
// 		fmt.Printf("Name =%s  isProduct=%d", f.Name, f.IsProduct)
// 	}
// }
