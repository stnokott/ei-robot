package store

import (
	"encoding/binary"
	"errors"
	"log"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

type Store struct {
	db *badger.DB
}

func encodeValue(t time.Time) []byte {
	return []byte(t.Format(time.RFC3339))
}

func decodeValue(b []byte) (time.Time, error) {
	return time.Parse(time.RFC3339, string(b))
}

func NewStore(dataDir string) (*Store, error) {
	opts := badger.DefaultOptions(dataDir)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	log.Printf("BadgerDB using data dir %s", dataDir)
	return &Store{
		db: db,
	}, nil
}

func (s *Store) Put(k int64, v time.Time) (err error) {
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(i2b(k), encodeValue(v))
	})
	if err == nil {
		log.Printf("key=%d, value=%s written to DB", k, v.Format(time.RFC3339))
	}
	return
}

var ErrKeyNotFound = errors.New("key not found")

func (s *Store) Get(k int64) (t time.Time, err error) {
	var v *badger.Item
	err = s.db.View(func(txn *badger.Txn) error {
		var err error
		v, err = txn.Get(i2b(k))
		return err
	})
	if err == nil {
		err = v.Value(func(val []byte) error {
			t, err = decodeValue(val)
			return err
		})
	} else if errors.Is(err, badger.ErrKeyNotFound) {
		err = ErrKeyNotFound
	}
	return
}

func (s *Store) Delete(k int64) (err error) {
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(i2b(k))
	})
	if err == nil {
		log.Printf("key=%d deleted from DB", k)
	}
	return
}

type Entry struct {
	ChatID int64
	T      time.Time
}

type ExpiryDates struct {
	Past     []*Entry
	Today    []*Entry
	Tomorrow []*Entry
}

func (s *Store) GetEggExpiries() (*ExpiryDates, error) {
	var past []*Entry
	var today []*Entry
	var tomorrow []*Entry
	now := time.Now().Truncate(24 * time.Hour)
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := b2i(item.Key())
			err := item.Value(func(v []byte) error {
				if t, err := decodeValue(v); err != nil {
					return err
				} else if time.Now().After(t) {
					// past
					past = append(past, &Entry{k, t})
				} else if now.Equal(t.Truncate(24 * time.Hour)) {
					// today
					today = append(today, &Entry{k, t})
				} else if now.Add(24 * time.Hour).Equal(t.Truncate(24 * time.Hour)) {
					// tomorrow
					tomorrow = append(tomorrow, &Entry{k, t})
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return &ExpiryDates{past, today, tomorrow}, nil
	}
}

func (s *Store) Close() error {
	return s.db.Close()
}

func i2b(v int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return b
}

func b2i(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}
