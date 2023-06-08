package entity

import (
	"time"

	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LinkOffserType string

const (
	LINK_OFFSET_AFTER = LinkOffserType("after")
)

func LinkOffserTypeFromProtobuf(offset proto.GetDialogLinksByIDRequest_OffsetType) LinkOffserType {
	// single option at the moment
	return LINK_OFFSET_AFTER
}

type LinkID struct {
	ID uint64
}

func LinkIDFromProtobuf(linkID *proto.LinkID) (*LinkID, error) {
	if linkID == nil || linkID.GetID() == 0 {
		return nil, ErrWrongRequestFormat
	}
	return &LinkID{
		ID: linkID.ID,
	}, nil
}

func LinkIDToProtobuf(linkID *LinkID) *proto.LinkID {
	return &proto.LinkID{
		ID: linkID.ID,
	}
}

type Link struct {
	LinkID    LinkID
	Link      string
	MessageID MessageID
	CreatedAt time.Time
}

func LinkToProtobuf(link *Link) *proto.Link {
	return &proto.Link{
		LinkID:    LinkIDToProtobuf(&link.LinkID),
		Link:      link.Link,
		MessageID: MessageIDToProtobuf(&link.MessageID),
		CreatedAt: timestamppb.New(link.CreatedAt),
	}
}

func LinksToProtobuf(links []*Link) []*proto.Link {
	result := make([]*proto.Link, 0, len(links))
	for _, link := range links {
		result = append(result, LinkToProtobuf(link))
	}
	return result
}
