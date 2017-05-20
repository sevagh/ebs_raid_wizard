package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EbsVol struct {
	EbsVolId     string
	VolumeName   string
	RaidLevel    int
	VolumeSize   int
	AttachedName string
	MountPath    string
	FsType       string
}

func MapEbsVolumes(ec2Instance *Ec2Instance) (map[string][]EbsVol, error) {
	drivesToMount := map[string][]EbsVol{}

	volumes, err := findEbsVolumes(ec2Instance)
	if err != nil {
		return drivesToMount, nil
	}

	log.Printf("Mapping EBS volumes")
	for _, volume := range volumes {
		drivesToMount[volume.VolumeName] = append(drivesToMount[volume.VolumeName], volume)
	}

	toDelete := []string{}

	for volName, volumes := range drivesToMount {
		//check if volName exists already
		if DoesLabelExist(PREFIX + "-" + volName) {
			log.Printf("Label already exists in /dev/disk/by-label")
			toDelete = append(toDelete, volName)
			continue
		}
		//check for volume mismatch
		volSize := volumes[0].VolumeSize
		mountPath := volumes[0].MountPath
		fsType := volumes[0].FsType
		raidLevel := volumes[0].RaidLevel
		if len(volumes) != volSize {
			return drivesToMount, fmt.Errorf("Found %d volumes, expected %d from VolumeSize tag", len(volumes), volSize)
		}
		for _, vol := range volumes[1:] {
			if volSize != vol.VolumeSize || mountPath != vol.MountPath || fsType != vol.FsType || raidLevel != vol.RaidLevel {
				return drivesToMount, fmt.Errorf("Mismatched tags among disks of same volume")
			}
		}
	}

	for _, volName := range(toDelete) {
		delete(drivesToMount, volName)
	}

	return drivesToMount, nil
}

func findEbsVolumes(ec2Instance *Ec2Instance) ([]EbsVol, error) {
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:" + PREFIX + "-IN:Prefix"),
				Values: []*string{
					aws.String(ec2Instance.Prefix),
				},
			},
			&ec2.Filter{
				Name: aws.String("tag:" + PREFIX + "-IN:NodeId"),
				Values: []*string{
					aws.String(ec2Instance.NodeId),
				},
			},
			&ec2.Filter{
				Name: aws.String("availability-zone"),
				Values: []*string{
					aws.String(ec2Instance.Az),
				},
			},
		},
	}

	volumes := []EbsVol{}

	result, err := ec2Instance.Ec2Client.DescribeVolumes(params)
	if err != nil {
		return volumes, err
	}

	for _, volume := range result.Volumes {
		ebsVolume := EbsVol{
			EbsVolId: *volume.VolumeId,
		}
		if len(volume.Attachments) > 0 {
			for _, attachment := range volume.Attachments {
				if *attachment.InstanceId != ec2Instance.InstanceId {
					return volumes, fmt.Errorf("Volume %s attached to different instance-id: %s", *volume.VolumeId, attachment.InstanceId)
				}
				ebsVolume.AttachedName = *attachment.Device
			}
		} else {
			ebsVolume.AttachedName = ""
		}
		for _, tag := range volume.Tags {
			switch *tag.Key {
			case PREFIX + "-IN:VolumeName":
				ebsVolume.VolumeName = *tag.Value
			case PREFIX + "-IN:RaidLevel":
				if ebsVolume.RaidLevel, err = strconv.Atoi(*tag.Value); err != nil {
					return volumes, fmt.Errorf("Couldn't parse RaidLevel tag as int: %v", err)
				}
			case PREFIX + "-IN:VolumeSize":
				if ebsVolume.VolumeSize, err = strconv.Atoi(*tag.Value); err != nil {
					return volumes, fmt.Errorf("Couldn't parse VolumeSize tag as int: %v", err)
				}
			case PREFIX + "-IN:MountPath":
				ebsVolume.MountPath = *tag.Value
			case PREFIX + "-IN:FsType":
				ebsVolume.FsType = *tag.Value
			case PREFIX + "-IN:NodeId": //do nothing
			case PREFIX + "-IN:Prefix": //do nothing
			default:
			}
		}
		volumes = append(volumes, ebsVolume)
	}
	return volumes, nil
}
