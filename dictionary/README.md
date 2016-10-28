RESO Data Dictionary Models
======

RESO data dictionary in Go

The attempt is to meet 1.5 compliance.

Based on the

http://www.reso.org/data-dictionary/

and Chris Ridenour's work mapping data dictionary (http://github.com/cridenour)

Find us at gophers.slack.com#gorets

```
	in := doms := ioutil.NopCloser(...)
	parser := xml.NewDecoder(in)

	// minidom isnt necessary but its crazy useful for massive streams
	md := minidom.MiniDom{
			EndFunc: minidom.QuitAt("REData"),
		}
	}
	err := md.Walk(parser, "PropertyListing", func(mini io.ReadCloser, err error) error {
		if err != nil {
			return nil
		}
		property := Property{}
		xml.NewDecoder(body).Decode(&property)
		return err
	}))

```
