package subcmd

import (
	"io/fs"
	"net/http"
	"time"

	gicfs "github.com/cyverse/go-irodsclient/fs"
	"github.com/cyverse/gocommands/cmd/flag"
	"github.com/cyverse/gocommands/commons"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var webServerCmd = &cobra.Command{
	Use:     "webserver [port]",
	Aliases: []string{"expose"},
	Short:   "Web serve iRODS",
	Long:    `This exposes iRODS through a web server.`,
	RunE:    processWebServerCommand,
	Args:    cobra.MinimumNArgs(1),
}

func AddWebServerCommand(rootCmd *cobra.Command) {
	// attach common flags
	flag.SetCommonFlags(webServerCmd)

	flag.SetTicketAccessFlags(webServerCmd)

	rootCmd.AddCommand(webServerCmd)
}

func processWebServerCommand(command *cobra.Command, args []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processWebServerCommand",
	})

	cont, err := flag.ProcessCommonFlags(command)
	if err != nil {
		return xerrors.Errorf("failed to process common flags: %w", err)
	}

	if !cont {
		return nil
	}

	// handle local flags
	_, err = commons.InputMissingFields()
	if err != nil {
		return xerrors.Errorf("failed to input missing fields: %w", err)
	}

	ticketAccessFlagValues := flag.GetTicketAccessFlagValues()

	appConfig := commons.GetConfig()
	syncAccount := false
	if len(ticketAccessFlagValues.Name) > 0 {
		logger.Debugf("use ticket: %s", ticketAccessFlagValues.Name)
		appConfig.Ticket = ticketAccessFlagValues.Name
		syncAccount = true
	}

	if syncAccount {
		err := commons.SyncAccount()
		if err != nil {
			return err
		}
	}

	// Create a file system
	account := commons.GetAccount()
	filesystem, err := gicfs.NewFileSystemWithDefault(account, "webserver")
	if err != nil {
		return xerrors.Errorf("failed to get iRODS FS Client: %w", err)
	}

	defer filesystem.Release()

	http.Handle("/", http.FileServerFS(&irodsfsasgofs{filesystem}))
	// http.Handle("/", http.FileServerFS(filesystem))

	logger.Fatal(http.ListenAndServe(args[0], nil))


	return nil
}


type irodsentryasgofi struct { // convert irods entry to go fs fileinfo
	entry *gicfs.Entry
}
func (gen irodsentryasgofi) Name() string { return gen.entry.Name }
func (gen irodsentryasgofi) Size() int64  { return gen.entry.Size }
func (gen irodsentryasgofi) Mode() fs.FileMode  { if(gen.entry.IsDir()) {return fs.ModeDir|0777}; return 0777 }
func (gen irodsentryasgofi) ModTime() time.Time  { return gen.entry.ModifyTime }
func (gen irodsentryasgofi) IsDir() bool  { return gen.entry.IsDir() }
func (gen irodsentryasgofi) Sys() any  { return nil }

type irodsfhasgofh struct {
	*gicfs.FileHandle
	// ifs *gicfs.FileSystem
	ifs *irodsfsasgofs
	name string
}
func (gfh irodsfhasgofh) Stat() (fs.FileInfo, error) {
	entry,err := gfh.ifs.Stat(gfh.name)
	return irodsentryasgofi{entry}, err
}

type irodsfsasgofs struct {
	*gicfs.FileSystem
	// gicfs.FileSystem
}
func (gfs *irodsfsasgofs) Open(name string) (fs.File, error) {
	ifh, err := gfs.OpenFile(name, "", "r")
	return irodsfhasgofh{ifh,gfs,name}, err
}