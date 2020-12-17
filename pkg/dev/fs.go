// This file is part of MinIO Direct CSI
// Copyright (c) 2020 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package dev

import (
	"encoding/binary"
	"errors"
	"math"
	"os"
)

var ErrNoFS = errors.New("No FS found")

type FSType string

type FSInfo struct {
	FSType        FSType  `json:"fsType,omitempty"`
	TotalCapacity uint64  `json:"totalCapacity,omitempty"`
	FreeCapacity  uint64  `json:"freeCapacity,omitempty"`
	FSBlockSize   uint64  `json:"fsBlockSize,omitempty"`
	Mounts        []Mount `json:"mounts,omitempty"`
}

func ProbeFS(devName string, logicalBlockSize uint64, offsetBlocks uint64) (*FSInfo, error) {
	ext4FSInfo, err := ProbeFSEXT4(devName, logicalBlockSize, offsetBlocks)
	if err != nil {
		if err != ErrNotEXT4 {
			return nil, err
		}
	}
	if ext4FSInfo != nil {
		return ext4FSInfo, nil
	}

	XFSInfo, err := ProbeFSXFS(devName, logicalBlockSize, offsetBlocks)
	if err != nil {
		if err != ErrNotXFS {
			return nil, err
		}
	}
	if XFSInfo != nil {
		return XFSInfo, nil
	}

	return nil, ErrNoFS
}

func ProbeFSEXT4(devName string, logicalBlockSize uint64, offsetBlocks uint64) (*FSInfo, error) {
	devPath := getBlockFile(devName)
	devFile, err := os.OpenFile(devPath, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return nil, err
	}
	defer devFile.Close()

	_, err = devFile.Seek(int64(logicalBlockSize*offsetBlocks), os.SEEK_CUR)
	if err != nil {
		return nil, err
	}

	ext4 := &EXT4SuperBlock{}
	err = binary.Read(devFile, binary.LittleEndian, ext4)
	if err != nil {
		return nil, err
	}
	if !ext4.Is() {
		return nil, ErrNotEXT4
	}

	fsBlockSize := uint64(math.Pow(2, float64(10+ext4.LogBlockSize)))
	fsInfo := &FSInfo{
		FSType:        FSTypeEXT4,
		FSBlockSize:   fsBlockSize,
		TotalCapacity: uint64(ext4.NumBlocks) * uint64(fsBlockSize),
		FreeCapacity:  uint64(ext4.FreeBlocks) * uint64(fsBlockSize),
		Mounts:        []Mount{},
	}

	return fsInfo, nil
}

func ProbeFSXFS(devName string, logicalBlockSize uint64, offsetBlocks uint64) (*FSInfo, error) {
	devPath := getBlockFile(devName)
	devFile, err := os.OpenFile(devPath, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return nil, err
	}
	defer devFile.Close()

	_, err = devFile.Seek(int64(logicalBlockSize*offsetBlocks), os.SEEK_CUR)
	if err != nil {
		return nil, err
	}

	xfs := &XFSSuperBlock{}
	err = binary.Read(devFile, binary.BigEndian, xfs)
	if err != nil {
		return nil, err
	}

	if !xfs.Is() {
		return nil, ErrNotXFS
	}

	fsInfo := &FSInfo{
		FSType:        FSTypeXFS,
		FSBlockSize:   uint64(xfs.BlockSize),
		TotalCapacity: uint64(xfs.TotalBlocks) * uint64(xfs.BlockSize),
		FreeCapacity:  uint64(xfs.FreeBlocks) * uint64(xfs.BlockSize),
		Mounts:        []Mount{},
	}

	return fsInfo, nil
}

