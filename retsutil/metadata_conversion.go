package retsutil

import (
	"reflect"
	"strings"

	"github.com/jpfielding/gorets/metadata"
	"github.com/jpfielding/gorets/rets"
)

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
	cm.SetForeignKeys(&ms.System)
	cm.SetResources(&ms.System)
	cm.SetFilters(&ms.System)
	return ms, nil
}

// SetForeignKeys ...
func (cm AsStandard) SetForeignKeys(sys *metadata.System) {
	for _, fk := range cm.Elements[metadata.MetaForeignKey.Name] {
		sys.MForeignKey = &metadata.MForeignKey{}
		metadata.FieldTransfer(fk.Attr).To(sys.MForeignKey)
		for _, entry := range fk.Entries() {
			tmp := &metadata.ForeignKey{}
			entry.SetFields(tmp)
			sys.MForeignKey.ForeignKey = append(sys.MForeignKey.ForeignKey, *tmp)
		}
	}
}

// SetFilters ...
func (cm AsStandard) SetFilters(sys *metadata.System) {
	for _, cd := range cm.Elements[metadata.MetaFilter.Name] {
		sys.MFilter = &metadata.MFilter{}
		metadata.FieldTransfer(cd.Attr).To(sys.MFilter)
		for _, entry := range cd.Entries() {
			tmp := &metadata.Filter{}
			entry.SetFields(tmp)
			cm.SetFilterTypes(tmp)
			sys.MFilter.Filter = append(sys.MFilter.Filter, *tmp)
		}
	}
}

// SetResources ...
func (cm AsStandard) SetResources(sys *metadata.System) {
	for _, cd := range cm.Elements[metadata.MetaResource.Name] {
		sys.MResource = &metadata.MResource{}
		metadata.FieldTransfer(cd.Attr).To(sys.MResource)
		for _, entry := range cd.Entries() {
			tmp := &metadata.Resource{}
			entry.SetFields(tmp)
			cm.SetClasses(tmp)
			cm.SetObjects(tmp)
			cm.SetLookups(tmp)
			cm.SetSearchHelps(tmp)
			cm.SetEditMasks(tmp)
			cm.SetUpdateHelps(tmp)
			cm.SetValidationExpressions(tmp)
			cm.SetValidationExternals(tmp)
			cm.SetValidationLookups(tmp)
			sys.MResource.Resource = append(sys.MResource.Resource, *tmp)
		}
	}
}

// SetValidationExpressions ...
func (cm AsStandard) SetValidationExpressions(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaValidationExpression.Name] {
		resourceID := string(resource.ResourceID)
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		resource.MValidationExpression = &metadata.MValidationExpression{}
		metadata.FieldTransfer(cd.Attr).To(resource.MValidationExpression)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ValidationExpression{}
			entry.SetFields(tmp)
			resource.MValidationExpression.ValidationExpression = append(resource.MValidationExpression.ValidationExpression, *tmp)
		}
	}
}

// SetUpdateHelps ...
func (cm AsStandard) SetUpdateHelps(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaUpdateHelp.Name] {
		resourceID := string(resource.ResourceID)
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		resource.MUpdateHelp = &metadata.MUpdateHelp{}
		metadata.FieldTransfer(cd.Attr).To(resource.MUpdateHelp)
		for _, entry := range cd.Entries() {
			tmp := &metadata.UpdateHelp{}
			entry.SetFields(tmp)
			resource.MUpdateHelp.UpdateHelp = append(resource.MUpdateHelp.UpdateHelp, *tmp)
		}
	}
}

// SetEditMasks ...
func (cm AsStandard) SetEditMasks(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaEditMask.Name] {
		resourceID := string(resource.ResourceID)
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		resource.MEditMask = &metadata.MEditMask{}
		metadata.FieldTransfer(cd.Attr).To(resource.MEditMask)
		for _, entry := range cd.Entries() {
			tmp := &metadata.EditMask{}
			entry.SetFields(tmp)
			resource.MEditMask.EditMask = append(resource.MEditMask.EditMask, *tmp)
		}
	}
}

