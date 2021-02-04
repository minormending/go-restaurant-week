package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/minormending/go-restaurant-week/models"
)

// RestaurantWeekResponse is the response from the API
type RestaurantWeekResponse struct {
	Data []struct {
		BlockID     string      `json:"blockId"`
		BlockOrder  int         `json:"blockOrder"`
		Passthrough interface{} `json:"passthrough"`
		GridItems   []struct {
			ChannelName      string `json:"channelName"`
			ChannelShortName string `json:"channelShortName"`
			DisplayTitle     string `json:"displayTitle"`
			Ecommerce        []struct {
				ButtonText    interface{} `json:"buttonText,omitempty"`
				MaxPrice      interface{} `json:"maxPrice,omitempty"`
				MinPrice      interface{} `json:"minPrice,omitempty"`
				PartnerID     interface{} `json:"partnerId"`
				PartnerName   string      `json:"partnerName"`
				PreferredLink string      `json:"preferredLink,omitempty"`
				URL           string      `json:"url,omitempty"`
			} `json:"ecommerce"`
			EndDate           interface{} `json:"endDate"`
			ID                string      `json:"id"`
			HasOfferAvailable bool        `json:"hasOfferAvailable"`
			IsAllDay          bool        `json:"isAllDay"`
			IsFeatured        bool        `json:"isFeatured"`
			IsOngoing         bool        `json:"isOngoing"`
			IsSponsored       bool        `json:"isSponsored"`
			ImagePath         string      `json:"imagePath"`
			Latitude          float64     `json:"latitude"`
			LogoPath          interface{} `json:"logoPath"`
			Longitude         float64     `json:"longitude"`
			LookupInfo        []struct {
				Ids        string `json:"ids"`
				LookupName string `json:"lookupName"`
			} `json:"lookupInfo"`
			NextOccurrence      interface{} `json:"nextOccurrence"`
			PrimaryCategoryID   string      `json:"primaryCategoryId"`
			PrimaryCategoryName string      `json:"primaryCategoryName"`
			PrimaryLocationName string      `json:"primaryLocationName"`
			ReviewScore         interface{} `json:"reviewScore"`
			Phone               string      `json:"phone"`
			ShortTitle          interface{} `json:"shortTitle"`
			SortTitle           string      `json:"sortTitle"`
			StartDate           interface{} `json:"startDate"`
			Summary             string      `json:"summary"`
			URL                 string      `json:"url"`
			Website             string      `json:"website"`
		} `json:"gridItems"`
	} `json:"data"`
	Lookup struct {
		Borough []struct {
			ID         string `json:"id"`
			IsSelected bool   `json:"isSelected"`
			Name       string `json:"name"`
			URLTitle   string `json:"urlTitle"`
		} `json:"borough"`
		Neighborhood []struct {
			BoroughID     string `json:"boroughId"`
			GooglePlaceID string `json:"googlePlaceId"`
			ID            string `json:"id"`
			IsSelected    bool   `json:"isSelected"`
			Name          string `json:"name"`
			URLTitle      string `json:"urlTitle"`
		} `json:"neighborhood"`
		Cuisine []struct {
			ID         string `json:"id"`
			IsSelected bool   `json:"isSelected"`
			Name       string `json:"name"`
			URLTitle   string `json:"urlTitle"`
		} `json:"cuisine"`
		RestaurantWeekAmenities []struct {
			ID         string `json:"id"`
			IsSelected bool   `json:"isSelected"`
			Name       string `json:"name"`
			URLTitle   string `json:"urlTitle"`
		} `json:"restaurantWeekAmenities"`
		RestaurantWeekMeals []struct {
			ID         string `json:"id"`
			IsSelected bool   `json:"isSelected"`
			Name       string `json:"name"`
			URLTitle   string `json:"urlTitle"`
		} `json:"restaurantWeekMeals"`
		RestaurantWeekDeliveryPlatforms []struct {
			ID         string `json:"id"`
			IsSelected bool   `json:"isSelected"`
			Name       string `json:"name"`
			URLTitle   string `json:"urlTitle"`
		} `json:"restaurantWeekDeliveryPlatforms"`
		Collections []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Image       string `json:"image"`
			IsSelected  bool   `json:"isSelected"`
			ShortTitle  string `json:"shortTitle"`
			SortOrder   int    `json:"sortOrder"`
			URLTitle    string `json:"urlTitle"`
			PartnerText string `json:"partnerText"`
		} `json:"collections"`
	} `json:"lookup"`
}

// GetRestaurantInfo returns restaurant week info
func GetRestaurantInfo() ([]*models.Restaurant, error) {
	res, err := http.Get("https://service.nycgo.com/nycgo/v2/body-grid-blocks?entryId=411&gridId=restaurant-week")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	payload := RestaurantWeekResponse{}
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&payload); err != nil {
		return nil, err
	}

	restaurants := []*models.Restaurant{}
	for _, r := range payload.Data[0].GridItems {
		restaurant := &models.Restaurant{
			ID:          r.ID,
			Name:        r.DisplayTitle,
			Latitude:    r.Latitude,
			Longitude:   r.Longitude,
			Description: r.Summary,
			Website:     r.Website,
			NYCLink:     fmt.Sprintf("https://www.nycgo.com/restaurant-week/browse/%s", r.ID),
		}
		cuisineIDs := []string{}
		for _, lookup := range r.LookupInfo {
			if lookup.LookupName == "cuisine" {
				cuisineIDs = append(cuisineIDs, strings.Split(lookup.Ids, ",")...)
			}
		}
		cuisines := []string{}
		for _, cuisine := range payload.Lookup.Cuisine {
			for _, cuisineID := range cuisineIDs {
				if cuisineID == cuisine.ID {
					cuisines = append(cuisines, cuisine.Name)
				}
			}
		}
		restaurant.Cuisine = strings.Join(cuisines, ", ")
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}
