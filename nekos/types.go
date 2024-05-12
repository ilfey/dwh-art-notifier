package nekos

type Response struct {
	Results []struct {
		ArtistHref *string `json:"artist_href"`
		ArtistName *string `json:"artist_name"`
		SourceURL  *string `json:"source_url"`
		AnimeName  *string `json:"anime_name"`
		URL        *string `json:"url"`
	} `json:"results"`
}
