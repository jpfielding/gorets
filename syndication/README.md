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
					attrs := map[string]string{}
					for _, v := range start.Attr {
						attrs[v.Name.Local] = v.Value
					}
					listings.ListingsKey = attrs["listingsKey"]
					listings.Version = attrs["version"]
					listings.VersionTimestamp = attrs["versionTimestamp"]
					listings.Language = attrs["lang"]
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
