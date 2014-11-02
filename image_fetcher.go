// go run boss.go --consumerkey <key> --consumersecret <secret> --appname <appname>
package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mrjones/oauth"
)

type yahooResponse struct {
	BossResponse bossResponse `json:"bossresponse"`
}

type bossResponse struct {
	Images bossImages `json:"images"`
}

type bossImages struct {
	Results []bossImageResult `json:"results"`
}

type bossImageResult struct {
	Url    string `json:"url"`
	Format string `json:"format"`
}

func findUrlsForWord(word string, consumer *oauth.Consumer, accessToken *oauth.AccessToken) []string {
	switch word {
	default:
		return getImageUrls(word, consumer, accessToken)
	case "to":
		return []string{"http://health.businessweekly.com.tw/images/Columns/0030.jpg",
			"https://dpstprenursery.files.wordpress.com/2013/06/numero-2-letras-y-numeros-numeros-pintado-por-johnnathan-97395131.jpg"}
	case "be":
		return []string{"http://medfordchurch.org/wp-content/uploads/2014/10/Be.jpg",
			"https://docs.google.com/a/klarna.com/spreadsheets/d/1b3uemLu1d9XHiubXNP9xcjlraVPllGRzLIm3PQFOckk/edit#gid=568444414",
			"http://www.discotecasdeibiza.net/wp-content/uploads/10_Logo-Be-.jpg"}
	case "was":
		return []string{"http://www.publicart.sg/sites/default/files/images/Was.jpg",
			"http://files.qrz.com/d/py6hd/WAS_MIXED.jpg"}
	case "were":
		return []string{"http://thumb7.shutterstock.com/display_pic_with_logo/890332/133260281/stock-photo-the-word-were-here-written-on-the-sand-as-the-background-133260281.jpg"}
	case "is":
		return []string{"http://www.indiescreen.us/wp-content/uploads/2013/03/IS-Logo-Box.jpg",
			"https://www.google.com/search?q=be&rlz=1C5CHFA_enIL532IL532&es_sm=91&source=lnms&tbm=isch&sa=X&ei=AEBWVO_lEI2raefMgJAF&ved=0CAkQ_AUoAg&biw=1144&bih=558#tbm=isch&q=is+clipart&facrc=_&imgdii=_&imgrc=V8XMPO4HL41bcM%253A%3Bo_pNHhiaq7NczM%3Bhttp%253A%252F%252Flessonpix.com%252Fdrawings%252F12126%252F380x380%252Fis.png%3Bhttp%253A%252F%252Flessonpix.com%252Fpictures%252F12126%252Fis%3B380%3B380"}
	case "are":
		return []string{"http://www.ncarb.org/~/media/Images/Widgets/are.jpg",
			"http://thumb9.shutterstock.com/display_pic_with_logo/609856/168608942/stock-photo-letters-and-symbols-in-fire-letter-r-168608942.jpg",
			"http://www.clipartbest.com/cliparts/9ip/brr/9ipbrrM9T.jpeg"}
	case "the":
		return []string{"https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcRvVQwo8Gg2uje7I3caJR0NDH38Jut95FQnAPIQ0-T4VawWaWhqBA",
			"https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcTVgJmIkoN_nquVc3GSlyhBy-wFEInGBvMNk-Y9LL1PxpoVpcIiQw"}
	case "a":
		return []string{"http://www.fhp.nu/Global/Bio/Rio/Arrangemang/A-lansering/A-st%C3%A4mpel_svart.jpg",
			"http://www.topuniversities.com/sites/default/files/A-level%20grades.jpg",
			"http://aurora-hall.ru/assets/images/about/aone.jpg"}
	case "an":
		return []string{"http://theartofdissent.net/wordpress/wp-content/uploads/2012/08/a-n.TwitterIcon.MAG_1.jpg",
			"https://pbs.twimg.com/profile_images/391958670/A_N_logo.jpg"}
	case "on":
		return []string{"https://c2.staticflickr.com/8/7272/7125837749_7c6ec31879_z.jpg",
			"https://lh3.ggpht.com/cgvbtpBknBOKkGaddXEtwjoVHcZff6ZMM1sKDTv2J7GLP_feCbOZ29_dVM3Ir_CNbg=w300",
			"https://www.uni-muenster.de/imperia/md/images/ziv/layoutbilder/power.jpeg"}
	case "of":
		return []string{"http://i.imgur.com/vdDxTSZ.jpg",
			"http://freightlinersf.com/wp-content/uploads/2014/08/roasted-chicken-clipart-438.jpg",
			"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTLnWzOLYIyXTewP-nT_kUYiqm-6sw49aX_ENS5gfVDLHPCU8gi6g"}
	case "my":
		return []string{"http://fc01.deviantart.net/fs71/i/2010/211/4/a/Little_My__by_swessu.jpg",
			"http://webfoundation.org/wp-content/uploads/2012/12/myworldlogo.jpg",
			"http://www.deafedge.com/image-files/basic-words-in-sign-language-my.jpg",
			"http://images.clipartpanda.com/jitters-clipart-10715678-little-boy-with-your-escholar-itens-going-to-school.jpg"}
	case "your":
		return []string{"http://images.clipartof.com/Royalty-Free-RF-Clipart-Illustration-Of-A-Colorful-Your-Word-102478518.jpg",
			"http://www.designindaba.com/sites/default/files/news/street_0.jpg",
			"http://2.bp.blogspot.com/-OoeEOfRF2JM/US4xjtzXA0I/AAAAAAAAAtY/F3E2_MykbO8/s1600/YOUR+7.JPG",
			"http://www.designindaba.com/sites/default/files/u940/Your%20streen%20Kzn_0.jpg"}
	case "his":
		return []string{"http://thoughtsonleadership.biz/wp-content/uploads/2009/11/iStock_000007836244Medium-300x263.jpg",
			"http://blogs.hawaiibusiness.com/wp-content/uploads/2013/05/H.I.S.-Logo.jpeg",
			"http://i1.cpcache.com/product/1031507654/gay_pride_his_his_mugs_pillow_case.jpg?color=White&height=460&width=460&qv=90",
			"http://www.legitreviews.com/images/reviews/1559/logo_HIS.jpg"}
	case "her":
		return []string{"http://www.impawards.com/2013/posters/her.jpg",
			"http://dormtainment.com/wp-content/uploads/2014/01/HER-Cover.jpg",
			"http://upload.wikimedia.org/wikipedia/commons/8/87/Spike_Jonze_-_Her.jpg",
			"http://turntherightcorner.files.wordpress.com/2013/12/her-2013-title-movie-logo.jpg"}
	case "their":
		return []string{"http://mcdn1.teacherspayteachers.com/thumbitem/FREE-Their-There-and-Theyre-Posters/original-674070-1.jpg"}
	case "hers":
		return []string{"http://www.hardbodynews.com/wp-content/uploads/2008/12/jessmuscleandfitness.jpg",
			"http://blog-imgs-29-origin.fc2.com/a/i/b/aibo0714/hers_20080912.jpg",
			"https://www.google.com/search?q=my&rlz=1C5CHFA_enIL532IL532&es_sm=91&biw=1152&bih=558&source=lnms&tbm=isch&sa=X&ei=ZEZWVN6cCcztaJfKgbAI&ved=0CAYQ_AUoAQ&dpr=1.25#tbm=isch&q=hers&facrc=_&imgdii=_&imgrc=BJJ2vu3Bc2nthM%253A%3BW6mxRFcCrA8Q4M%3Bhttp%253A%252F%252Fwww.crystaljin.com%252Fwp-content%252Fuploads%252F2011%252F03%252FHersMuscleandFitness-MarchApril2011-cover.jpg%3Bhttp%253A%252F%252Fwww.pushgreen.net%252Ftag%252Fmuscle-and-fitness-hers-magazine-founder-eric%3B300%3B375"}
	case "this":
		return []string{"http://images.clipartof.com/small/210437-Royalty-Free-RF-Clipart-Illustration-Of-A-Retro-Black-And-White-Man-Carrying-A-Try-This-Sign.jpg"}
	case "me":
		return []string{"http://tigerwires.com/cms/lib3/PA01000001/Centricity/Domain/699/all-about-me-boy.jpg",
			"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTbX5aHsTJKydf1CvFwQRpywNBHul2Zj00W5ntNO-0Wac5kF20oag",
			"https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcTeVi392rLmG4arCzK-89-_ldS2dwXxFGmOCePvneoFx3ehujMl"}
	case "you":
		return []string{"http://cruciblecacophony.wikispaces.com/file/view/446px-Uncle_Sam_(pointing_finger).jpg/101768939/291x372/446px-Uncle_Sam_(pointing_finger).jpg",
			"http://images.sodahead.com/polls/000865673/polls_You_Logo_Black_xlarge.jpeg",
			"http://biblestudyoutlines.org/wp-content/uploads/2012/08/argument-clipart-1024x896.jpg"}
	case "he":
		return []string{"http://thoughtsonleadership.biz/wp-content/uploads/2009/11/iStock_000007836244Medium-300x263.jpg",
			"https://img0.etsystatic.com/000/0/5297554/il_fullxfull.251484430.jpg",
			"http://geezees.com/newsite/wp-content/uploads/NEW-he-she-vow-art-41.jpg"}
	case "she":
		return []string{"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTex4uJ5de2Ll3PZYB3qMPJHPh2Ww5udYTAsifEOKcfNJ9ZbrFV",
			"http://indulgy.ccio.co/5F/a1/p6/80783387035465957tahMtREic.jpg",
			"https://img0.etsystatic.com/000/0/5297554/il_fullxfull.251484430.jpg",
			"http://geezees.com/newsite/wp-content/uploads/NEW-he-she-vow-art-41.jpg"}
	case "they":
		return []string{"http://workwiththey.com/themes/site_themes/workwiththeydev/img/they-og.jpg",
			"http://static.giantbomb.com/uploads/original/2/21249/899580-they_logo.jpg",
			"http://us.cdn3.123rf.com/168nwm/qingwa/qingwa1110/qingwa111000237/11027373-the-word-they-carved-in-stone.jpg",
			"http://www.theyrep.com/wp-content/uploads/2012/04/THEY-website-about.jpg"}
	case "us":
		return []string{"http://www.walsingham.com/sites/default/files/photosymbols/individuals.jpg",
			"http://www.a2zscraplets.com.au/images/super/Us%20SCript%20Word.jpg"}
	case "in":
		return []string{"https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcRzky1Q6VHt4GMgEcD5iLRc3dIcv_Lptye5bNvcuLEekZvjFWfG",
			"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRUE1HfQa183x_RNqoWOGZ94HSq3sK-8NX5OlsNDFEI6oAeFzQU2Q",
			"http://media-cdn.tripadvisor.com/media/photo-s/01/ca/34/9e/hollywood-inn-express.jpg"}
	case "at":
		return []string{"https://encrypted-tbn2.gstatic.com/images?q=tbn:ANd9GcQsjCD6VdYRqZVeWUeMbUa0ciYZe6ASeatFSFSp7D9o-651vhsA",
			"https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcREAXvJAAH4JjOIH5RyXTR-kPWLNgGAIn8yj7ei6njYAcf45qK9CA"}
	case "that":
		return []string{"http://homepage2.nifty.com/that/pict/that-uzu.jpg",
			"http://1.bp.blogspot.com/-EVdkXyu7-uM/TbHDOJ2VcyI/AAAAAAAAOLE/_y5-rz-Nglk/s1600/t7.jpg"}
	case "if":
		return []string{"http://salmanshaheen.com/wp-content/uploads/2013/01/MDG-IF-campaign-against-009.jpg",
			"http://www.exceltraining.org.uk/excel-training-images/excel-if-statement.jpg",
			"http://upload.wikimedia.org/wikipedia/en/1/15/Msi_if_cover.jpg",
			"http://www.meg.ie/wp-content/uploads/2011/11/billryderjones_if_big.jpg"}
	case "else":
		return []string{"http://www.pretamarque.com/wp-content/uploads/2012/04/Joes-Jeans-Else.jpeg"}
	case "where":
		return []string{"http://m5.paperblog.com/i/6/61771/how-to-prepare-for-a-trip-L-HTeG64.jpeg",
			"http://4.bp.blogspot.com/_S3TFiuLoYtg/S_PdfSptbZI/AAAAAAAAAMk/nRrI4LB-3dc/s1600/Where+to+start.jpg",
			"http://sr.photos3.fotosearch.com/bthumb/CSP/CSP571/k5718115.jpg"}
	case "how":
		return []string{"https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcQojQilbSXMxich1CQTdgo3_v6HYYpnjlqzEeBVTxmS4aJJPefu1g",
			"http://ec.l.thumbs.canstockphoto.com/canstock7517808.jpg"}
	case "do":
		return []string{"http://venturebeat.files.wordpress.com/2011/11/do-com-logo.jpg",
			"http://www-cdf.fnal.gov/~canelli/gif_jpg/d0_ii.jpg",
			"http://payload.cargocollective.com/1/2/71370/857585/Do_Logo_800.jpg"}
	case "what":
		return []string{"https://patrecia1939.files.wordpress.com/2013/04/royalty-free-confused-clipart-illustration-215206.jpg",
			"https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcRb44ObPb1-o-bI4QunDTFK33bgU-6Hau62dlAOGDjxw_m0E0ms",
			"http://images.clipartpanda.com/confusion-clipart-royalty-free-confused-clipart-illustration-215196.jpg"}
	case "who":
		return []string{"http://www.lsac-lb.org/upload/homepage/profile-photo.jpg"}
	case "why":
		return []string{"http://generic.pixmac.com/4/cartoon-gorilla-scratching-his-head-cute-clipart-87469570.jpg",
			"http://www.clipartguide.com/_named_clipart_images/0511-1111-0316-0745_Funny_Looking_Cartoon_Man_with_the_Question_Scratching_His_Head_clipart_image.jpg",
			"http://www.clipartbest.com/cliparts/RTA/kb8/RTAkb8kTL.jpeg"}
	case "we":
		return []string{"http://www.fitnessblender.com/media/uploads/article-images/m_wii_console.jpg",
			"http://fc07.deviantart.net/fs44/i/2009/096/d/1/Just_ME_and_YOU_by_Add1ct3d.jpg",
			"http://www.wetheworld.org/img/we.jpg",
			"http://i301.photobucket.com/albums/nn44/pdubbleyah/weeee.jpg"}
	case "when":
		return []string{"http://i.kinja-img.com/gawker-media/image/upload/s--tx7OMnHv--/18fbt8v3pearxjpg.jpg",
			"http://www.suitqaisdiaries.com/wp-content/uploads/2013/12/Carpe-Diem-When.jpg"}
	case "as":
		return []string{"http://integrationacres.com/products/images/goat_1.jpg",
			"https://encrypted-tbn3.gstatic.com/images?q=tbn:ANd9GcQVNSsuOvLaSE48k-WMq7qKubBgCwxQHbsjkf_0uQtCI2Q05hy7"}
	case "let":
		return []string{
			"http://www.idiomsbykids.com/taylor/mrtaylor/class20022003/idioms/idioms2003/idioms5/let%20up..jpg"}
	case "for":
		return []string{
			"http://www.eceeldek.com/logo_kurumsal/for_logo.jpg",
			"http://itssoverycheri.com/wp-content/uploads/2014/08/4.jpeg",
			"http://cdn.mundodastribos.com/673943-Curiosidades-sobre-o-ser-humano-que-voc%C3%AA-n%C3%A3o-sabe5.jpg"}
	case "by":
		return []string{
			"http://rumol.org/wp-content/uploads/domen_by.jpg",
			"http://idahobusinessreview.com/files/2014/02/20-something-waving-goodbye.jpg"}
	}
}

