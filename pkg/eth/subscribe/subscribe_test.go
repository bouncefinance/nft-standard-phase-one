package subscribe

import (
	"Ankr-gin-ERC721/pkg/setting"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"testing"
)

func TestSubscribe721(t *testing.T) {
	chainID := 56
	ch := make(chan *types.Header)
	sub, err := setting.ETHClients[chainID].SubscribeNewHead(context.Background(), ch)
	if err != nil {
		t.Error(err)
		return
	}

	for {
		select {
		case e:=<-sub.Err():
			fmt.Printf("error %s\n",e)
		case c := <-ch:
			fmt.Println(c.Number.String())
		}

	}

}
