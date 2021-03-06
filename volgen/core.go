//Core file for graph and volfile generation

package volgen

import (
	"fmt"
	"github.com/gluster/glusterd2/volume"
	"os"
	"strings"
)

const (
	workdir = "/var/lib/glusterd"
)

var (
	GenerateVolfileFunc = GenerateVolfile
)

// GenerateVolfile function will do all task from graph generation to volfile generation
func GenerateVolfile(vinfo *volume.Volinfo) error {
	var cpath string
	var err error
	var f *os.File

	graph := GenerateGraph(vinfo, "CLIENT")

	//Generate client volfile
	err = getClientFilePath(vinfo, &cpath)
	if err != nil {
		return err
	}

	f, err = os.Create(cpath)
	if err != nil {
		return err
	}
	graph.DumpGraph(f)
	f.Close()

	//Generate volfile for server
	for _, b := range vinfo.Bricks {
		sgraph := GenerateGraph(vinfo, "SERVER")

		//Generate volfile for server
		err = getServerFilePath(vinfo, &cpath, b)
		if err != nil {
			return err
		}

		f, err = os.Create(cpath)
		if err != nil {
			return err
		}

		sgraph.DumpGraph(f)
		f.Close()
	}
	return nil
}

func getVolumeDir(vinfo *volume.Volinfo, dir *string) {
	*dir = fmt.Sprintf("%s/vols/%s", workdir, vinfo.Name)
}

func getServerFilePath(vinfo *volume.Volinfo, path *string, brickinfo volume.Brickinfo) error {
	var vdir string
	getVolumeDir(vinfo, &vdir)

	//TODO : take the hostname from brickinfo
	hname := brickinfo.Hostname
	// Create volume directory (/var/lib/glusterd/vols/<VOLNAME>)
	err := os.MkdirAll(vdir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	*path = fmt.Sprintf("%s/%s.%s.%s.vol", vdir, vinfo.Name, hname, strings.Replace(brickinfo.Path, "/", "-", -1))
	return nil
}

func getClientFilePath(vinfo *volume.Volinfo, path *string) error {
	var vdir string
	getVolumeDir(vinfo, &vdir)

	// Create volume directory (/var/lib/glusterd/vols/<VOLNAME>)
	err := os.MkdirAll(vdir, os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	*path = fmt.Sprintf("%s/trusted-%s.tcp-fuse.vol", vdir, vinfo.Name)
	return err
}