// SetSearchHelps ...
func (cm AsStandard) SetSearchHelps(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaSearchHelp.Name] {
		resourceID := string(resource.ResourceID)
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		resource.MSearchHelp = &metadata.MSearchHelp{}
		metadata.FieldTransfer(cd.Attr).To(resource.MSearchHelp)
		for _, entry := range cd.Entries() {
			tmp := &metadata.SearchHelp{}
			entry.SetFields(tmp)
			resource.MSearchHelp.SearchHelp = append(resource.MSearchHelp.SearchHelp, *tmp)
		}
	}
}

// SetClasses ...
func (cm AsStandard) SetClasses(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaClass.Name] {
		resourceID := string(resource.ResourceID)
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		resource.MClass = &metadata.MClass{}
		metadata.FieldTransfer(cd.Attr).To(resource.MClass)
		for _, entry := range cd.Entries() {
			tmp := &metadata.Class{}
			entry.SetFields(tmp)
			cm.SetTables(resourceID, tmp)
			cm.SetUpdates(resourceID, tmp)
			cm.SetColumnGroupSets(resourceID, tmp)
			cm.SetColumnGroups(string(resource.ResourceID), tmp)
			resource.MClass.Class = append(resource.MClass.Class, *tmp)
		}
	}
}

// SetColumnGroups ...
func (cm AsStandard) SetColumnGroups(resourceID string, class *metadata.Class) {
	for _, cd := range cm.Elements[metadata.MetaColumnGroup.Name] {
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		classID := string(class.ClassName)
		if cd.Attr["Class"] != classID {
			continue
		}
		class.MColumnGroup = &metadata.MColumnGroup{}
		metadata.FieldTransfer(cd.Attr).To(class.MColumnGroup)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ColumnGroup{}
			entry.SetFields(tmp)
			cm.SetColumnGroupControls(resourceID, classID, tmp)
			cm.SetColumnGroupTables(resourceID, classID, tmp)
			cm.SetColumnGroupNormalizations(resourceID, classID, tmp)
			class.MColumnGroup.ColumnGroup = append(class.MColumnGroup.ColumnGroup, *tmp)
		}
	}
}

// SetColumnGroupControls ...
func (cm AsStandard) SetColumnGroupControls(resourceID, classID string, cg *metadata.ColumnGroup) {
	for _, cd := range cm.Elements[metadata.MetaColumnGroupControl.Name] {
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		if cd.Attr["Class"] != classID {
			continue
		}
		if cd.Attr["ColumnGroup"] != string(cg.ColumnGroupName) {
			continue
		}
		cg.MColumnGroupControl = &metadata.MColumnGroupControl{}
		metadata.FieldTransfer(cd.Attr).To(cg.MColumnGroupControl)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ColumnGroupControl{}
			entry.SetFields(tmp)
			cg.MColumnGroupControl.ColumnGroupControl = append(cg.MColumnGroupControl.ColumnGroupControl, *tmp)
		}
	}
}

// SetColumnGroupTables ...
func (cm AsStandard) SetColumnGroupTables(resourceID, classID string, cg *metadata.ColumnGroup) {
	for _, cd := range cm.Elements[metadata.MetaColumnGroupTable.Name] {
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		if cd.Attr["Class"] != classID {
			continue
		}
		if cd.Attr["ColumnGroup"] != string(cg.ColumnGroupName) {
			continue
		}
		cg.MColumnGroupTable = &metadata.MColumnGroupTable{}
		metadata.FieldTransfer(cd.Attr).To(cg.MColumnGroupTable)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ColumnGroupTable{}
			entry.SetFields(tmp)
			cg.MColumnGroupTable.ColumnGroupTable = append(cg.MColumnGroupTable.ColumnGroupTable, *tmp)
		}
	}
}

// SetColumnGroupNormalizations ...
func (cm AsStandard) SetColumnGroupNormalizations(resourceID, classID string, cg *metadata.ColumnGroup) {
	for _, cd := range cm.Elements[metadata.MetaColumnGroupNormalization.Name] {
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		if cd.Attr["Class"] != classID {
			continue
		}
		if cd.Attr["ColumnGroup"] != string(cg.ColumnGroupName) {
			continue
		}
		cg.MColumnGroupNormalization = &metadata.MColumnGroupNormalization{}
		metadata.FieldTransfer(cd.Attr).To(cg.MColumnGroupTable)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ColumnGroupNormalization{}
			entry.SetFields(tmp)
			cg.MColumnGroupNormalization.ColumnGroupNormalization = append(cg.MColumnGroupNormalization.ColumnGroupNormalization, *tmp)
		}
	}
}

