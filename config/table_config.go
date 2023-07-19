package config

type TableConfig struct {
	shorttermduration       uint16
	longtermduration        uint16
	initvacuum              bool
	autovacuum              bool
	tablespacename          string
	tablespace              bool
	tablespacepath          string
	shorttermtablespacename string
	shorttermtablespace     bool
	shorttermtablespacepath string
	longtermtablespacename  string
	longtermtablespace      bool
	longtermtablespacepath  string
}

// GetShorttermDuration is a function that returns the short-term duration.
func (tableConf *TableConfig) GetShorttermDuration() uint16 {
	return tableConf.shorttermduration
}

// GetLongtermDuration is a function that returns the long-term duration.
func (tableConf *TableConfig) GetLongtermDuration() uint16 {
	return tableConf.longtermduration
}

// IsInitvacuum is a function that returns whether to initialize vacuum.
func (tableConf *TableConfig) IsInitvacuum() bool {
	return tableConf.initvacuum
}

// IsAutovacuum is a function that returns whether to use auto vacuum.
func (tableConf *TableConfig) IsAutovacuum() bool {
	return tableConf.autovacuum
}

// GetTableSpaceName is a function that returns the name of the tablespace.
func (tableConf *TableConfig) GetTableSpaceName() string {
	return tableConf.tablespacename
}

// GetTableSpace is a function that returns whether to use tablespace.
func (tableConf *TableConfig) GetTableSpaceCount() bool {
	return tableConf.tablespace
}

// GetTableSpacePath is a function that returns the path of the tablespace.
func (tableConf *TableConfig) GetTableSpacePath() string {
	return tableConf.tablespacepath
}

// GetShorttermTableSpaceName is a function that returns the name of the short-term tablespace.
func (tableConf *TableConfig) GetShorttermTableSpaceName() string {
	return tableConf.shorttermtablespacename
}

// GetShorttermTableSpace is a function that returns whether to use short-term tablespace.
func (tableConf *TableConfig) GetShorttermTableSpace() bool {
	return tableConf.shorttermtablespace
}

// GetShorttermTableSpacePath is a function that returns the path of the short-term tablespace.
func (tableConf *TableConfig) GetShorttermTableSpacePath() string {
	return tableConf.shorttermtablespacepath
}

// GetLongtermTableSpaceName is a function that returns the name of the long-term tablespace.
func (tableConf *TableConfig) GetLongtermTableSpaceName() string {
	return tableConf.longtermtablespacename
}

// GetLongtermTableSpace is a function that returns whether to use long-term tablespace.
func (tableConf *TableConfig) GetLongtermTableSpace() bool {
	return tableConf.longtermtablespace
}

// GetLongtermTableSpacePath is a function that returns the path of the long-term tablespace.
func (tableConf *TableConfig) GetLongtermTableSpacePath() string {
	return tableConf.longtermtablespacepath
}

// IsLongtermTablespace is a function that returns whether to use long-term tablespace.
func (tableConf *TableConfig) IsLongtermTablespace() bool {
	return tableConf.longtermtablespace
}
