package weblevel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"net/http"
	"net/url"
)

const (
	Port       = 33745
	PathSets   = "/sets"
	PathGet    = "/get"
	PathDel    = "/del"
	PathStat   = "/stat"
	PathSearch = "/search"
)

type (
	WebLevel struct {
		db   *leveldb.DB
		mux  *http.ServeMux
		port int
	}
	Options struct {
		Port     int
		DBPath   string
		Database *leveldb.DB
	}
	ValKV struct {
		Key string `json:"key"`
		Val string `json:"val"`
	}
)

func NewServer(opt *Options) (*WebLevel, error) {
	if opt == nil {
		return nil, errors.New("option is nil")
	}

	if opt.Port <= 0 || opt.Port >= 65535 {
		opt.Port = Port
	}

	if opt.DBPath != "" {
		newDB, err := leveldb.OpenFile(opt.DBPath, nil)
		if err != nil {
			return nil, err
		}
		opt.Database = newDB
	}

	if opt.Database == nil {
		return nil, errors.New("db is nil")
	}
	return &WebLevel{db: opt.Database, mux: http.NewServeMux(), port: opt.Port}, nil
}

func (w *WebLevel) Bootstrap() {
	w.mux.HandleFunc(PathDel, func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		key := query.Get("key")
		logger.Warn("del", key)
		w.del(key)
	})

	w.mux.HandleFunc(PathGet, func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		key := query.Get("key")
		val, err := w.get(key)
		if err != nil {
			writer.Write(ctool.Failed[string]().JSON())
			return
		}
		writer.Write(ctool.SuccessWith(val).JSON())
	})

	// TODO /sets 错请求成 //sets body会被清空 url参数会保留
	w.mux.HandleFunc(PathSets, func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var vals []ValKV
		err := decoder.Decode(&vals)
		if err != nil {
			logger.Error(err)
			return
		}

		//logger.Info(vals)
		for _, val := range vals {
			w.set(val.Key, val.Val)
		}
	})

	w.mux.HandleFunc(PathSearch, func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		prefix := query.Get("prefix")
		cache := w.rangeKey(prefix)

		marshal, _ := json.Marshal(cache)
		writer.Write(marshal)
	})

	logger.Info("start on", w.port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", w.port), w.mux)
	if err != nil {
		logger.Error(err)
	}
}

func (w *WebLevel) del(key string) {
	err := w.db.Delete([]byte(key), nil)
	if err != nil {
		logger.Error(key, err)
	}
}

func (w *WebLevel) get(key string) (string, error) {
	value, err := w.db.Get([]byte(key), nil)
	if err != nil {
		unescape, err2 := url.QueryUnescape(key)
		logger.Warn(unescape, err, err2)
		return "", err
	}
	return string(value), nil
}

func (w *WebLevel) set(key, value string) {
	err := w.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		logger.Error(key, value, err)
	}
}

func (w *WebLevel) rangeKey(prefix string) map[string]string {
	result := make(map[string]string)
	iter := w.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for iter.Next() {
		result[string(iter.Key())] = string(iter.Value())
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		logger.Error(err)
	}
	return result
}
