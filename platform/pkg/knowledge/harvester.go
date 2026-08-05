package knowledge

import (
	"context"
	"fmt"
	"time"

	"github.com/wgmasvix-hue/ChengetAi-learn-/platform/pkg/dspace"
	"go.uber.org/zap"
)

// Harvester handles knowledge harvesting from DSpace
type Harvester struct {
	dspaceClient *dspace.Client
	logger       *zap.SugaredLogger
}

// NewHarvester creates a new knowledge harvester
func NewHarvester(dspaceClient *dspace.Client, logger *zap.SugaredLogger) *Harvester {
	return &Harvester{
		dspaceClient: dspaceClient,
		logger:       logger,
	}
}

// HarvestedResource represents a harvested resource from DSpace
type HarvestedResource struct {
	DSpaceID    string
	Title       string
	Description string
	Creator     string
	ResourceType string
	Subject     string
	Language    string
	PublicationDate *time.Time
	FileURL     string
	MimeType    string
	FileSize    int64
	Metadata    map[string]interface{}
}

// HarvestCommunity harvests all resources from a community
func (h *Harvester) HarvestCommunity(ctx context.Context, communityID string) ([]HarvestedResource, error) {
	h.logger.Infow("Starting community harvest", "community_id", communityID)

	// Get collections in community
	collections, err := h.dspaceClient.GetCollectionsByID(ctx, communityID)
	if err != nil {
		return nil, fmt.Errorf("getting collections: %w", err)
	}

	var resources []HarvestedResource

	// Harvest each collection
	for _, collection := range collections {
		collResources, err := h.HarvestCollection(ctx, collection.ID)
		if err != nil {
			h.logger.Errorw("Error harvesting collection", "collection_id", collection.ID, "error", err)
			continue
		}
		resources = append(resources, collResources...)
	}

	h.logger.Infow("Community harvest complete", "community_id", communityID, "resources_count", len(resources))
	return resources, nil
}

// HarvestCollection harvests all resources from a collection
func (h *Harvester) HarvestCollection(ctx context.Context, collectionID string) ([]HarvestedResource, error) {
	h.logger.Infow("Starting collection harvest", "collection_id", collectionID)

	var allResources []HarvestedResource
	const limit = 50
	offset := 0

	for {
		// Get items from collection
		items, err := h.dspaceClient.GetItemsByCollection(ctx, collectionID, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("getting items: %w", err)
		}

		if len(items) == 0 {
			break
		}

		// Convert items to resources
		for _, item := range items {
			resource := h.convertItemToResource(item)
			allResources = append(allResources, resource)
		}

		offset += limit
	}

	h.logger.Infow("Collection harvest complete", "collection_id", collectionID, "resources_count", len(allResources))
	return allResources, nil
}

// convertItemToResource converts a DSpace item to a HarvestedResource
func (h *Harvester) convertItemToResource(item dspace.Item) HarvestedResource {
	resource := HarvestedResource{
		DSpaceID: item.ID,
		Title:    item.Name,
		Metadata: make(map[string]interface{}),
	}

	// Extract metadata
	for _, m := range item.Metadata {
		if len(m.Value) > 0 {
			switch m.Key {
			case "dc.description":
				resource.Description = m.Value[0]
			case "dc.creator":
				resource.Creator = m.Value[0]
			case "dc.type":
				resource.ResourceType = m.Value[0]
			case "dc.subject":
				resource.Subject = m.Value[0]
			case "dc.language.iso":
				resource.Language = m.Value[0]
			case "dc.issued":
				if t, err := time.Parse("2006-01-02", m.Value[0]); err == nil {
					resource.PublicationDate = &t
				}
			}
			// Store all metadata
			resource.Metadata[m.Key] = m.Value
		}
	}

	// Extract bitstream information if available
	if len(item.Bitstreams) > 0 {
		bitstream := item.Bitstreams[0]
		resource.FileURL = fmt.Sprintf("%s/bitstream/%s", item.Handle, bitstream.Name)
		resource.MimeType = bitstream.MimeType
		resource.FileSize = bitstream.Size
	}

	return resource
}
