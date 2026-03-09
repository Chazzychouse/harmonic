package tags

type TrackMeta struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	AlbumArtist string `json:"album_artist"`
	Genre       string `json:"genre"`
	Year        int    `json:"year"`
	TrackNum    int    `json:"track_num"`
	TrackTotal  int    `json:"track_total"`
	DiscNum     int    `json:"disc_num"`
	DiscTotal   int    `json:"disc_total"`
	HasArt      bool   `json:"has_art"`
}
