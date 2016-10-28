RESO Syndication
======

RETS syndication in Go

The attempt is to meet 2012-03 compliance.

Based on the

http://www.reso.org/schemas-for-syndication/

and Chris Ridenour's work mapping syndication (http://github.com/cridenour)

Find us at gophers.slack.com#gorets

```
	in := ioutil.NopCloser(...)
	parser := xml.NewDecoder(in)
	listings := Listings{}

	// minidom isnt necessary but its crazy useful for massive streams
	md := minidom.MiniDom{
			StartFunc: func(start xml.StartElement) {
				switch start.Name.Local {
				case "Listings":
					for _, v := range start.Attr {
						switch v.Name.Local {
						case "listingsKey":
							listings.ListingsKey = v.Value
						case "version":
							listings.version = v.Value
						case "versionTimestamp":
							listings.VersionTimestamp = v.Value
						case "lang":
							listings.Language = v.Value
						}
					}
				}
			},
			// quit on the the xml tag
			EndFunc: minidom.QuitAt("Listings"),
		}
	}
	err := md.Walk(parser, "Listing", ToListing(func(l Listing, err error) error {
		// .... process the listing here
		return err
	}))

```
