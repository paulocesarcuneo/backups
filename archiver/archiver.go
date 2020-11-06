package archiver

import (
	"backups/commands"
	"backups/tarutils"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"strings"
)

const ArchivePath = "/tmp/archiver"

type Archiver struct {
	md5s map[string]string
}

func NewArchiver() Archiver {
	return Archiver{md5s: make(map[string]string)}
}

func (archiver *Archiver) Handle(request commands.Command) commands.Command {
	switch archive := request.(type) {
	case commands.Archive:
		return archiver.transfer(archive)
	default:
		return commands.Error{Error: "Unhable to handle command"}
	}
}

func (arch *Archiver) transfer(archReq commands.Archive) commands.Command {
	sourcePath := archReq.Path
	tarPath := targetPath(sourcePath)
	currentMD5, err := tar(sourcePath, tarPath)
	log.Println("Archiver: archived path ", sourcePath, " md5 ", currentMD5, "err ", err)
	if arch.md5s[sourcePath] == currentMD5 {
		return commands.Acknowledge{}
	} else {
		result, err := openFileAsTransfer(tarPath)
		if err != nil {
			log.Println("Archiver: Transfer ", err)
			return commands.Error{Error: err.Error()}
		}
		arch.md5s[sourcePath] = currentMD5
		return result
	}
}

func targetPath(sourcePath string) string {
	return ArchivePath + "/" + strings.ReplaceAll(sourcePath, "/", "_")
}

func tar(path string, dst string) (string, error) {
	file, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer file.Close()
	log.Println("making tar gz for", path)
	md5writer := md5.New()
	err = tarutils.TarFolder(path, file, md5writer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(md5writer.Sum(nil)), nil
}

func openFileAsTransfer(tarPath string) (commands.Command, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return commands.Transfer{
		Size:   fi.Size(),
		Reader: file,
	}, nil
}

func init() {
	path := ArchivePath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
