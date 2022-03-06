package sonarcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	pl "github.com/reinoudk/go-sonarcloud/sonarcloud/project_links"
)

type dataSourceProjectLinksType struct{}

func (d dataSourceProjectLinksType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This datasource retrieves the list of links for the given project.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"project_key": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The key of the project.",
			},
			"links": {
				Computed:    true,
				Description: "The links of this project.",
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:        types.StringType,
						Computed:    true,
						Description: "ID of the link.",
					},
					"name": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The name the link.",
					},
					"type": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The type of the link.",
					},
					"url": {
						Type:        types.StringType,
						Computed:    true,
						Description: "The url of the link.",
					},
				}),
			},
		},
	}, nil
}

func (d dataSourceProjectLinksType) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return dataSourceProjectLinks{
		p: *(p.(*provider)),
	}, nil
}

type dataSourceProjectLinks struct {
	p provider
}

func (d dataSourceProjectLinks) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var config DataProjectLinks
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := pl.SearchRequest{
		ProjectKey: config.ProjectKey.Value,
	}

	response, err := d.p.client.ProjectLinks.Search(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read the project's links",
			fmt.Sprintf("The Search request returned an error: %+v", err),
		)
		return
	}

	links := make([]DataProjectLink, len(response.Links))
	for i, link := range response.Links {
		links[i] = DataProjectLink{
			Id:   types.String{Value: link.Id},
			Name: types.String{Value: link.Name},
			Type: types.String{Value: link.Type},
			Url:  types.String{Value: link.Url},
		}
	}

	result := DataProjectLinks{
		ID:         types.String{Value: config.ProjectKey.Value},
		ProjectKey: config.ProjectKey,
		Links:      links,
	}

	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}
