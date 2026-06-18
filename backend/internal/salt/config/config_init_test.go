package config

import (
	"os"
	"testing"

	"github.com/keepsty/go_rds/internal/cluster/models"
)

func TestFileTemplate(t *testing.T) {
	// 创建临时目录作为测试工作目录
	tmpDir, err := os.MkdirTemp("", "salt-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建默认模板目录
	defaultDir := tmpDir + "/salt/config/default"
	stateSlsDir := tmpDir + "/salt/config/state_sls"
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(stateSlsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 写入测试模板文件
	tplContent := `[mysqld]
port=[[ .Port ]]
server_id=[[ .ServerId ]]
basedir=[[ .BaseDir ]]
datadir=[[ .Datadir ]]
`
	if err := os.WriteFile(defaultDir+"/mysql_80_cnf", []byte(tplContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 保存并恢复工作目录
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	tests := []struct {
		name    string
		e       interface{}
		files   []*models.SaltStateFiles
		wantErr bool
	}{
		{
			name: "mysql_80_cnf_template",
			e: &models.SaltMysqlHostPost{
				Port:      3306,
				ServerId:  1003306,
				BaseDir:   "/usr/local/mysql_8032",
				Datadir:   "/data/mysql_3306/data",
				Version:   "8032",
				MysqlIp:   "192.168.1.101",
				Host:      "node1",
				MysqlDir:  "/data/mysql_3306",
			},
			files: []*models.SaltStateFiles{
				{
					FilePath:       defaultDir,
					FileName:       "mysql_80_cnf",
					TargetFilePath: tmpDir + "/output/node1/mysql_3306",
					TargetFileName: "my_3306.cnf",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FileTemplate(tt.e, tt.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(tt.files) > 0 {
				// 验证输出文件已创建
				outPath := tt.files[0].TargetFilePath + "/" + tt.files[0].TargetFileName
				if _, err := os.Stat(outPath); os.IsNotExist(err) {
					t.Errorf("输出文件未创建: %s", outPath)
				}
			}
		})
	}
}
