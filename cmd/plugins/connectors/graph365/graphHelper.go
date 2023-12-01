package main

import (
	"time"
)

/*func (g *Graph365) InitializeGraphForUserAuth() error {
	clientId := g.ClientID
	tenantId := g.TenantID
	scopes := []string{"ExternalItem.Read.All", "Files.Read.All", "Sites.Read.All", "User.Read.All", "User.ReadWrite.All"}
	g.Scopes = scopes

	// Create the device code credential
	credential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientId,
		TenantID: tenantId,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			fmt.Println(message.Message)
			return nil
		},
	})
	if err != nil {
		return err
	}

	g.deviceCodeCredential = credential

	// Create an auth provider using the credential
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, g.graphUserScopes)
	if err != nil {
		return err
	}

	// Create a request adapter using the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return err
	}

	// Create a Graph client using request adapter
	_ = msgraphsdk.NewGraphServiceClient(adapter)

	return nil
}*/

type DriveSearchResponse struct {
	Value []struct {
		SearchTerms    []interface{} `json:"searchTerms"`
		HitsContainers []struct {
			Hits []struct {
				HitID    string `json:"hitId"`
				Rank     int    `json:"rank"`
				Summary  string `json:"summary"`
				Resource struct {
					OdataType      string `json:"@odata.type"`
					Size           int    `json:"size"`
					FileSystemInfo struct {
						CreatedDateTime      time.Time `json:"createdDateTime"`
						LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
					} `json:"fileSystemInfo"`
					ListItem struct {
						OdataType string `json:"@odata.type"`
						Fields    struct {
						} `json:"fields"`
						ID string `json:"id"`
					} `json:"listItem"`
					ID        string `json:"id"`
					CreatedBy struct {
						User struct {
							DisplayName string `json:"displayName"`
							Email       string `json:"email"`
						} `json:"user"`
					} `json:"createdBy"`
					CreatedDateTime time.Time `json:"createdDateTime"`
					LastModifiedBy  struct {
						User struct {
							DisplayName string `json:"displayName"`
							Email       string `json:"email"`
						} `json:"user"`
					} `json:"lastModifiedBy"`
					LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
					Name                 string    `json:"name"`
					ParentReference      struct {
						DriveID       string `json:"driveId"`
						ID            string `json:"id"`
						SharepointIds struct {
							ListID           string `json:"listId"`
							ListItemID       string `json:"listItemId"`
							ListItemUniqueID string `json:"listItemUniqueId"`
						} `json:"sharepointIds"`
						SiteID string `json:"siteId"`
					} `json:"parentReference"`
					WebURL string `json:"webUrl"`
				} `json:"resource"`
			} `json:"hits"`
			Total                int  `json:"total"`
			MoreResultsAvailable bool `json:"moreResultsAvailable"`
		} `json:"hitsContainers"`
	} `json:"value"`
	OdataContext string `json:"@odata.context"`
}

type DriveSearchRequestItem struct {
	EntityTypes []string `json:"entityTypes"`
	Query       struct {
		QueryString string `json:"queryString"`
	} `json:"query"`
	From int `json:"from"`
	Size int `json:"size"`
}

type DriveSearchRequest struct {
	Requests []DriveSearchRequestItem `json:"requests"`
}
