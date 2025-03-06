package bigtableclient

import (
	"context"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/bigtable"
	"github.com/blueturbo-ad/go-utils/environment"
	k8sclient "github.com/blueturbo-ad/go-utils/k8s_tool/k8s_client"
)

func TestBigtableclient(t *testing.T) {
	os.Setenv("POD_NAMESPACE", "dsp-ns")
	environment.Init()
	k8sclient.GetSingleton().SetUp()
	t.Run("TestBigTableClient", func(t *testing.T) {
		cliobj := GetSingleton()
		if err := cliobj.UpdateLoadK8sConfigMap("bigtable", "Dev"); err != nil {
			t.Errorf("GetSingleton() = %v; want nil", err)
		}

	})
	t.Run("test bigtable insert", func(t *testing.T) {
		ctx := context.Background()
		cliobj := GetSingleton()
		if err := cliobj.UpdateLoadK8sConfigMap("bigtable", "Dev"); err != nil {
			t.Errorf("GetSingleton() = %v; want nil", err)
		}
		client := cliobj.GetClient()
		tb1 := client.Open("test-02")

		// 插入一条数据
		rowKey := "row-key-1"
		mut := bigtable.NewMutation()
		mut.Set("campaign_id", "id", bigtable.Now(), []byte("12345"))
		mut.Set("price", "value", bigtable.Now(), []byte("100"))

		err := tb1.Apply(ctx, rowKey, mut)
		if err != nil {
			t.Errorf("Could not apply mutation: %v", err)
		}

	})
	t.Run("test bigtable read", func(t *testing.T) {
		ctx := context.Background()
		cliobj := GetSingleton()
		if err := cliobj.UpdateLoadK8sConfigMap("bigtable", "Dev"); err != nil {
			t.Errorf("GetSingleton() = %v; want nil", err)
		}
		client := cliobj.GetClient()
		tb1 := client.Open("test-02")

		// 插入一条数据
		rowKey := "row-key-1"
		err := tb1.ReadRows(ctx, bigtable.PrefixRange(rowKey), func(row bigtable.Row) bool {
			for _, ris := range row {
				for _, ri := range ris {
					fmt.Printf("Row: %s, Column: %s, Value: %s\n", row.Key(), ri.Column, string(ri.Value))
				}
			}
			return true
		})
		if err != nil {
			t.Errorf("Could not read row: %v", err)
		}
	})
	// t.Run("Test mutation and read", func(t *testing.T) {
	// 	ctx := context.Background()
	// 	cliobj := GetSingleton()
	// 	if err := cliobj.UpdateLoadK8sConfigMap("bigtable", "Dev"); err != nil {
	// 		t.Errorf("GetSingleton() = %v; want nil", err)
	// 	}
	// 	client := cliobj.GetClient()
	// 	tb1 := client.Open("test-01")
	// 	var wg sync.WaitGroup
	// 	// 高并发情况下进行原子增量操作
	// 	for i := 0; i < 100; i++ {
	// 		wg.Add(1)
	// 		go func(i int) {
	// 			defer wg.Done()
	// 			mut := bigtable.NewMutation()
	// 			cond := bigtable.RowFilter(bigtable.ChainFilters(
	// 				bigtable.RowKeyFilter("campaign_id"),
	// 				bigtable.FamilyFilter(columnFamily),
	// 			))
	// 			// mut.Set("cf1", "counter", bigtable.Now(), []byte("1"))
	// 			// rowKey := "row-key"
	// 			// err := tb1.Apply(ctx, rowKey, mut)
	// 			// if err != nil {
	// 			// 	log.Printf("Could not apply mutation: %v", err)
	// 			// }
	// 		}(i)
	// 	}

	// })
}
