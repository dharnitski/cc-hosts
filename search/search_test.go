package search_test

import (
	"fmt"
	"testing"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/search"
	"github.com/dharnitski/cc-hosts/vertices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearcher_GetTargets(t *testing.T) {
	t.Parallel()

	eOffsets := edges.Offsets{}
	err := eOffsets.Load(fmt.Sprintf("../data/%s", access.EdgesOffsetsFile))
	require.NoError(t, err)
	out := edges.NewEdges("../data/edges", eOffsets)

	eOffsets = edges.Offsets{}
	err = eOffsets.Load(fmt.Sprintf("../data/%s", access.EdgesReversedOffsetFile))
	require.NoError(t, err)
	in := edges.NewEdges("../data/edges_reversed", eOffsets)

	vOffsets := vertices.Offsets{}
	err = vOffsets.Load(fmt.Sprintf("../data/%s", access.VerticesOffsetsFile))
	require.NoError(t, err)
	v := vertices.NewVertices("../data/vertices", vOffsets)

	searcher := search.NewSearcher(v, out, in)
	results, err := searcher.GetTargets(t.Context(), "binaryedge.io")
	assert.NoError(t, err)
	assert.Equal(t, []string{"40fy.io", "app.binaryedge.io", "blog.binaryedge.io", "cloudflare.com", "coalitioninc.com", "cyberfables.io", "d1ehrggk1349y0.cloudfront.net", "facebook.com", "fonts.googleapis.com", "github.com", "linkedin.com", "maps.googleapis.com", "slack.binaryedge.io", "support.cloudflare.com", "twitter.com"}, results.Out)
	assert.Equal(t, []string{}, results.In)
	assert.Equal(t, "binaryedge.io", results.Target)
}

var socialNetworkKeywords = []string{
	"2mdn.net",
	"3gppnetwork.org",
	"a2z.com",
	"aaplimg.com",
	"akadns.net",
	"akamai.net",
	"akamaiedge.net",
	"amazon-adsystem.com",
	"amazon.com",
	"amazonaws.com",
	"ampproject.org",
	"android.com",
	"app-analytics-services.com",
	"app-measurement.com",
	"apple-dns.net",
	"apple.com",
	"applovin.com",
	"appsflyersdk.com",
	"atlassian.net",
	"azure.com",
	"baidu.com",
	"bing.com",
	"bsky.app",
	"bytefcdn-oversea.com",
	"bytefcdn-ttpeu.com",
	"casalemedia.com",
	"cdn-apple.com",
	"cdn77.org",
	"cdninstagram.com",
	"cloudflare-dns.com",
	"cloudflare.com",
	"cloudforce.com",
	"cloudfront.net",
	"crashlytics.com",
	"criteo.com",
	"crunchbase.com",
	"digicert.com",
	"dns.google",
	"documentforce.com",
	"doubleclick.net",
	"doubleverify.com",
	"facebook.com",
	"fastly.net",
	"fbcdn.net",
	"force.com",
	"ggpht.com",
	"gmail.com",
	"goo.gl",
	"google-analytics.com",
	"google.com",
	"googleadservices.com",
	"googleapis.com",
	"googlesyndication.com",
	"googletagmanager.com",
	"googleusercontent.com",
	"googlevideo.com",
	"gstatic.com",
	"gvt1.com",
	"gvt2.com",
	"icims.com",
	"icloud.com",
	"instagram.com",
	"lencr.org",
	"linkedin.com",
	"live.com",
	"mailchi.mp",
	"mailchimp.com",
	"microsoft.com",
	"microsoftonline.com",
	"miui.com",
	"msftncsi.com",
	"msn.com",
	"mzstatic.com",
	"netflix.com",
	"ntp.org",
	"office.com",
	"office365.com",
	"one.one",
	"pangle.io",
	"pinterest.com",
	"qlivecdn.com",
	"qq.com",
	"rbxcdn.com",
	"roblox.com",
	"rocket-cdn.com",
	"root-servers.net",
	"salesforce.com",
	"samsung.com",
	"sentry.io",
	"sharepoint.com",
	"shopify.com",
	"skype.com",
	"snapchat.com",
	"spotify.com",
	"taboola.com",
	"threads.net",
	"tiktok.com",
	"tiktokcdn-eu.com",
	"tiktokcdn-us.com",
	"tiktokcdn.com",
	"tiktokv.com",
	"trafficmanager.net",
	"ttlivecdn.com",
	"twitter.com",
	"ui.com",
	"unity3d.com",
	"visualforce.com",
	"vungle.com",
	"whatsapp.net",
	"wikipedia.org",
	"windows.com",
	"windows.net",
	"windowsupdate.com",
	"wordpress.com",
	"wordpress.org",
	"wp.com",
	"wpengine.com",
	"x.com",
	"xiaomi.com",
	"yahoo.com",
	"yelp-ir.com",
	"yelp-support.com",
	"yelp.ca",
	"yelp.com",
	"youtube.com",
	"ytimg.com",
	"zendesk.com",
}
