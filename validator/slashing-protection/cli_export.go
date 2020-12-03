package slashingprotection

import (
	"encoding/json"
	"path/filepath"

	"github.com/prysmaticlabs/prysm/shared/cmd"
	"github.com/prysmaticlabs/prysm/shared/fileutil"
	"github.com/prysmaticlabs/prysm/validator/db/kv"
	"github.com/prysmaticlabs/prysm/validator/flags"
	export "github.com/prysmaticlabs/prysm/validator/slashing-protection/local/standard-protection-format"
	"github.com/urfave/cli/v2"
)

const (
	jsonExportFileName = "slashing_protection.json"
)

// ExportSlashingProtectionJSONCli extracts a validator's slashing protection
// history from their database and formats it into an EIP-3076 standard JSON
// file via a CLI entrypoint to make it easy to migrate machines or eth2 clients.
//
// Steps:
// 1. Parse a path to the validator's datadir from the CLI context.
// 2. Open the validator database file.
// 3. Call the function which actually exports the data from
// from the validator's db into an EIP standard slashing protection format
// 4. Format and save the JSON file to a user's specified output directory.
func ExportSlashingProtectionJSONCli(cliCtx *cli.Context) error {
	datadir := cliCtx.String(cmd.DataDirFlag.Name)
	validatorDB, err := kv.NewKVStore(datadir, nil)
	if err != nil {
		return err
	}
	eipJSON, err := export.ExportStandardProtectionJSON(cliCtx.Context, validatorDB)
	if err != nil {
		return err
	}
	outputDir := cliCtx.String(flags.SlashingProtectionExportDirFlag.Name)
	exists, err := fileutil.HasDir(outputDir)
	if err != nil {
		return err
	}
	if !exists {
		if err := fileutil.MkdirAll(outputDir); err != nil {
			return err
		}
	}
	outputFilePath := filepath.Join(outputDir, jsonExportFileName)
	encoded, err := json.Marshal(eipJSON)
	if err != nil {
		return err
	}
	return fileutil.WriteFile(outputFilePath, encoded)
}