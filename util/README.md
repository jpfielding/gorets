RETS Utilities
======

RETS utilities for

- pull compact metadata incrementally
- converting comact metaadata to StandardXML

```
	compact := &retsutil.IncrementalCompact{}
	_ = compact.Load(sess, ctx, urls.GetMetadata)
	// convert it to the Standard XML model
	return retsutil.AsStandard(*compact).Convert()
```
