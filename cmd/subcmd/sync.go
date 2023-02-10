package subcmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cyverse/gocommands/commons"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync i:[collection] [local dir] or sync [local dir] i:[collection]",
	Short: "Sync local directory with iRODS collection",
	Long:  `This synchronizes a local directory with the given iRODS collection.`,
	RunE:  processSyncCommand,
}

func AddSyncCommand(rootCmd *cobra.Command) {
	// attach common flags
	commons.SetCommonFlags(syncCmd)

	syncCmd.Flags().Bool("progress", false, "Display progress bar")
	syncCmd.Flags().Bool("no_hash", false, "Compare files without using md5 hash")

	rootCmd.AddCommand(syncCmd)
}

func processSyncCommand(command *cobra.Command, args []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processSyncCommand",
	})

	cont, err := commons.ProcessCommonFlags(command)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	if !cont {
		return nil
	}

	// handle local flags
	_, err = commons.InputMissingFields()
	if err != nil {
		logger.Error(err)
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	progress := false
	progressFlag := command.Flags().Lookup("progress")
	if progressFlag != nil {
		progress, err = strconv.ParseBool(progressFlag.Value.String())
		if err != nil {
			progress = false
		}
	}

	noHash := false
	noHashFlag := command.Flags().Lookup("no_hash")
	if noHashFlag != nil {
		noHash, err = strconv.ParseBool(noHashFlag.Value.String())
		if err != nil {
			noHash = false
		}
	}

	config := commons.GetConfig()
	replicate := !config.NoReplication

	// Create a file system
	account := commons.GetAccount()
	filesystem, err := commons.GetIRODSFSClient(account)
	if err != nil {
		logger.Error(err)
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	defer filesystem.Release()

	if len(args) < 2 {
		err := fmt.Errorf("not enough input arguments")
		logger.Error(err)
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	targetPath := "i:./"
	sourcePaths := args[:]

	if len(args) >= 2 {
		targetPath = args[len(args)-1]
		sourcePaths = args[:len(args)-1]
	}

	parallelJobManager := commons.NewParallelJobManager(filesystem, commons.MaxThreadNumDefault, progress)
	parallelJobManager.Start()

	for _, sourcePath := range sourcePaths {
		if strings.HasPrefix(sourcePath, "i:") {
			if strings.HasPrefix(targetPath, "i:") {
				// copy
				err = copyOne(parallelJobManager, sourcePath[2:], targetPath[2:], true, false, true, noHash)
				if err != nil {
					logger.Error(err)
					fmt.Fprintln(os.Stderr, err.Error())
					return nil
				}
			} else {
				// get
				err = getOne(parallelJobManager, sourcePath[2:], targetPath, false, true, noHash)
				if err != nil {
					logger.Error(err)
					fmt.Fprintln(os.Stderr, err.Error())
					return nil
				}
			}
		} else {
			if strings.HasPrefix(targetPath, "i:") {
				// put
				err = putOne(parallelJobManager, sourcePath, targetPath[2:], false, replicate, true, noHash)
				if err != nil {
					logger.Error(err)
					fmt.Fprintln(os.Stderr, err.Error())
					return nil
				}
			} else {
				// local to local
				err := fmt.Errorf("syncing between local files/directories is not supported")
				logger.Error(err)
				fmt.Fprintln(os.Stderr, err.Error())
				return nil
			}
		}
	}

	parallelJobManager.DoneScheduling()
	err = parallelJobManager.Wait()
	if err != nil {
		logger.Error(err)
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}

	return nil
}
