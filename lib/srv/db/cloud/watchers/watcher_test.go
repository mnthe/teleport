/*
Copyright 2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package watchers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/teleport/lib/srv/db/cloud"
	"github.com/gravitational/teleport/lib/srv/db/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/require"
)

// TestWatcher tests cloud databases watcher.
func TestWatcher(t *testing.T) {
	ctx := context.Background()

	rdsInstance1, rdsDatabase1 := makeRDSInstance(t, "instance-1", "us-east-1", map[string]string{"env": "prod"})
	rdsInstance2, _ := makeRDSInstance(t, "instance-2", "us-east-2", map[string]string{"env": "prod"})
	rdsInstance3, _ := makeRDSInstance(t, "instance-3", "us-east-1", map[string]string{"env": "dev"})

	auroraCluster1, auroraDatabase1 := makeRDSCluster(t, "cluster-1", "us-east-1", services.RDSEngineModeProvisioned, map[string]string{"env": "prod"})
	auroraCluster2, auroraDatabase2 := makeRDSCluster(t, "cluster-2", "us-east-2", services.RDSEngineModeProvisioned, map[string]string{"env": "dev"})
	auroraCluster3, _ := makeRDSCluster(t, "cluster-3", "us-east-2", services.RDSEngineModeProvisioned, map[string]string{"env": "prod"})
	auroraClusterUnsupported, _ := makeRDSCluster(t, "serverless", "us-east-1", services.RDSEngineModeServerless, map[string]string{"env": "prod"})

	watcher, err := NewWatcher(ctx, WatcherConfig{
		AWSMatchers: []services.AWSMatcher{
			{
				Types:   []string{services.AWSMatcherRDS},
				Regions: []string{"us-east-1"},
				Tags:    types.Labels{"env": []string{"prod"}},
			},
			{
				Types:   []string{services.AWSMatcherRDS},
				Regions: []string{"us-east-2"},
				Tags:    types.Labels{"env": []string{"dev"}},
			},
		},
		Clients: &common.TestCloudClients{
			RDSPerRegion: map[string]rdsiface.RDSAPI{
				"us-east-1": &cloud.RDSMock{
					DBInstances: []*rds.DBInstance{rdsInstance1, rdsInstance3},
					DBClusters:  []*rds.DBCluster{auroraCluster1, auroraClusterUnsupported},
				},
				"us-east-2": &cloud.RDSMock{
					DBInstances: []*rds.DBInstance{rdsInstance2},
					DBClusters:  []*rds.DBCluster{auroraCluster2, auroraCluster3},
				},
			},
		},
	})
	require.NoError(t, err)

	go watcher.fetchAndSend()
	select {
	case databases := <-watcher.DatabasesC():
		require.Equal(t, types.Databases{
			rdsDatabase1, auroraDatabase1, auroraDatabase2}, databases)
	case <-time.After(time.Second):
		t.Fatal("didn't receive databases after 1 second")
	}
}

func makeRDSInstance(t *testing.T, name, region string, labels map[string]string) (*rds.DBInstance, types.Database) {
	instance := &rds.DBInstance{
		DBInstanceArn:        aws.String(fmt.Sprintf("arn:aws:rds:%v:1234567890:db:%v", region, name)),
		DBInstanceIdentifier: aws.String(name),
		DbiResourceId:        aws.String(uuid.New()),
		Engine:               aws.String(services.RDSEnginePostgres),
		Endpoint: &rds.Endpoint{
			Address: aws.String("localhost"),
			Port:    aws.Int64(5432),
		},
		TagList: labelsToTags(labels),
	}
	database, err := services.NewDatabaseFromRDSInstance(instance)
	require.NoError(t, err)
	return instance, database
}

func makeRDSCluster(t *testing.T, name, region, engineMode string, labels map[string]string) (*rds.DBCluster, types.Database) {
	cluster := &rds.DBCluster{
		DBClusterArn:        aws.String(fmt.Sprintf("arn:aws:rds:%v:1234567890:cluster:%v", region, name)),
		DBClusterIdentifier: aws.String(name),
		DbClusterResourceId: aws.String(uuid.New()),
		Engine:              aws.String(services.RDSEngineAuroraMySQL),
		EngineMode:          aws.String(engineMode),
		Endpoint:            aws.String("localhost"),
		Port:                aws.Int64(3306),
		TagList:             labelsToTags(labels),
	}
	database, err := services.NewDatabaseFromRDSCluster(cluster)
	require.NoError(t, err)
	return cluster, database
}

func labelsToTags(labels map[string]string) (tags []*rds.Tag) {
	for key, val := range labels {
		tags = append(tags, &rds.Tag{
			Key:   aws.String(key),
			Value: aws.String(val),
		})
	}
	return tags
}
