package archiver

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"
	"log"
	"bufio"
	"backups/commands"
	"backups/tarutils"
)

const ArchivePath = "/tmp/archiver"

type Archiver struct {
	md5s map[string]string
}

func NewArchiver() Archiver {
	return Archiver{md5s: make(map[string]string)}
}

func (arch *Archiver) Transfer(archReq commands.Archive, writer *bufio.Writer) {
	sourcePath := archReq.Path
	tarPath := targetPath(sourcePath)
	currentMD5, err := tar(sourcePath, tarPath)
	log.Println("Archiver: archived path ",sourcePath, " md5 ", currentMD5, "err ", err)
	if arch.md5s[sourcePath] == currentMD5 {
		commands.WriteCommand(writer, commands.Acknowledge{})
	} else {
		sendFile(writer, sourcePath)
		arch.md5s[sourcePath] = currentMD5
	}
}

func targetPath(sourcePath string) string {
	return  ArchivePath + "/" + strings.ReplaceAll(sourcePath, "/", "_")
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

func sendFile(writer *bufio.Writer, tarPath string)  {
	file, err:= os.Open(tarPath)
	if err != nil {
		commands.WriteCommand(writer, commands.Acknowledge{Err: err})
		return
	}
	fi, err := file.Stat()
	if err != nil {
		commands.WriteCommand(writer, commands.Acknowledge{Err: err})
		return
	}
	defer file.Close()
	commands.WriteCommand(writer, commands.Transfer{Size: int(fi.Size()), Reader:file})
	// TODO maybe delete file after transfer
}
