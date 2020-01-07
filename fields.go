package transmissionRPC

type GetField int

const (
	ActivityDate GetField = iota
	AddedDate
	BandwidthPriority
	Comment
	CorruptEver
	Creator
	DateCreated
	DesiredAvailable
	DoneDate
	DownloadDir
	DownloadedEver
	DownloadLimit
	DownloadLimited
	EditDate
	Error
	ErrorString
	Eta
	EtaIdle
	Files
	FileStats
	HashString
	HaveUnchecked
	HaveValid
	HonorsSessionLimits
	ID
	IsFinished
	IsPrivate
	IsStalled
	Labels
	LeftUntilDone
	MagnetLink
	ManualAnnounceTime
	MaxConnectedPeers
	MetadataPercentComplete
	Name
	PeerLimit
	Peers
	PeersConnected
	PeersFrom
	PeersGettingFromUs
	PeersSendingToUs
	PercentDone
	Pieces
	PieceCount
	PieceSize
	Priorities
	QueuePosition
	RateDownload
	RateUpload
	RecheckProgress
	SecondsDownloading
	SecondsSeeding
	SeedIdleLimit
	SeedIdleMode
	SeedRatioLimit
	SeedRatioMode
	SizeWhenDone
	StartDate
	Status
	Trackers
	TrackerStats
	TotalSize
	TorrentFile
	UploadedEver
	UploadLimit
	UploadLimited
	UploadRatio
	Wanted
	Webseeds
	WebseedsSendingToUs
)

func (f GetField) String() string {
	_strings := []string{
		"activityDate",
		"addedDate",
		"bandwidthPriority",
		"comment",
		"corruptEver",
		"creator",
		"dateCreated",
		"desiredAvailable",
		"doneDate",
		"downloadDir",
		"downloadedEver",
		"downloadLimit",
		"downloadLimited",
		"editDate",
		"error",
		"errorString",
		"eta",
		"etaIdle",
		"files",
		"fileStats",
		"hashString",
		"haveUnchecked",
		"haveValid",
		"honorsSessionLimits",
		"id",
		"isFinished",
		"isPrivate",
		"isStalled",
		"labels",
		"leftUntilDone",
		"magnetLink",
		"manualAnnounceTime",
		"maxConnectedPeers",
		"metadataPercentComplete",
		"name",
		"peer-limit",
		"peers",
		"peersConnected",
		"peersFrom",
		"peersGettingFromUs",
		"peersSendingToUs",
		"percentDone",
		"pieces",
		"pieceCount",
		"pieceSize",
		"priorities",
		"queuePosition",
		"rateDownload",
		"rateUpload",
		"recheckProgress",
		"secondsDownloading",
		"secondsSeeding",
		"seedIdleLimit",
		"seedIdleMode",
		"seedRatioLimit",
		"seedRatioMode",
		"sizeWhenDone",
		"startDate",
		"status",
		"trackers",
		"trackerStats",
		"totalSize",
		"torrentFile",
		"uploadedEver",
		"uploadLimit",
		"uploadLimited",
		"uploadRatio",
		"wanted",
		"webseeds",
		"webseedsSendingToUs",
	}

	if f < ActivityDate || f > WebseedsSendingToUs {
		return ""
	}

	return _strings[f]
}

func FieldList(f ...GetField) []string {
	var tmp []string

	for _, i := range f {
		tmp = append(tmp, i.String())
	}

	return tmp
}

//=====
func TorrentStatus(ID int) string {
	_strings := []string{
		"Torrent is stopped",
		"Queued to check files",
		"Checking files",
		"Queued to download",
		"Downloading",
		"Queued to seed",
		"Seeding",
	}

	return _strings[ID]
}