// SetColumnGroupSets ...
func (cm AsStandard) SetColumnGroupSets(resource string, class *metadata.Class) {
	for _, cd := range cm.Elements[metadata.MetaColumnGroupSet.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != string(class.ClassName) {
			continue
		}
		class.MColumnGroupSet = &metadata.MColumnGroupSet{}
		metadata.FieldTransfer(cd.Attr).To(class.MColumnGroupSet)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ColumnGroupSet{}
			entry.SetFields(tmp)
			class.MColumnGroupSet.ColumnGroupSet = append(class.MColumnGroupSet.ColumnGroupSet, *tmp)
		}
	}
}

// SetUpdates ...
func (cm AsStandard) SetUpdates(resource string, class *metadata.Class) {
	for _, cd := range cm.Elements[metadata.MetaUpdate.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != string(class.ClassName) {
			continue
		}
		class.MUpdate = &metadata.MUpdate{}
		metadata.FieldTransfer(cd.Attr).To(class.MUpdate)
		for _, entry := range cd.Entries() {
			tmp := &metadata.Update{}
			entry.SetFields(tmp)
			classID := string(class.ClassName)
			cm.SetUpdateTypes(resource, classID, tmp)
			class.MUpdate.Update = append(class.MUpdate.Update, *tmp)
		}
	}
}

// SetUpdateTypes ...
func (cm AsStandard) SetUpdateTypes(resource, class string, update *metadata.Update) {
	for _, cd := range cm.Elements[metadata.MetaUpdateType.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != class {
			continue
		}
		if cd.Attr["Update"] != string(update.MetadataEntryID) {
			continue
		}
		update.MUpdateType = &metadata.MUpdateType{}
		metadata.FieldTransfer(cd.Attr).To(update.MUpdateType)
		for _, entry := range cd.Entries() {
			tmp := &metadata.UpdateType{}
			entry.SetFields(tmp)
			update.MUpdateType.UpdateType = append(update.MUpdateType.UpdateType, *tmp)
		}
	}
}

// SetTables ...
func (cm AsStandard) SetTables(resource string, class *metadata.Class) {
	for _, cd := range cm.Elements[metadata.MetaTable.Name] {
		if cd.Attr["Resource"] != resource {
			continue
		}
		if cd.Attr["Class"] != string(class.ClassName) {
			continue
		}
		class.MTable = &metadata.MTable{}
		metadata.FieldTransfer(cd.Attr).To(class.MTable)
		for _, entry := range cd.Entries() {
			tmp := &metadata.Field{}
			entry.SetFields(tmp)
			class.MTable.Field = append(class.MTable.Field, *tmp)
		}
	}
}

// SetObjects ...
func (cm AsStandard) SetObjects(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaObject.Name] {
		if cd.Attr["Resource"] != string(resource.ResourceID) {
			continue
		}
		resource.MObject = &metadata.MObject{}
		metadata.FieldTransfer(cd.Attr).To(resource.MObject)
		for _, entry := range cd.Entries() {
			tmp := &metadata.Object{}
			entry.SetFields(tmp)
			resource.MObject.Object = append(resource.MObject.Object, *tmp)
		}
	}
}

// SetLookups ...
func (cm AsStandard) SetLookups(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaLookup.Name] {
		if cd.Attr["Resource"] != string(resource.ResourceID) {
			continue
		}
		resource.MLookup = &metadata.MLookup{}
		metadata.FieldTransfer(cd.Attr).To(resource.MLookup)
		for _, entry := range cd.Entries() {
			tmp := &metadata.Lookup{}
			entry.SetFields(tmp)
			cm.SetLookupTypes(tmp)
			resource.MLookup.Lookup = append(resource.MLookup.Lookup, *tmp)
		}
	}
}

