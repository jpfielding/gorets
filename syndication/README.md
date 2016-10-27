gorets
======

RETS syndication in Go

The attempt is to meet 2012-03 compliance.

Based on the
http://www.reso.org/schemas-for-syndication/

Find us at gophers.slack.com#gorets



```
	in := doms := ioutil.NopCloser(...)
	parser := xml.NewDecoder(listings)
	listings := Listings{}

	// minidom isnt necessary but its crazy useful for massive streams
	md := minidom.MiniDom{
			StartFunc: func(start xml.StartElement) {
				switch start.Name.Local {
				case "RETS":
					for _, v := range start.Attr {
						switch v.Name.Local {
						case "Listings":
							parser.DecodeElement(&listings, &start)
						case "Dislaimer":
							listings.Disclaimer = start.Value
						}
					}
				}
			},
			// quit on the the xml tag
			EndFunc: minidom.QuitAt("Listings"),
		}
	}
	err := md.Walk(parser, "Listing", func(segment io.ReadCloser, err error) error {
		tmp := Listing{}
		xml.NewDecoder(segment).Decode(&tmp)
		// .... process the listing here
		return err
	})

```
