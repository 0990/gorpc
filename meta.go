package gorpc

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type context struct {
	name string
	data interface{}
}

type MessageMeta struct {
	Codec   Codec
	Type    reflect.Type
	ID      int
	ctxList []*context
}

func (self *MessageMeta) TypeName() string {
	if self == nil {
		return ""
	}

	if self.Type.Kind() == reflect.Ptr {
		return self.Type.Elem().Name()
	}

	return self.Type.Name()
}

func (self *MessageMeta) FullName() string {
	if self == nil {
		return ""
	}

	rtype := self.Type
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	var sb strings.Builder
	sb.WriteString(path.Base(rtype.PkgPath()))
	sb.WriteString(".")
	sb.WriteString(rtype.Name())
	return sb.String()
}

func (self *MessageMeta) NewType() interface{} {
	if self.Type == nil {
		return nil
	}

	return reflect.New(self.Type).Interface()
}

func (self *MessageMeta) SetContext(name string, data interface{}) *MessageMeta {
	for _, ctx := range self.ctxList {
		if ctx.name == name {
			ctx.data = data
			return self
		}
	}

	self.ctxList = append(self.ctxList, &context{
		name: name,
		data: data,
	})
	return self
}

func (self *MessageMeta) GetContext(name string) (interface{}, bool) {
	for _, ctx := range self.ctxList {
		if ctx.name == name {
			return ctx.data, true
		}
	}
	return nil, false
}

var (
	metaByFullName = map[string]*MessageMeta{}
	metaByID       = map[int]*MessageMeta{}
	metaByType     = map[reflect.Type]*MessageMeta{}
)

func RegisterMessageMeta(meta *MessageMeta) *MessageMeta {
	if _, ok := metaByType[meta.Type]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by type:%d name:%s", meta.ID, meta.Type.Name()))
	} else {
		metaByType[meta.Type] = meta
	}

	if _, ok := metaByFullName[meta.FullName()]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by fullname:%s", meta.FullName()))
	} else {
		metaByFullName[meta.FullName()] = meta
	}

	if meta.ID == 0 {
		panic("message meta require `ID` field:" + meta.TypeName())
	}

	if prev, ok := metaByID[meta.ID]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by id:%d type:%s,pre type:%s", meta.ID, meta.TypeName(), prev.TypeName()))
	} else {
		metaByID[meta.ID] = meta
	}
	return meta
}

func MessageMetaByFullName(name string) *MessageMeta {
	if v, ok := metaByFullName[name]; ok {
		return v
	}
	return nil
}

func MessageMetaVisit(nameRule string, callback func(meta *MessageMeta) bool) error {
	exp, err := regexp.Compile(nameRule)
	if err != nil {
		return err
	}

	for name, meta := range metaByFullName {
		if exp.MatchString(name) {
			if !callback(meta) {
				return nil
			}
		}
	}
	return nil
}

func MessageMetaByType(t reflect.Type) *MessageMeta {
	if t == nil {
		return nil
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v, ok := metaByType[t]; ok {
		return v
	}
	return nil
}

func MessageMetaByMsg(msg interface{}) *MessageMeta {
	if msg == nil {
		return nil
	}

	return MessageMetaByType(reflect.TypeOf(msg))
}

func MessageMetaByID(id int) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}
	return nil
}

func MessageToName(msg interface{}) string {
	if msg == nil {
		return ""
	}

	meta := MessageMetaByMsg(msg)
	if meta == nil {
		return ""
	}
	return meta.TypeName()
}

func MessageToID(msg interface{}) int {
	if msg == nil {
		return 0
	}

	meta := MessageMetaByMsg(msg)
	if meta == nil {
		return 0
	}
	return int(meta.ID)
}

func MessageSize(msg interface{}) int {
	if msg == nil {
		return 0
	}

	meta := MessageMetaByType(reflect.TypeOf(msg))

	if meta == nil {
		return 0
	}

	raw, err := meta.Codec.Encode(msg)
	if err != nil {
		return 0
	}

	return len(raw.([]byte))
}

func MessageToString(msg interface{}) string {
	if msg == nil {
		return ""
	}

	if stringer, ok := msg.(interface {
		String() string
	}); ok {
		return stringer.String()
	}
	return ""
}
