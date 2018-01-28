package kraken

import (
	"net/http"
	"testing"
)

func TestGetClip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		slug       string
		respBody   string
	}{
		{
			http.StatusNotFound,
			"slug-that-does-not-exist",
			`{"error":"Not Found","status":404,"message":"Clip does not exist"}`,
		},
		{
			http.StatusOK,
			"AwkwardHelplessSalamanderSwiftRage",
			`{"slug":"AwkwardHelplessSalamanderSwiftRage","tracking_id":"157589949","url":"https://clips.twitch.tv/AwkwardHelplessSalamanderSwiftRage?tt_medium=clips_api\u0026tt_content=url","embed_url":"https://clips.twitch.tv/embed?clip=AwkwardHelplessSalamanderSwiftRage\u0026tt_medium=clips_api\u0026tt_content=embed","embed_html":"\u003ciframe src='https://clips.twitch.tv/embed?clip=AwkwardHelplessSalamanderSwiftRage\u0026tt_medium=clips_api\u0026tt_content=embed' width='640' height='360' frameborder='0' scrolling='no' allowfullscreen='true'\u003e\u003c/iframe\u003e","broadcaster":{"id":"67955580","name":"chewiemelodies","display_name":"ChewieMelodies","channel_url":"https://www.twitch.tv/chewiemelodies","logo":"https://static-cdn.jtvnw.net/jtv_user_pictures/chewiemelodies-profile_image-2e631188c0919167-150x150.jpeg"},"curator":{"id":"53834192","name":"blacknova03","display_name":"BlackNova03","channel_url":"https://www.twitch.tv/blacknova03","logo":"https://static-cdn.jtvnw.net/jtv_user_pictures/3328691c42d67df4-profile_image-150x150.jpeg"},"vod":null,"broadcast_id":"26860897456","game":"Creative","language":"en","title":"babymetal","views":114,"duration":60,"created_at":"2017-11-30T22:34:18Z","thumbnails":{"medium":"https://clips-media-assets.twitch.tv/157589949-preview-480x272.jpg","small":"https://clips-media-assets.twitch.tv/157589949-preview-260x147.jpg","tiny":"https://clips-media-assets.twitch.tv/157589949-preview-86x45.jpg"}}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient("cid", newMockHandler(testCase.statusCode, testCase.respBody))
		resp, err := c.GetClip(testCase.slug)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusNotFound {
			if resp.Error != "Not Found" {
				t.Errorf("expected error to be %s, got %s", "Not Found", resp.Error)
			}

			if resp.StatusCode != http.StatusNotFound {
				t.Errorf("expected error code to be %d, got %d", http.StatusNotFound, resp.StatusCode)
			}

			if resp.ErrorMessage != "Clip does not exist" {
				t.Errorf("expected error message to be %s, got %s", "Clip does not exist", resp.ErrorMessage)
			}
		}

		if resp.StatusCode == http.StatusOK && resp.Data.Slug != testCase.slug {
			t.Errorf("expected slug to be %s, got %s", testCase.slug, resp.Data.Slug)
		}
	}
}

func TestGetTopClip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode int
		channel    string
		slug       string
		respBody   string
	}{
		{
			http.StatusOK,
			"summit1g",
			"EncouragingPluckySlothSSSsss",
			`{"clips":[{"slug":"EncouragingPluckySlothSSSsss","tracking_id":"182509178","url":"https://clips.twitch.tv/EncouragingPluckySlothSSSsss?tt_medium=clips_api\u0026tt_content=url","embed_url":"https://clips.twitch.tv/embed?clip=EncouragingPluckySlothSSSsss\u0026tt_medium=clips_api\u0026tt_content=embed","embed_html":"\u003ciframe src='https://clips.twitch.tv/embed?clip=EncouragingPluckySlothSSSsss\u0026tt_medium=clips_api\u0026tt_content=embed' width='640' height='360' frameborder='0' scrolling='no' allowfullscreen='true'\u003e\u003c/iframe\u003e","broadcaster":{"id":"26490481","name":"summit1g","display_name":"summit1g","channel_url":"https://www.twitch.tv/summit1g","logo":"https://static-cdn.jtvnw.net/jtv_user_pictures/200cea12142f2384-profile_image-150x150.png"},"curator":{"id":"143839181","name":"nb00ts","display_name":"nB00ts","channel_url":"https://www.twitch.tv/nb00ts","logo":"https://static-cdn.jtvnw.net/jtv_user_pictures/nb00ts-profile_image-b1996340ee6a7ae4-150x150.jpeg"},"vod":{"id":"222004532","url":"https://www.twitch.tv/videos/222004532?t=5h7m50s","offset":18470,"preview_image_url":"https://static-cdn.jtvnw.net/s3_vods/fb22aab1431d75543ea8_summit1g_27381519904_780015163/thumb/thumb0-320x240.jpg"},"broadcast_id":"27381519904","game":"Sea of Thieves","language":"en","title":"summit and fat tim discover how to usemaps","views":78322,"duration":59.33,"created_at":"2018-01-25T04:04:15Z","thumbnails":{"medium":"https://clips-media-assets.twitch.tv/182509178-preview-480x272.jpg","small":"https://clips-media-assets.twitch.tv/182509178-preview-260x147.jpg","tiny":"https://clips-media-assets.twitch.tv/182509178-preview-86x45.jpg"}}],"_cursor":"MQ=="}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient("cid", newMockHandler(testCase.statusCode, testCase.respBody))

		resp, err := c.GetTopClips(&TopClipsParams{
			Channel: testCase.channel,
			Limit:   1,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code %d, but got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.Data.Clips[0].Slug != testCase.slug {
			t.Errorf("expected slug to be %s, got %s", testCase.slug, resp.Data.Clips[0].Slug)
		}

		if resp.Data.Clips[0].Broadcaster.Name != testCase.channel {
			t.Errorf("expected channel to be %s, got %s", testCase.channel, resp.Data.Clips[0].Broadcaster.Name)
		}
	}
}

func TestGetFollowedClip(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		statusCode  int
		accessToken string
		slug        string
		respBody    string
	}{
		{
			http.StatusUnauthorized,
			"",
			"",
			`{"error":"Unauthorized","status":401,"message":"invalid oauth token"}`,
		},
		{
			http.StatusOK,
			"my-valid-access-token",
			"PlayfulTemperedFloofGivePLZ",
			`{"clips":[{"slug":"PlayfulTemperedFloofGivePLZ","tracking_id":"93288497","url":"https://clips.twitch.tv/PlayfulTemperedFloofGivePLZ?tt_medium=clips_api\u0026tt_content=url","embed_url":"https://clips.twitch.tv/embed?clip=PlayfulTemperedFloofGivePLZ\u0026tt_medium=clips_api\u0026tt_content=embed","embed_html":"\u003ciframe src='https://clips.twitch.tv/embed?clip=PlayfulTemperedFloofGivePLZ\u0026tt_medium=clips_api\u0026tt_content=embed' width='640' height='360' frameborder='0' scrolling='no' allowfullscreen='true'\u003e\u003c/iframe\u003e","broadcaster":{"id":"71166086","name":"deadmau5","display_name":"deadmau5","channel_url":"https://www.twitch.tv/deadmau5","logo":"https://static-cdn.jtvnw.net/jtv_user_pictures/deadmau5-profile_image-ee72d3d05d3b99a8-150x150.jpeg"},"curator":{"id":"23844496","name":"bloodtake","display_name":"Bloodtake","channel_url":"https://www.twitch.tv/bloodtake","logo":"https://static-cdn.jtvnw.net/jtv_user_pictures/b8665ac6-6069-4b40-b459-1c50e72dcdab-profile_image-150x150.png"},"vod":{"id":"155025475","url":"https://www.twitch.tv/videos/155025475?t=5h25m31s","offset":19531,"preview_image_url":"https://static-cdn.jtvnw.net/s3_vods/4aa296a0007691e1c23c_deadmau5_25616507600_669702217/thumb/thumb0-320x240.jpg"},"broadcast_id":"25616507600","game":"Creative","language":"en","title":"Deadmau5 Finds Sandstorm","views":419449,"duration":29.75,"created_at":"2017-06-28T20:15:23Z","thumbnails":{"medium":"https://clips-media-assets.twitch.tv/25616507600-offset-19556-preview-480x272.jpg","small":"https://clips-media-assets.twitch.tv/25616507600-offset-19556-preview-260x147.jpg","tiny":"https://clips-media-assets.twitch.tv/25616507600-offset-19556-preview-86x45.jpg"}}],"_cursor":"MQ=="}`,
		},
	}

	for _, testCase := range testCases {
		c := newMockClient("cid", newMockHandler(testCase.statusCode, testCase.respBody))
		c.SetAccessToken(testCase.accessToken)

		resp, err := c.GetFollowedClips(&FollowedClipsParams{
			Limit: 1,
		})
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != testCase.statusCode {
			t.Errorf("expected status code %d, but got %d", testCase.statusCode, resp.StatusCode)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if resp.Error != "Unauthorized" {
				t.Errorf("expected error to be %s, got %s", "Unauthorized", resp.Error)
			}

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected error code to be %d, got %d", http.StatusNotFound, resp.StatusCode)
			}

			if resp.ErrorMessage != "invalid oauth token" {
				t.Errorf("expected error message to be %s, got %s", "invalid oauth token", resp.ErrorMessage)
			}
		}

		if resp.StatusCode == http.StatusOK && resp.Data.Clips[0].Slug != testCase.slug {
			t.Errorf("expected slug to be %s, got %s", testCase.slug, resp.Data.Clips[0].Slug)
		}
	}
}