// SetLookupTypes ...
func (cm AsStandard) SetLookupTypes(lookup *metadata.Lookup) {
	for _, cd := range cm.Elements[metadata.MetaLookupType.Name] {
		if cd.Attr["Lookup"] != string(lookup.LookupName) {
			continue
		}
		lookup.MLookupType = &metadata.MLookupType{}
		metadata.FieldTransfer(cd.Attr).To(lookup.MLookupType)
		for _, entry := range cd.Entries() {
			tmp := &metadata.LookupType{}
			entry.SetFields(tmp)
			lookup.MLookupType.LookupType = append(lookup.MLookupType.LookupType, *tmp)
		}
	}
}

// SetFilterTypes ...
func (cm AsStandard) SetFilterTypes(filter *metadata.Filter) {
	for _, cd := range cm.Elements[metadata.MetaFilterType.Name] {
		if cd.Attr["Filter"] != string(filter.FilterID) {
			continue
		}
		filter.MFilterType = &metadata.MFilterType{}
		metadata.FieldTransfer(cd.Attr).To(filter.MFilterType)
		for _, entry := range cd.Entries() {
			tmp := &metadata.FilterType{}
			entry.SetFields(tmp)
			filter.MFilterType.FilterType = append(filter.MFilterType.FilterType, *tmp)
		}
	}
}

// SetValidationLookups ...
func (cm AsStandard) SetValidationLookups(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaValidationLookup.Name] {
		resourceID := string(resource.ResourceID)
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		resource.MValidationLookup = &metadata.MValidationLookup{}
		metadata.FieldTransfer(cd.Attr).To(resource.MValidationLookup)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ValidationLookup{}
			entry.SetFields(tmp)
			cm.SetValidationLookupTypes(resourceID, tmp)
			resource.MValidationLookup.ValidationLookup = append(resource.MValidationLookup.ValidationLookup, *tmp)
		}
	}
}

// SetValidationLookupTypes ...
func (cm AsStandard) SetValidationLookupTypes(resourceID string, lookup *metadata.ValidationLookup) {
	for _, cd := range cm.Elements[metadata.MetaValidationLookupType.Name] {
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		if cd.Attr["Lookup"] != string(lookup.ValidationLookupName) {
			continue
		}
		lookup.MValidationLookupType = &metadata.MValidationLookupType{}
		metadata.FieldTransfer(cd.Attr).To(lookup.MValidationLookupType)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ValidationLookupType{}
			entry.SetFields(tmp)
			lookup.MValidationLookupType.ValidationLookupType = append(lookup.MValidationLookupType.ValidationLookupType, *tmp)
		}
	}
}

// SetValidationExternals ...
func (cm AsStandard) SetValidationExternals(resource *metadata.Resource) {
	for _, cd := range cm.Elements[metadata.MetaValidationExternal.Name] {
		resourceID := string(resource.ResourceID)
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		resource.MValidationExternal = &metadata.MValidationExternal{}
		metadata.FieldTransfer(cd.Attr).To(resource.MValidationExternal)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ValidationExternal{}
			entry.SetFields(tmp)
			cm.SetValidationExternalTypes(resourceID, tmp)
			resource.MValidationExternal.ValidationExternal = append(resource.MValidationExternal.ValidationExternal, *tmp)
		}
	}
}

// SetValidationExternalTypes ...
func (cm AsStandard) SetValidationExternalTypes(resourceID string, validation *metadata.ValidationExternal) {
	for _, cd := range cm.Elements[metadata.MetaValidationExternalType.Name] {
		if cd.Attr["Resource"] != resourceID {
			continue
		}
		if cd.Attr["ValidationExternal"] != string(validation.ValidationExternalName) {
			continue
		}
		validation.MValidationExternalType = &metadata.MValidationExternalType{}
		metadata.FieldTransfer(cd.Attr).To(validation.MValidationExternalType)
		for _, entry := range cd.Entries() {
			tmp := &metadata.ValidationExternalType{}
			entry.SetFields(tmp)
			validation.MValidationExternalType.ValidationExternalType = append(validation.MValidationExternalType.ValidationExternalType, *tmp)
		}
	}
}

func (cm AsStandard) setFields(foo interface{}, fields map[string]string) {
	for k, v := range fields {
		val := reflect.ValueOf(foo).Elem().FieldByNameFunc(func(n string) bool {
			return strings.ToLower(n) == strings.ToLower(k)
		})
		if val.IsValid() {
			val.SetString(v)
		}
	}
}