func getImageUrls(word string, consumer *oauth.Consumer, accessToken *oauth.AccessToken) []string {
	response, err := consumer.Get(
		"https://yboss.yahooapis.com/ysearch/images",
		map[string]string{"q": word, "sites": "", "format": "json", "dimensions": "medium", "count": "10"},
		accessToken)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	unparsedResponse, err := ioutil.ReadAll(response.Body)

	parsedResponse := yahooResponse{}
	json.NewDecoder(strings.NewReader(string(unparsedResponse))).Decode(&parsedResponse)

	images := make([]string, 0, 35)
	for i := range parsedResponse.BossResponse.Images.Results {
		imageResult := parsedResponse.BossResponse.Images.Results[i]
		if imageResult.Format == "jpeg" || imageResult.Format == "jpg" {
			images = append(images, imageResult.Url)
		}
	}

	return images
}

func main() {
	c := oauth.NewConsumer(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.login.yahoo.com/oauth/v2/get_request_token",
			AuthorizeTokenUrl: "https://api.login.yahoo.com/oauth/v2/request_auth",
			AccessTokenUrl:    "https://api.login.yahoo.com/oauth/v2/get_token",
		})

	accessToken := &oauth.AccessToken{}
	accessToken.Token = os.Getenv("ACCESS_TOKEN")
	accessToken.Secret = os.Getenv("TOKEN_SECRET")
	accessToken.AdditionalData = map[string]string{"oauth_expires_in": "3600", "oauth_session_handle": os.Getenv("SESSION_HANDLE"), "oauth_authorization_expires_in": "732555938", "xoauth_yahoo_guid": os.Getenv("YAHOO_GUID")}

	m := martini.Classic()

	//curl -X POST  -k -u token: "http://localhost:3000/get_image.json" -d '{"urls": ["http://talks.golang.org/2013/advconc/gopherswim.jpg", "http://www.unixstickers.com/image/cache/data/stickers/golang/golang.sh-600x600.png"], "height_px": 400}' > test.png
	m.Get("/v1/imageUrls/:word", func(params martini.Params) []byte {

		imageUrls := findUrlsForWord(params["word"], c, accessToken)

		jsonString, _ := json.Marshal(imageUrls)

		return jsonString
	})

	m.Run()
}
