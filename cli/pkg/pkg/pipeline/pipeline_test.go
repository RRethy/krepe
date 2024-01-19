package pipeline

import (
	_ "github.com/stretchr/testify/assert"
	_ "testing"
)

// func TestPipelineRun(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		pipeline *Pipeline
// 		wantErr  bool
// 	}{
// 		{
// 			name: "succeeds with valid pipeline",
// 			pipeline: &Pipeline{
// 				Name: "test",
// 				Steps: []Step{
// 					{
// 						Name:   "test",
// 						Script: "echo 'hello world'",
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "fails with invalid pipeline",
// 			pipeline: &Pipeline{
// 				Name: "test",
// 				Steps: []Step{
// 					{
// 						Name:   "test",
// 						Script: "exit 1",
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "fails with invalid pipeline",
// 			pipeline: &Pipeline{
// 				Name: "test",
// 				Steps: []Step{
// 					{
// 						Name:   "test",
// 						Script: "exit 1",
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "fails with invalid pipeline",
// 			pipeline: &Pipeline{
// 				Name: "test",
// 				Steps: []Step{
// 					{
// 						Name:   "test",
// 						Script: "exit 1",
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := tt.pipeline.Run()
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }
