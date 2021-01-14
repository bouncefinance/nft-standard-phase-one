package runtime

import (
	"Ankr-gin-ERC721/pkg/eventLoop"
	"context"
)

var (
	EventLoop *eventLoop.EventLoop
)

const SUBSCRIBE_CONTRACT_721 = "SUBSCRIBE_CONTRACT_721_"
const SUBSCRIBE_CONTRACT_1155 = "SUBSCRIBE_CONTRACT_1155_"

func init() {
	EventLoop = eventLoop.New(context.Background())
}
