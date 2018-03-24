RETS Utilities
======

RETS utilities for:

- pulling compact metadata incrementally
- converting compact metadata to StandardXML

```go
	compact := &retsutil.IncrementalCompact{}
	_ = compact.Load(sess, ctx, urls.GetMetadata)
	// convert it to the Standard XML model
	return retsutil.AsStandard(*compact).Convert()
```
