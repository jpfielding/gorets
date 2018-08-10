package util

import (
	"github.com/jpfielding/gorets/pkg/metadata"
	"github.com/jpfielding/gorets/pkg/rets"
)

// TODO need to restructure this to return the sub types and pass the
// results back.  maybe a builder pattern?  maybe if done right an interface
// for returning results could allow a hook into either providing
// a function that allows incremental compact or standard loading for this?

// AsStandard ...
type AsStandard rets.CompactMetadata

// Convert ...
func (cm AsStandard) Convert() (*metadata.MSystem, error) {
	ms := &metadata.MSystem{}
	ms.Date = metadata.DateTime(cm.MSystem.Date)
	ms.Version = metadata.Version(cm.MSystem.Version)
	ms.System.ID = cm.MSystem.System.ID
	ms.System.Description = cm.MSystem.System.Description
	ms.System.MetadataID = cm.MSystem.System.MetadataID
	ms.System.TimeZoneOffset = cm.MSystem.System.TimeZoneOffset
	ms.System.Comments = cm.MSystem.Comments
	// TODO figure out how compact system sub-elements exists
	ms.System.MForeignKey = cm.MForeignKey()
	ms.System.MResource = cm.MResource()
	ms.System.MFilter = cm.MFilter()
	return ms, nil
}

// MForeignKey ...
func (cm AsStandard) MForeignKey() *metadata.MForeignKey {
	mfk := metadata.MForeignKey{}
	for _, fk := range cm.Elements[metadata.MIForeignKey.Name] {
		metadata.FieldTransfer(fk.Attr).To(&mfk)
		for _, entry := range fk.Entries() {
			tmp := metadata.ForeignKey{}
			entry.SetFields(&tmp)
			mfk.ForeignKey = append(mfk.ForeignKey, tmp)
		}
	}
	return &mfk
}

// MFilter ...
func (cm AsStandard) MFilter() *metadata.MFilter {
	mf := metadata.MFilter{}
	mi := metadata.MIFilter
	for _, cd := range cm.Elements[mi.Name] {
		metadata.FieldTransfer(cd.Attr).To(&mf)
		for _, entry := range cd.Entries() {
			tmp := metadata.Filter{}
			entry.SetFields(&tmp)
			id := mi.ID(tmp)
			tmp.MFilterType = cm.MFilterType(id)
			mf.Filter = append(mf.Filter, tmp)
		}
	}
	return &mf
}

// MResource ...
func (cm AsStandard) MResource() *metadata.MResource {
	mr := metadata.MResource{}
	mi := metadata.MIResource
	for _, cd := range cm.Elements[mi.Name] {
		metadata.FieldTransfer(cd.Attr).To(&mr)
		for _, entry := range cd.Entries() {
			tmp := metadata.Resource{}
			entry.SetFields(&tmp)
			id := mi.ID(&tmp)
			tmp.MClass = cm.MClass(id)
			tmp.MObject = cm.MObject(id)
			tmp.MLookup = cm.MLookup(id)
			tmp.MSearchHelp = cm.MSearchHelp(id)
			tmp.MEditMask = cm.MEditMask(id)
			tmp.MUpdateHelp = cm.MUpdateHelp(id)
			tmp.MValidationExpression = cm.MValidationExpression(id)
			tmp.MValidationExternal = cm.MValidationExternal(id)
			tmp.MValidationLookup = cm.MValidationLookup(id)
			mr.Resource = append(mr.Resource, tmp)
		}
	}
	return &mr
}

// MValidationExpression ...
func (cm AsStandard) MValidationExpression(resource string) *metadata.MValidationExpression {
	mve := metadata.MValidationExpression{}
	mi := metadata.MIValidationExpression
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mve)
		for _, entry := range cd.Entries() {
			tmp := metadata.ValidationExpression{}
			entry.SetFields(&tmp)
			mve.ValidationExpression = append(mve.ValidationExpression, tmp)
		}
	}
	return &mve
}

// MUpdateHelp ...
func (cm AsStandard) MUpdateHelp(resource string) *metadata.MUpdateHelp {
	muh := metadata.MUpdateHelp{}
	for _, cd := range cm.Elements[metadata.MIUpdateHelp.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&muh)
		for _, entry := range cd.Entries() {
			tmp := metadata.UpdateHelp{}
			entry.SetFields(&tmp)
			muh.UpdateHelp = append(muh.UpdateHelp, tmp)
		}
	}
	return &muh
}

// MEditMask ...
func (cm AsStandard) MEditMask(resource string) *metadata.MEditMask {
	mem := metadata.MEditMask{}
	for _, cd := range cm.Elements[metadata.MIEditMask.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mem)
		for _, entry := range cd.Entries() {
			tmp := metadata.EditMask{}
			entry.SetFields(&tmp)
			mem.EditMask = append(mem.EditMask, tmp)
		}
	}
	return &mem
}

