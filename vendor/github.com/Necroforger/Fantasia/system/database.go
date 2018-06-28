package system

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/Necroforger/Fantasia/models"
	"github.com/boltdb/bolt"
)

// Bucket name constants
const (
	BucketGuilds = "guilds"
)

// Database wraps bolt.DB
type Database struct {
	*bolt.DB
}

// GetData retrieves data from the database
func (d *Database) GetData(bucket, key string, data interface{}) error {
	return d.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))
		if bkt == nil {
			return errors.New("Bucket does not exist")
		}

		gdata := bkt.Get([]byte(key))
		if gdata == nil {
			return errors.New("Key does not exist")
		}

		return gob.NewDecoder(bytes.NewReader(gdata)).Decode(data)
	})
}

// SaveData saves data to the database
func (d *Database) SaveData(bucket, key string, data interface{}) error {
	return d.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		var encoded bytes.Buffer
		err = gob.NewEncoder(&encoded).Encode(data)
		if err != nil {
			return err
		}

		return bkt.Put([]byte(key), encoded.Bytes())
	})
}

// GetGuild retreives saved guild settings from the database
func (d *Database) GetGuild(id string) (*models.Guild, error) {
	var guild *models.Guild
	err := d.GetData(BucketGuilds, id, guild)
	return guild, err
}

// SaveGuild saves a guild to the database
func (d *Database) SaveGuild(id string, guild *models.Guild) error {
	return d.SaveData(BucketGuilds, id, guild)
}
