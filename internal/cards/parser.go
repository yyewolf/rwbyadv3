package cards

import (
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Card struct {
	ID         string   `yaml:"id"`
	Name       string   `yaml:"name"`
	Type       CardType `yaml:"type"`
	Team       string   `yaml:"team"`
	Categories []string `yaml:"categories"`
	BaseStats  struct {
		Damage  int `yaml:"damage"`
		Healing int `yaml:"healing"`
		Armor   int `yaml:"armor"`
		Health  int `yaml:"health"`
		Speed   int `yaml:"speed"`
	} `yaml:"base_stats"`

	// For variants
	Parent *Card `yaml:"-"`

	// Unused after parsing
	Variants []*Card `yaml:"variants"`
}

func (a *Card) doVariants() []*Card {
	var b []*Card

	// Firstly, the card itself is the first variant
	b = append(b, a)

	for _, variant := range a.Variants {
		// Then all variant inherit the card's properties
		c := *a

		// For all non-empty fields, overwrite the card's properties
		v := reflect.ValueOf(variant).Elem()
		for i := 0; i < v.NumField(); i++ {
			if !v.Field(i).IsZero() {
				reflect.ValueOf(&c).Elem().Field(i).Set(v.Field(i))
			}
		}

		c.Parent = a

		b = append(b, &c)
	}

	for _, card := range b {
		card.Variants = b
	}

	return b
}

func parseCard(location string) []*Card {
	var card Card

	data, err := os.ReadFile(location)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to read card file")
	}

	err = yaml.Unmarshal(data, &card)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to unmarshal card file")
	}

	return card.doVariants()
}

var Cards map[string]*Card

func ParseCards(location string) {
	// Location points to a folder that only contains .yml files
	folder, err := os.Open(location)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open cards folder")
	}

	files, err := folder.ReadDir(-1)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to read cards folder")
	}

	cards := make([]*Card, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Parse the file if it's a .yml file
		card := parseCard(location + "/" + file.Name())
		cards = append(cards, card...)
	}

	// Create a map of cards
	cardMap := make(map[string]*Card)

	for _, card := range cards {
		if _, ok := cardMap[card.ID]; ok {
			logrus.WithField("id", card.ID).Fatal("Duplicate card ID")
		}

		cardMap[card.ID] = card
	}

	Cards = cardMap
}