const (
	None                uint32 = 0x0
	ADFS_SUPER_MAGIC           = 0xadf5
	AFFS_SUPER_MAGIC           = 0xadff
	AFS_SUPER_MAGIC            = 0x5346414f
	ANON_INODE_FS_MAGIC        = 0x09041934 /* Anonymous inode FS (for
	   pseudofiles that have no name;
	   e.g., epoll, signalfd, bpf) */
	AUTOFS_SUPER_MAGIC    = 0x0187
	BDEVFS_MAGIC          = 0x62646576
	BEFS_SUPER_MAGIC      = 0x42465331
	BFS_MAGIC             = 0x1badface
	BINFMTFS_MAGIC        = 0x42494e4d
	BPF_FS_MAGIC          = 0xcafe4a11
	BTRFS_SUPER_MAGIC     = 0x9123683e
	BTRFS_TEST_MAGIC      = 0x73727279
	CGROUP_SUPER_MAGIC    = 0x27e0eb   /* Cgroup pseudo FS */
	CGROUP2_SUPER_MAGIC   = 0x63677270 /* Cgroup v2 pseudo FS */
	CIFS_MAGIC_NUMBER     = 0xff534d42
	CODA_SUPER_MAGIC      = 0x73757245
	COH_SUPER_MAGIC       = 0x012ff7b7
	CRAMFS_MAGIC          = 0x28cd3d45
	DEBUGFS_MAGIC         = 0x64626720
	DEVFS_SUPER_MAGIC     = 0x1373 /* Linux 2.6.17 and earlier */
	DEVPTS_SUPER_MAGIC    = 0x1cd1
	ECRYPTFS_SUPER_MAGIC  = 0xf15f
	EFIVARFS_MAGIC        = 0xde5e81e4
	EFS_SUPER_MAGIC       = 0x00414a53
	EXT_SUPER_MAGIC       = 0x137d /* Linux 2.0 and earlier */
	EXT2_OLD_SUPER_MAGIC  = 0xef51
	EXT2_SUPER_MAGIC      = 0xef53
	EXT3_SUPER_MAGIC      = 0xef53
	EXT4_SUPER_MAGIC      = 0xef53
	F2FS_SUPER_MAGIC      = 0xf2f52010
	FUSE_SUPER_MAGIC      = 0x65735546
	FUTEXFS_SUPER_MAGIC   = 0xbad1dea /* Unused */
	HFS_SUPER_MAGIC       = 0x4244
	HOSTFS_SUPER_MAGIC    = 0x00c0ffee
	HPFS_SUPER_MAGIC      = 0xf995e849
	HUGETLBFS_MAGIC       = 0x958458f6
	ISOFS_SUPER_MAGIC     = 0x9660
	JFFS2_SUPER_MAGIC     = 0x72b6
	JFS_SUPER_MAGIC       = 0x3153464a
	MINIX_SUPER_MAGIC     = 0x137f     /* original minix FS */
	MINIX_SUPER_MAGIC2    = 0x138f     /* 30 char minix FS */
	MINIX2_SUPER_MAGIC    = 0x2468     /* minix V2 FS */
	MINIX2_SUPER_MAGIC2   = 0x2478     /* minix V2 FS, 30 char names */
	MINIX3_SUPER_MAGIC    = 0x4d5a     /* minix V3 FS, 60 char names */
	MQUEUE_MAGIC          = 0x19800202 /* POSIX message queue FS */
	MSDOS_SUPER_MAGIC     = 0x4d44
	MTD_INODE_FS_MAGIC    = 0x11307854
	NCP_SUPER_MAGIC       = 0x564c
	NFS_SUPER_MAGIC       = 0x6969
	NILFS_SUPER_MAGIC     = 0x3434
	NSFS_MAGIC            = 0x6e736673
	NTFS_SB_MAGIC         = 0x5346544e
	OCFS2_SUPER_MAGIC     = 0x7461636f
	OPENPROM_SUPER_MAGIC  = 0x9fa1
	OVERLAYFS_SUPER_MAGIC = 0x794c7630
	PIPEFS_MAGIC          = 0x50495045
	PROC_SUPER_MAGIC      = 0x9fa0 /* /proc FS */
	PSTOREFS_MAGIC        = 0x6165676c
	QNX4_SUPER_MAGIC      = 0x002f
	QNX6_SUPER_MAGIC      = 0x68191122
	RAMFS_MAGIC           = 0x858458f6
	REISERFS_SUPER_MAGIC  = 0x52654973
	ROMFS_MAGIC           = 0x7275
	SECURITYFS_MAGIC      = 0x73636673
	SELINUX_MAGIC         = 0xf97cff8c
	SMACK_MAGIC           = 0x43415d53
	SMB_SUPER_MAGIC       = 0x517b
	SMB2_MAGIC_NUMBER     = 0xfe534d42
	SOCKFS_MAGIC          = 0x534f434b
	SQUASHFS_MAGIC        = 0x73717368
	SYSFS_MAGIC           = 0x62656572
	SYSV2_SUPER_MAGIC     = 0x012ff7b6
	SYSV4_SUPER_MAGIC     = 0x012ff7b5
	TMPFS_MAGIC           = 0x01021994
	TRACEFS_MAGIC         = 0x74726163
	UDF_SUPER_MAGIC       = 0x15013346
	UFS_MAGIC             = 0x00011954
	USBDEVICE_SUPER_MAGIC = 0x9fa2
	V9FS_MAGIC            = 0x01021997
	VXFS_SUPER_MAGIC      = 0xa501fcf5
	XENFS_SUPER_MAGIC     = 0xabba1974
	XENIX_SUPER_MAGIC     = 0x012ff7b4
	XFS_SUPER_MAGIC       = 0x58465342
	_XIAFS_SUPER_MAGIC    = 0x012fd16d /* Linux 2.0 and earlier */
)
