package main

import (
	"fmt"
	"os"

	"github.com/vmware/go-nfs-client/nfs"
	"github.com/vmware/go-nfs-client/nfs/rpc"
)

const MOUNT_DIRECTORY = "/mnt/nfs"

func WriteFile(path string, data []byte) error {
	mount, err := nfs.DialMount("nfs.example.com:/share")
	if err != nil {
		return fmt.Errorf("failed to connect to NFS server: %v", err)
	}
	defer mount.Close()

	target, err := mount.Mount(MOUNT_DIRECTORY, rpc.AuthNull)
	if err != nil {
		return fmt.Errorf("failed to mount NFS directory: %v", err)
	}
	defer func() {
		err := mount.Unmount()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to unmount: %v\n", err)
		}
	}()

	file, err := target.OpenFile(path, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func main() {
	err := WriteFile("test.txt", []byte("hello world"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
		os.Exit(1)
	}
}
