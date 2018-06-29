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

// Error variables
var (
	ErrNotFound = errors.New("not found")
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
			return ErrNotFound
		}

		gdata := bkt.Get([]byte(key))
		if gdata == nil {
			return ErrNotFound
		}

		return gob.NewDecoder(bytes.NewReader(gdata)).Decode(data)
	})
}

// SaveData saves data to the database
func (d *Database) SaveData(bucket, key string, data interface{}) error {
	return d.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return ErrNotFound
		}

		var encoded bytes.Buffer
		err = gob.NewEncoder(&encoded).Encode(data)
		if err != nil {
			return ErrNotFound
		}

		return bkt.Put([]byte(key), encoded.Bytes())
	})
}

// CreateGuildIfNotExists gets or creates a guild config if it does not exist
func (d *Database) CreateGuildIfNotExists(id string) (*models.Guild, error) {
	guild, err := d.GetGuild(id)
	if err == ErrNotFound { // If the guild is not found, create it
		g := models.NewGuild()
		err = d.SaveGuild(id, g)
		if err != nil {
			return nil, err
		}
		return g, nil
	}
	return guild, err
}

// GetGuild retreives saved guild settings from the database
func (d *Database) GetGuild(id string) (*models.Guild, error) {
	guild := &models.Guild{}
	err := d.GetData(BucketGuilds, id, guild)
	return guild, err
}

// SaveGuild saves a guild to the database
func (d *Database) SaveGuild(id string, guild *models.Guild) error {
	return d.SaveData(BucketGuilds, id, guild)
}
