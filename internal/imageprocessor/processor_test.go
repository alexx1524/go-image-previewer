package imageprocessor

import (
	"bytes"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func ReadTestFile(fileName string) []byte {
	data, _ := os.ReadFile(fileName)
	return data
}

func Test_processor_Crop(t *testing.T) {
	type args struct {
		width  int
		height int
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Crop success",
			args: args{
				width:  500,
				height: 400,
				data:   ReadTestFile("../../tests/static/_gopher_original_1024x504.jpg"),
			},
			wantErr: false,
		},
		{
			name: "Source image is smaller than preview size",
			args: args{
				width:  2000,
				height: 400,
				data:   ReadTestFile("../../tests/static/_gopher_original_1024x504.jpg"),
			},
			wantErr: true,
		},
		{
			name: "Wrong type format",
			args: args{
				width:  500,
				height: 400,
				data:   ReadTestFile("../../tests/static/test.txt"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &processor{}

			got, err := p.Crop(tt.args.width, tt.args.height, tt.args.data)

			if (err != nil) != tt.wantErr {
				t.Errorf("Crop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				config, _, err := image.DecodeConfig(bytes.NewReader(got))

				require.NoError(t, err)
				require.Equal(t, tt.args.width, config.Width)
				require.Equal(t, tt.args.height, config.Height)
			}
		})
	}
}
