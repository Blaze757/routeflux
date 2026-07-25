package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newRoutingCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routing",
		Short: "Manage geosite/geoip routing rules and geo data files",
	}

	cmd.AddCommand(
		newRoutingUpdateGeoCmd(opts),
		newRoutingGetCmd(opts),
		newRoutingSetCmd(opts),
		newRoutingUpdateGeoWithCheckCmd(opts),
	)

	return cmd
}

func newRoutingUpdateGeoCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "update-geo",
		Short: "Download geosite.dat and geoip.dat from configured URLs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.service.UpdateGeoData(cmd.Context()); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Geo data updated successfully")
			return nil
		},
	}
}

func newRoutingGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show current routing settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			settings, err := opts.service.GetSettings()
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(),
				"direct-geosite=%s\ndirect-geoip=%s\ngeo-update-enabled=%t\ngeo-update-interval=%s\ngeo-update-urls=%s\n",
				strings.Join(settings.Routing.DirectGeosite, ", "),
				strings.Join(settings.Routing.DirectGeoIP, ", "),
				settings.GeoUpdate.Enabled,
				settings.GeoUpdate.Interval,
				formatGeoURLs(settings.GeoUpdate.URLs),
			)
			return nil
		},
	}
}

func formatGeoURLs(urls map[string]string) string {
	parts := make([]string, 0, len(urls))
	for k, v := range urls {
		parts = append(parts, k+"="+v)
	}
	return strings.Join(parts, ", ")
}

func newRoutingSetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Update a routing setting",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			var settingsKey string
			switch key {
			case "direct-geosite":
				settingsKey = "routing.direct-geosite"
			case "direct-geoip":
				settingsKey = "routing.direct-geoip"
			case "geo-update-enabled":
				settingsKey = "geo-update.enabled"
			case "geo-update-interval":
				settingsKey = "geo-update.interval"
			case "geo-url-geoip":
				settingsKey = "geo-update.url-geoip"
			case "geo-url-geosite":
				settingsKey = "geo-update.url-geosite"
			default:
				return fmt.Errorf("unknown routing setting %q: supported keys are direct-geosite, direct-geoip, geo-update-enabled, geo-update-interval, geo-url-geoip, geo-url-geosite", key)
			}

			settings, err := opts.service.SetSetting(settingsKey, value)
			if err != nil {
				return err
			}
			return printOutput(cmd, opts.jsonOutput, settings, fmt.Sprintf("Updated %s=%s", key, value))
		},
	}
}

func newRoutingUpdateGeoWithCheckCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:    "check-update",
		Short:  "Check and update geo data if stale",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.service.CheckGeoUpdate(cmd.Context()); err != nil {
				return err
			}
			return nil
		},
	}
}
