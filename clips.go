package kraken

// Clip ...
type Clip struct {
	Slug        string         `json:"slug"`
	TrackingID  string         `json:"tracking_id"`
	URL         string         `json:"url"`
	EmbedURL    string         `json:"embed_url"`
	EmbedHTML   string         `json:"embed_html"`
	Broadcaster ClipUser       `json:"broadcaster"`
	Curator     ClipUser       `json:"curator"`
	VOD         ClipVOD        `json:"vod"`
	Game        string         `json:"game"`
	Language    string         `json:"language"`
	Title       string         `json:"title"`
	Views       int            `json:"views"`
	Duration    float64        `json:"duration"`
	CreatedAt   string         `json:"created_at"`
	Thumbnails  ClipThumbnails `json:"thumbnails"`
}

// ClipUser ...
type ClipUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	ChannelURL  string `json:"channel_url"`
	Logo        string `json:"logo"`
}

// ClipVOD ...
type ClipVOD struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// ClipThumbnails ...
type ClipThumbnails struct {
	Medium string `json:"medium"`
	Small  string `json:"small"`
	Tiny   string `json:"tiny"`
}

// ClipResponse ...
type ClipResponse struct {
	ResponseCommon
	Data Clip
}

// GetClip gets the clip referenced by the slug.
func (c *Client) GetClip(slug string) (*ClipResponse, error) {
	resp, err := c.get("/clips/"+slug, &Clip{}, nil)
	if err != nil {
		return nil, err
	}

	clip := &ClipResponse{}
	clip.Error = resp.Error
	clip.ErrorStatus = resp.ErrorStatus
	clip.ErrorMessage = resp.ErrorMessage
	clip.StatusCode = resp.StatusCode
	clip.Data = *resp.Data.(*Clip)

	return clip, nil
}

// ManyClips ...
type ManyClips struct {
	Clips  []Clip `json:"clips"`
	Cursor string `json:"_cursor"`
}

// ClipsResponse ...
type ClipsResponse struct {
	ResponseCommon
	Data ManyClips
}

// TopClipsParams ...
type TopClipsParams struct {
	Channel  string `query:"channel"`
	Cursor   string `query:"cursor"`
	Game     string `query:"game"`
	Language string `query:"language"`
	Limit    int    `query:"limit,10"`
	Period   string `query:"period,week"`
	Trending bool   `query:"trending"`
}

// GetTopClips gets the top clips which meet a specified set of parameters.
// Note that if both channel and game are specified, game is ignored.
func (c *Client) GetTopClips(params *TopClipsParams) (*ClipsResponse, error) {
	return c.manyClipsRequest("/clips/top", params)
}

// FollowedClipsParams ...
type FollowedClipsParams struct {
	Cursor   string `query:"cursor"`
	Limit    int    `query:"limit,10"`
	Trending bool   `query:"trending"`
}

// GetFollowedClips the top clips for the games followed by a specified user,
// identified by an OAuth token.
func (c *Client) GetFollowedClips(params *FollowedClipsParams) (*ClipsResponse, error) {
	return c.manyClipsRequest("/clips/followed", params)
}

func (c *Client) manyClipsRequest(path string, params interface{}) (*ClipsResponse, error) {
	resp, err := c.get(path, &ManyClips{}, params)
	if err != nil {
		return nil, err
	}

	clips := &ClipsResponse{}
	clips.Error = resp.Error
	clips.ErrorStatus = resp.ErrorStatus
	clips.ErrorMessage = resp.ErrorMessage
	clips.StatusCode = resp.StatusCode
	clips.Data.Clips = resp.Data.(*ManyClips).Clips
	clips.Data.Cursor = resp.Data.(*ManyClips).Cursor

	return clips, nil
}
