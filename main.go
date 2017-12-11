package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
	cli "gopkg.in/urfave/cli.v2"
)

func loadJSONFile(path string) (*geojson.FeatureCollection, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse config file")
	}

	fc, err := geojson.UnmarshalFeatureCollection(raw)
	if err != nil {
		return nil, errors.Wrap(err, "should unmarshal feature collection without issue")

	}
	return fc, nil
}

func main() {
	app := cli.App{
		Name:     "geo-airport",
		Version:  "v1.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name: "John Daily Jr.",
			},
		},
		HelpName:  "geo-airport",
		Usage:     "maintain geojson airport",
		UsageText: "geo-airport --file ./airports.json process",
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "path to airports json `FILE`",
			EnvVars: []string{"FILE"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "iata",
			Aliases: []string{"i"},
			Usage:   "generate iata geojson",
			Action: func(c *cli.Context) error {
				// load airport file
				featureCollection, err := loadJSONFile(c.String("file"))
				if err != nil {
					fmt.Println("file path required")
					os.Exit(1)
				}
				for _, feature := range featureCollection.Features {
					iata, err := feature.PropertyString("iata")
					if err != nil {
						fmt.Println("error getting iata code")
						continue
					}
					rawJSON, err := feature.MarshalJSON()
					writeErr := ioutil.WriteFile(fmt.Sprintf("iata/%s.geojson", iata), rawJSON, 0644)
					if writeErr != nil {
						fmt.Println("error encountered writing to file")
						continue
					}
				}

				return nil
			},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "icao",
			Aliases: []string{"a"},
			Usage:   "generate icao geojson",
			Action: func(c *cli.Context) error {
				// load airport file
				featureCollection, err := loadJSONFile(c.String("file"))
				if err != nil {
					fmt.Println("file path required")
					os.Exit(1)
				}
				for _, feature := range featureCollection.Features {
					iata, err := feature.PropertyString("icao")
					if err != nil {
						fmt.Println("error getting iata code")
						continue
					}
					rawJSON, err := feature.MarshalJSON()
					writeErr := ioutil.WriteFile(fmt.Sprintf("icao/%s.geojson", iata), rawJSON, 0644)
					if writeErr != nil {
						fmt.Println("error encountered writing to file")
						continue
					}
				}

				return nil
			},
		},
	}
	app.Run(os.Args)
}
