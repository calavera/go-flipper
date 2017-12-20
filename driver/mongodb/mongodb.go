package mongodb

import (
	"github.com/calavera/go-flipper/driver"
	"github.com/calavera/go-flipper/feature"
	"github.com/calavera/go-flipper/gates"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const defaultCollectionName = "flipper"

type config struct {
	URL        string `mapstructure:"url"`
	Database   string `mapstructure:"database"`
	Collection string `mapstructure:"collection"`
}

// Driver is a store driver that keeps features and gates in mongoDB.
type Driver struct {
	collection *mgo.Collection
}

// NewDriver initializes a new mongoDB driver.
func NewDriver() *Driver {
	return &Driver{}
}

// Configure configures the mongodb driver.
// These are the options for this driver:
//   - url: string url to the mongoDB cluster (required)
//   - database: database name (optional - default to the database in the url, or "test" if also empty)
//   - collection: collection name (optional - default "flipper")
func (a *Driver) Configure(c map[string]interface{}) error {
	var conf config
	if err := mapstructure.Decode(c, &conf); err != nil {
		return errors.Wrap(err, "error decoding Mongodb's driver configuration")
	}

	if conf.URL == "" {
		return errors.New("invalid connection URL for Mongodb's driver")
	}

	if conf.Collection == "" {
		conf.Collection = defaultCollectionName
	}

	session, err := mgo.Dial(conf.URL)
	if err != nil {
		return errors.Wrap(err, "error connecting to Mongodb")
	}

	a.collection = session.DB(conf.Database).C(conf.Collection)

	return nil
}

// Enable opens a feature for a give gate.
func (a *Driver) Enable(feature feature.Feature, gate gates.Gate) error {
	var err error
	key := string(gate.Key())

	if g, ok := gate.(gates.IntGateType); ok {
		set := bson.M{"$set": bson.M{key: g.IntValue()}}
		_, err = a.collection.UpsertId(feature.Name, set)
	} else if _, ok := gate.(gates.BoolGateType); ok {
		set := bson.M{"$set": bson.M{key: true}}
		_, err = a.collection.UpsertId(feature.Name, set)
	} else if g, ok := gate.(gates.SetGateType); ok {
		set := make([]string, 0, len(g.SetValue()))
		for k := range g.SetValue() {
			set = append(set, k)
		}
		up := bson.M{"$addToSet": bson.M{key: bson.M{"$each": set}}}
		_, err = a.collection.UpsertId(feature.Name, up)
	} else {
		err = errors.Errorf("unsupported data type: %v", gate.Key())
	}

	return err
}

// Disable closes a feature for a given gate.
func (a *Driver) Disable(feature feature.Feature, gate gates.Gate) error {
	var err error
	key := string(gate.Key())

	if g, ok := gate.(gates.IntGateType); ok {
		set := bson.M{"$set": bson.M{key: g.IntValue()}}
		_, err = a.collection.UpsertId(feature.Name, set)
	} else if _, ok := gate.(gates.BoolGateType); ok {
		err = a.collection.RemoveId(feature.Name)
	} else if g, ok := gate.(gates.SetGateType); ok {
		set := make([]string, 0, len(g.SetValue()))
		for k := range g.SetValue() {
			set = append(set, k)
		}
		up := bson.M{"$pull": bson.M{key: bson.M{"$in": set}}}
		_, err = a.collection.UpsertId(feature.Name, up)
	} else {
		err = errors.Errorf("unsupported data type: %v", gate.Key())
	}

	return err
}

// Get returns the enabled gates for a feature given a set of gate keys.
// Gates are skipped if they are not open for a feature.
func (a *Driver) Get(feature feature.Feature, keys []gates.GateKey) ([]gates.Gate, error) {
	var g []gates.Gate

	result := bson.M{}
	if err := a.collection.FindId(feature.Name).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return g, nil
		}
		return nil, err
	}

	for _, t := range keys {
		v, ok := result[string(t)]
		if !ok {
			continue
		}

		switch t {
		case gates.BoolGateKey:
			g = append(g, gates.NewBoolGate(ok))
		case gates.ActorGateKey:
			gs, ok := v.([]interface{})
			if !ok {
				return nil, errors.Errorf("unexpected set value stored: %v", v)
			}
			set := gates.Set{}
			for _, k := range gs {
				s := k.(string)
				set[s] = s
			}
			g = append(g, gates.NewActorGate(set))
		case gates.GroupGateKey:
			gs, ok := v.([]interface{})
			if !ok {
				return nil, errors.Errorf("unexpected set value stored: %v", v)
			}
			set := gates.Set{}
			for _, k := range gs {
				s := k.(string)
				set[s] = s
			}
			g = append(g, gates.NewGroupGate(set))
		case gates.PercentageOfActorsGateKey:
			gi, ok := v.(int)
			if !ok {
				return nil, errors.Errorf("unexpected int value: %v", v)
			}
			g = append(g, gates.NewPercentageOfActorsGate(gi))
		case gates.PercentageOfTimeGateKey:
			gi, ok := v.(int)
			if !ok {
				return nil, errors.Errorf("unexpected int value: %v", v)
			}
			g = append(g, gates.NewPercentageOfTimeGate(gi))
		default:
			return nil, errors.Errorf("unsupported gate: %v", t)
		}
	}

	return g, nil
}

func init() {
	driver.Init("mongodb", NewDriver())
}
