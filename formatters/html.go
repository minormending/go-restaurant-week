package formatters

import (
	"html/template"
	"io"

	"github.com/minormending/go-restaurant-week/models"
)

var pageTemplate = `<!DOCTYPE html><html>
<head>
	<title>NYC Restaurant Week 2021</title>
	<script src="https://polyfill.io/v3/polyfill.min.js?features=default"></script>
	<script src="https://maps.googleapis.com/maps/api/js?key={{.APIKey}}&callback=initMap&libraries=&v=weekly" defer></script>
	<style type="text/css">
		#map { 
			height: 100%; 
		}
		html, 
		body {
			height: 100%;
			margin: 0;
			padding: 0;
		}
	</style>
	<script>
		function initMap() {
			const nyc = { lat: 40.7448362, lng: -73.9584712 };
			const map = new google.maps.Map(document.getElementById("map"), {
				zoom: 13,
				center: nyc,
			});
			const markerIcon = "iVBORw0KGgoAAAANSUhEUgAAABQAAAAgCAYAAAASYli2AAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyZpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuNi1jMDE0IDc5LjE1Njc5NywgMjAxNC8wOC8yMC0wOTo1MzowMiAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENDIDIwMTQgKFdpbmRvd3MpIiB4bXBNTTpJbnN0YW5jZUlEPSJ4bXAuaWlkOkQxNjg4MUQzREIyQzExRTVCN0E2RkU3MTY5RTFGOTMzIiB4bXBNTTpEb2N1bWVudElEPSJ4bXAuZGlkOkQxNjg4MUQ0REIyQzExRTVCN0E2RkU3MTY5RTFGOTMzIj4gPHhtcE1NOkRlcml2ZWRGcm9tIHN0UmVmOmluc3RhbmNlSUQ9InhtcC5paWQ6RDE2ODgxRDFEQjJDMTFFNUI3QTZGRTcxNjlFMUY5MzMiIHN0UmVmOmRvY3VtZW50SUQ9InhtcC5kaWQ6RDE2ODgxRDJEQjJDMTFFNUI3QTZGRTcxNjlFMUY5MzMiLz4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz6klcgWAAAC+UlEQVR42qyWz08TURDH5+2PtkpQ8aI3SYwgGjWCSlBRMJBwQA9qYg8m/ggYoxfjxX/ARE/oQSOIB2Piz2gaD0YPGEBpaymiIMgPLWAbUw+gxgot7e4+Z9YnEukPKvuSb/b1ze6nM/tmZh/jPtUFAKtRCVjYUFEBRdMhX5FhLedgLITGGEjI0pXgF/7epsI0WDDiCfQwbwkrycuFQiuA3yKwTJIYcLBoEEsyODCrgMSSuIVAYikZbuAtft7+zGsM0e+aMqmwaivbjTua0gklzY7Fzl7Wm68+NPrxp5kFDXeMN6cPSu8azsj1mBmOpO8xFfCp13iOsFc4Je/eCg3RGtlSbkwqg6uNt+FlANUnoENiPiBs2YXssEEIL59Q37E8zSpipQkK3RC27Dw8vlei+p78AzM36fd8UtiyA25aw+pfNCnL/12nNbJlDbTbIH/HRkbvyolaIeSkNbJl/Q7Nf5OgAC93KUyxlCNJ6ZNbmWcR5My7nsHigREkVEuB1GlBx+5vQYvA/FwqNbrrmqZ1e9SELpA68pn75Za2bsfYVMVEVdFoid2uKyAp/wXzvuPuLUe0SzKlXG//aHx0qiK8rzJULGuakqGrJYN5t9dp13EaImAcpfUPjEXHflSHazYEilUHejrPjOro4Z6d9doNipgki+/xFEF7+wLRj5HycG1JcLMq62omTxHmLj8xAxtGTSicc4Mx9lN0FnjgegmJeCm/fa7r1CIlvhgUm9m6/7ZxHfdONj0TsADqA8HIOTOxe7r9XJQXQQddT3zek9f2XIkl1CjSza/4TB1wAzzdCYI1z4KNz5w80EMCQihoOkhP5qLWo/YfdlZfjLVClLvxNh8G4wHecavIg7ajqHLUyjnlOxs4EhieDV1H0EMHKi/E2lmUtwP33tvmw7VjqF1JYclqudPXQeHTO6WuPHj/Uav//OPam13h0tdlzs5GEaa5AZHIpDbnjJPu/CO6zCpUgZiHBGw81WlNzpBm9FBMfEa/ooLpYDR+CTAAX0lK4qp7J98AAAAASUVORK5CYII=";
			const infowindow = new google.maps.InfoWindow();

			{{range $r := .Restaurants}}
			const marker{{$r.ID}} = new google.maps.Marker({
				title: "{{$r.Name}}",
				position: { lat: {{$r.Latitude}}, lng: {{$r.Longitude}} },
				map: map,
				icon: "data:image/png;base64," + markerIcon,
			});
			marker{{$r.ID}}.addListener("click", () => {
				const site = {{$r.Website}};
				const contentString =
					'<div id="content">' +
					'<div id="siteNotice">' +
					"</div>" +
					'<h3 id="firstHeading" class="firstHeading">{{$r.Name}} ({{$r.Cuisine}})</h3>' +
					'<div id="bodyContent">' +
					"<p>{{$r.Description}}</p>" +
					'<p>Website: <a href="'+site+'" target="_blank">'+site+'</a></p>' +
					'<p>Details: <a href="https://www.nycgo.com/restaurant-week/browse/{{$r.ID}}" target="_blank">https://www.nycgo.com/restaurant-week/browse/{{$r.ID}}</a></p>' +
					"</div>" +
					"</div>";
				infowindow.setContent(contentString);
				infowindow.open(map, marker{{$r.ID}});
			});
			{{end}}
		}
</script>
</head>
<body>
    <div id="map"></div>
</body>
</html>
`

type context struct {
	APIKey      string
	Restaurants []*restaurantContext
}

type restaurantContext struct {
	ID          template.JS
	Name        template.JS
	Latitude    float64
	Longitude   float64
	Cuisine     template.JS
	Description template.JS
	Website     template.URL
}

// ToHTML generates an single page map of the restaurants
func ToHTML(wr io.Writer, apiKey string, restaurants []*models.Restaurant) error {
	t, err := template.New("htmlPage").Parse(pageTemplate)
	if err != nil {
		return err
	}

	contextItems := []*restaurantContext{}
	for _, r := range restaurants {
		contextItems = append(contextItems, &restaurantContext{
			ID:          template.JS(r.ID),
			Name:        template.JS(r.Name),
			Latitude:    r.Latitude,
			Longitude:   r.Longitude,
			Cuisine:     template.JS(r.Cuisine),
			Description: template.JS(r.Description),
			Website:     template.URL(r.Website),
		})
	}

	return t.Execute(wr, context{
		APIKey:      apiKey,
		Restaurants: contextItems,
	})
}
