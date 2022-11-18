package utils

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gososy/sorpc/log"
)

func Pb2Json(msg proto.Message) (string, error) {
	var m = jsonpb.Marshaler{
		EmitDefaults: true,
		OrigName:     true,
	}
	j, err := m.MarshalToString(msg)
	if err != nil {
		log.Error("proto MarshalToString err:", err)
		return "", err
	}
	return j, nil
}

func Pb2JsonV2(msg proto.Message) (string, error) {
	j, err := json.Marshal(msg)
	if err != nil {
		log.Error("proto MarshalToString v2 err:", err)
		return "", err
	}
	return string(j), nil
}

func Pb2JsonDoNotSkipDefaults(msg proto.Message) (string, error) {
	var marshalerSkipDefaults = jsonpb.Marshaler{
		EmitDefaults: true,
		OrigName:     true,
	}
	j, err := marshalerSkipDefaults.MarshalToString(msg)
	if err != nil {
		log.Error("proto MarshalToString err:", err)
		return "", err
	}
	return j, nil
}

func Pb2JsonSkipDefaults(msg proto.Message) (string, error) {
	var marshalerSkipDefaults = jsonpb.Marshaler{
		EmitDefaults: false,
		OrigName:     true,
	}
	j, err := marshalerSkipDefaults.MarshalToString(msg)
	if err != nil {
		log.Error("proto MarshalToString err:", err)
		return "", err
	}
	return j, nil
}

func Json2Pb(j string, msg proto.Message) error {
	var u = jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	err := u.Unmarshal(bytes.NewReader([]byte(j)), msg)
	if err != nil {
		log.Errorf("UnmarshalString err:%s, json %s", err, j)
		return err
	}
	return nil
}