// MSearchHelp ...
func (cm AsStandard) MSearchHelp(resource string) *metadata.MSearchHelp {
	msh := metadata.MSearchHelp{}
	mi := metadata.MISearchHelp
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&msh)
		for _, entry := range cd.Entries() {
			tmp := metadata.SearchHelp{}
			entry.SetFields(&tmp)
			msh.SearchHelp = append(msh.SearchHelp, tmp)
		}
	}
	return &msh
}

// MClass ...
func (cm AsStandard) MClass(resource string) *metadata.MClass {
	mc := metadata.MClass{}
	mi := metadata.MIClass
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mc)
		for _, entry := range cd.Entries() {
			tmp := metadata.Class{}
			entry.SetFields(&tmp)
			id := mi.ID(&tmp)
			tmp.MTable = cm.MTable(resource, id)
			tmp.MUpdate = cm.MUpdate(resource, id)
			tmp.MColumnGroupSet = cm.MColumnGroupSet(resource, id)
			tmp.MColumnGroup = cm.MColumnGroup(resource, id)
			mc.Class = append(mc.Class, tmp)
		}
	}
	return &mc
}

// MColumnGroup ...
func (cm AsStandard) MColumnGroup(resource, class string) *metadata.MColumnGroup {
	mcg := metadata.MColumnGroup{}
	mi := metadata.MIColumnGroup
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mcg)
		for _, entry := range cd.Entries() {
			tmp := metadata.ColumnGroup{}
			entry.SetFields(&tmp)
			id := mi.ID(&tmp)
			tmp.MColumnGroupControl = cm.MColumnGroupControl(resource, class, id)
			tmp.MColumnGroupTable = cm.MColumnGroupTable(resource, class, id)
			tmp.MColumnGroupNormalization = cm.MColumnGroupNormalization(resource, class, id)
			mcg.ColumnGroup = append(mcg.ColumnGroup, tmp)
		}
	}
	return &mcg
}

// MColumnGroupControl ...
func (cm AsStandard) MColumnGroupControl(resource, class, cg string) *metadata.MColumnGroupControl {
	mcgc := metadata.MColumnGroupControl{}
	mi := metadata.MIColumnGroupControl
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		if cd.Attr["ColumnGroup"] != cg {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mcgc)
		for _, entry := range cd.Entries() {
			tmp := metadata.ColumnGroupControl{}
			entry.SetFields(&tmp)
			mcgc.ColumnGroupControl = append(mcgc.ColumnGroupControl, tmp)
		}
	}
	return &mcgc
}

// MColumnGroupTable ...
func (cm AsStandard) MColumnGroupTable(resource, class, cg string) *metadata.MColumnGroupTable {
	mcgt := metadata.MColumnGroupTable{}
	for _, cd := range cm.Elements[metadata.MIColumnGroupTable.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		if cd.Attr["ColumnGroup"] != cg {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mcgt)
		for _, entry := range cd.Entries() {
			tmp := metadata.ColumnGroupTable{}
			entry.SetFields(&tmp)
			mcgt.ColumnGroupTable = append(mcgt.ColumnGroupTable, tmp)
		}
	}
	return &mcgt
}

// MColumnGroupNormalization ...
func (cm AsStandard) MColumnGroupNormalization(resource, class, cg string) *metadata.MColumnGroupNormalization {
	mcgn := metadata.MColumnGroupNormalization{}
	for _, cd := range cm.Elements[metadata.MIColumnGroupNormalization.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		if cd.Attr["ColumnGroup"] != cg {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mcgn)
		for _, entry := range cd.Entries() {
			tmp := metadata.ColumnGroupNormalization{}
			entry.SetFields(&tmp)
			mcgn.ColumnGroupNormalization = append(mcgn.ColumnGroupNormalization, tmp)
		}
	}
	return &mcgn
}

// MColumnGroupSet ...
func (cm AsStandard) MColumnGroupSet(resource, class string) *metadata.MColumnGroupSet {
	mcgs := metadata.MColumnGroupSet{}
	for _, cd := range cm.Elements[metadata.MIColumnGroupSet.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mcgs)
		for _, entry := range cd.Entries() {
			tmp := metadata.ColumnGroupSet{}
			entry.SetFields(&tmp)
			mcgs.ColumnGroupSet = append(mcgs.ColumnGroupSet, tmp)
		}
	}
	return &mcgs
}

// MUpdate ...
func (cm AsStandard) MUpdate(resource, class string) *metadata.MUpdate {
	mu := metadata.MUpdate{}
	mi := metadata.MIUpdate
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mu)
		for _, entry := range cd.Entries() {
			tmp := metadata.Update{}
			entry.SetFields(&tmp)
			id := mi.ID(&tmp)
			tmp.MUpdateType = cm.MUpdateType(resource, class, id)
			mu.Update = append(mu.Update, tmp)
		}
	}
	return &mu
}

