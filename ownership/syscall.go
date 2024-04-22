package ownership

import (
	"k8s.io/klog/v2"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func Ownership() {
	directory := "/home/arnob"
	fileInfo, err := os.Stat(directory)
	if err != nil {
		klog.Infoln("Error:", err)
		return
	}
	if !fileInfo.IsDir() {
		klog.Infoln("Not a directory:", directory)
		return
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		klog.Infoln("Error reading directory contents:", err)
		return
	}
	for _, entry := range entries {
		fileInfo, err := entry.Info()
		if err != nil {
			klog.Infof("Error getting file info for %s: %v\n", entry.Name(), err)
			continue
		}
		uid, gid := getOwnershipInfo(fileInfo)
		klog.Infof("name=%v , uid=%v , gid=%v\n", entry.Name(), uid, gid)
	}
}

func getOwnershipInfo(fileInfo os.FileInfo) (string, string) {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	uid := strconv.Itoa(int(stat.Uid))
	gid := strconv.Itoa(int(stat.Gid))
	return uid, gid
}

func printOwnershipInfo(name string, fileInfo os.FileInfo) {
	// Get file's Unix attributes
	stat := fileInfo.Sys().(*syscall.Stat_t)

	// Extract UID and GID from the file's Unix attributes
	uid := strconv.Itoa(int(stat.Uid))
	gid := strconv.Itoa(int(stat.Gid))

	klog.Infof("uid=%v , gid=%v\n", uid, gid)

	// Get username corresponding to UID
	user, err := user.LookupId(uid)
	if err != nil {
		klog.Infof("Error getting user for %s: %v\n", name, err)
		return
	}

	// Get group name corresponding to GID
	//group, err := user.LookupGroupId(gid)
	group, err := user.GroupIds()
	if err != nil {
		klog.Infof("Error getting group for %s: %v\n", name, err)
		return
	}

	// Print ownership information
	klog.Infof("Owner: %s\tGroup: %s\tFile: %s\n", user.Username, group, name)
}
