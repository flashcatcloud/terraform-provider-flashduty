# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default:

- **provider/provider.tf** - Example file for the provider index page
- **resources/`full resource name`/resource.tf** - Example file for the named resource page
- **data-sources/`full data source name`/data-source.tf** - Example file for the named data source page

## Running Examples

To run any example:

1. Set your Flashduty APP key:
   ```bash
   export FLASHDUTY_APP_KEY="your-app-key"
   ```

2. Navigate to the example directory:
   ```bash
   cd examples/resources/flashduty_team
   ```

3. Initialize and apply:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```