// MUpdateType ...
func (cm AsStandard) MUpdateType(resource, class, update string) *metadata.MUpdateType {
	mut := metadata.MUpdateType{}
	mi := metadata.MIUpdateType
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		if cd.Attr["Update"] != update {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mut)
		for _, entry := range cd.Entries() {
			tmp := metadata.UpdateType{}
			entry.SetFields(&tmp)
			mut.UpdateType = append(mut.UpdateType, tmp)
		}
	}
	return &mut
}

// MTable ...
func (cm AsStandard) MTable(resource, class string) *metadata.MTable {
	mt := metadata.MTable{}
	mi := metadata.MITable
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mt)
		for _, entry := range cd.Entries() {
			tmp := metadata.Field{}
			entry.SetFields(&tmp)
			mt.Field = append(mt.Field, tmp)
		}
	}
	return &mt
}

// MObject ...
func (cm AsStandard) MObject(resource string) *metadata.MObject {
	mo := metadata.MObject{}
	mi := metadata.MIObject
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mo)
		for _, entry := range cd.Entries() {
			tmp := metadata.Object{}
			entry.SetFields(&tmp)
			mo.Object = append(mo.Object, tmp)
		}
	}
	return &mo
}

// MLookup ...
func (cm AsStandard) MLookup(resource string) *metadata.MLookup {
	ml := metadata.MLookup{}
	mi := metadata.MILookup
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&ml)
		for _, entry := range cd.Entries() {
			tmp := metadata.Lookup{}
			entry.SetFields(&tmp)
			id := mi.ID(&tmp)
			tmp.MLookupType = cm.MLookupType(resource, id)
			ml.Lookup = append(ml.Lookup, tmp)
		}
	}
	return &ml
}

// MLookupType ...
func (cm AsStandard) MLookupType(resource, lookup string) *metadata.MLookupType {
	mlt := metadata.MLookupType{}
	mi := metadata.MILookupType
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Lookup"] != lookup {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mlt)
		for _, entry := range cd.Entries() {
			tmp := metadata.LookupType{}
			entry.SetFields(&tmp)
			mlt.LookupType = append(mlt.LookupType, tmp)
		}
	}
	return &mlt
}

// MFilterType ...
func (cm AsStandard) MFilterType(filter string) *metadata.MFilterType {
	mft := metadata.MFilterType{}
	mi := metadata.MIFilterType
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Filter"] != filter {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mft)
		for _, entry := range cd.Entries() {
			tmp := metadata.FilterType{}
			entry.SetFields(&tmp)
			mft.FilterType = append(mft.FilterType, tmp)
		}
	}
	return &mft
}

// MValidationLookup ...
func (cm AsStandard) MValidationLookup(resource string) *metadata.MValidationLookup {
	mvl := metadata.MValidationLookup{}
	mi := metadata.MIValidationLookup
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mvl)
		for _, entry := range cd.Entries() {
			tmp := metadata.ValidationLookup{}
			entry.SetFields(&tmp)
			id := mi.ID(&tmp)
			tmp.MValidationLookupType = cm.MValidationLookupType(resource, id)
			mvl.ValidationLookup = append(mvl.ValidationLookup, tmp)
		}
	}
	return &mvl
}

// MValidationLookupType ...
func (cm AsStandard) MValidationLookupType(resource, validationLookup string) *metadata.MValidationLookupType {
	mvlt := metadata.MValidationLookupType{}
	mi := metadata.MIValidationLookupType
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["ValidationLookup"] != validationLookup {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mvlt)
		for _, entry := range cd.Entries() {
			tmp := metadata.ValidationLookupType{}
			entry.SetFields(&tmp)
			mvlt.ValidationLookupType = append(mvlt.ValidationLookupType, tmp)
		}
	}
	return &mvlt
}

// MValidationExternal ...
func (cm AsStandard) MValidationExternal(resource string) *metadata.MValidationExternal {
	mve := metadata.MValidationExternal{}
	mi := metadata.MIValidationExternal
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mve)
		for _, entry := range cd.Entries() {
			tmp := metadata.ValidationExternal{}
			entry.SetFields(&tmp)
			id := mi.ID(&tmp)
			tmp.MValidationExternalType = cm.MValidationExternalType(resource, id)
			mve.ValidationExternal = append(mve.ValidationExternal, tmp)
		}
	}
	return &mve
}

// MValidationExternalType ...
func (cm AsStandard) MValidationExternalType(resource, validationExternal string) *metadata.MValidationExternalType {
	mvet := metadata.MValidationExternalType{}
	mi := metadata.MIValidationExternalType
	for _, cd := range cm.Elements[mi.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["ValidationExternal"] != validationExternal {
			continue
		}
		metadata.FieldTransfer(cd.Attr).To(&mvet)
		for _, entry := range cd.Entries() {
			tmp := metadata.ValidationExternalType{}
			entry.SetFields(&tmp)
			mvet.ValidationExternalType = append(mvet.ValidationExternalType, tmp)
		}
	}
	return &mvet
}
