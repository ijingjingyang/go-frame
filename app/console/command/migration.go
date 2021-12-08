package command

import (
	"fmt"
	"github.com/ijingjingyang/go-frame/client"
	"github.com/ijingjingyang/go-frame/conf"
	"github.com/jjonline/go-lib-backend/migrate"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "数据库迁移相关",
		Long:  `数据库迁移相关`,
	}
	migrateCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "查看迁移文件列表和状态",
		Long:  `查看迁移文件列表和状态`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getMigration().Status()
		},
	})
	migrateCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "创建迁移文件",
		Long:  `创建迁移文件`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getMigration().Create(args[0])
		},
		Args: cobra.ExactArgs(1),
	})
	migrateCmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "执行迁移文件",
		Long:  `执行迁移文件`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return getMigration().ExecUp()
		},
	})
	migrateCmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "迁移文件回滚：不指定文件名，默认回滚migrations表最后一条迁移记录",
		Long:  `迁移文件回滚：不指定文件名，默认回滚migrations表最后一条迁移记录`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := "" // 文件名需包含后缀
			if len(args) > 0 {
				filename = args[0]
			}
			return getMigration().ExecDown(filename)
		},
	})
	RootCmd.AddCommand(migrateCmd)
}

func getMigration() *migrate.Migrate {
	dbHandle, _ := client.DB.DB()
	return migrate.New(migrate.Config{
		Dir:       "migrations",
		TableName: fmt.Sprintf("%smigrations", conf.Config.Database.Prefix),
		DB:        dbHandle,
	})
}
