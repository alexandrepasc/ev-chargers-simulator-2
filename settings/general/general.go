//nolint:unused,gofmt,revive,gocritic,deadcode,goimports,nolintlint,whitespace,staticcheck,wsl
package general

import (
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
)

/*
 */
func General(fp string, fl *flags.Flags) {}

/*
Checks if the general settings file exists in the defined path.

If exists returns true, if not returns false.

fp	-	Defined general file path (string)
*/
func isGeneral(fp string) {

	_, err := os.Stat(fp + "/" + generalFile)

	if err != nil {}
}
