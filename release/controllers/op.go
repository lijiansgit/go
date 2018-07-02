package controllers

import (
	"errors"
	"fmt"
	"release/libs"
	"release/models"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	// 检查配置文件中的int类型错误配置
	if _, err := beego.AppConfig.Int("preBackupNum"); err != nil {
		panic(err)
	}
	if _, err := beego.AppConfig.Int("pageSize"); err != nil {
		panic(err)
	}
}

type Release struct {
	localUsername string
	autoLogin     string
	privateKey    string
	preUsername   string
	prePassword   string
	preHost       string
	proHost       []string
	port          string
	localPath     string
	prePath       string
	preBackupPath string
	preBackupNum  int
	preFileTmp    string
	proLink       string
	proPathOne    string
	proPathTwo    string
	syncCmdFormat string
	syncExclude   string
	preSyncCmd    string
	sshPre        *libs.Ssh
	email         string
	//proSync       map[string]string
}

func NewRelease() *Release {
	localUsername := beego.AppConfig.String("localUsername")
	autoLogin := beego.AppConfig.String("autoLogin")
	privateKey := beego.AppConfig.String("privateKey")
	preUsername := beego.AppConfig.String("preUsername")
	prePassword := beego.AppConfig.String("prePassword")
	preHost := beego.AppConfig.String("preHost")
	port := beego.AppConfig.String("port")
	localPath := beego.AppConfig.String("localPath")
	prePath := beego.AppConfig.String("prePath")
	preBackupPath := beego.AppConfig.String("preBackupPath")
	preBackupNum, _ := beego.AppConfig.Int("preBackupNum")
	preFileTmp := beego.AppConfig.String("preFileTmp")
	proLink := beego.AppConfig.String("proLink")
	proPathOne := beego.AppConfig.String("proPathOne")
	proPathTwo := beego.AppConfig.String("proPathTwo")
	email := beego.AppConfig.String("Email")
	exclude := []string{
		"--exclude .git/",
		"--exclude /upload",
	}
	syncExclude := strings.Join(exclude, " ")
	var syncCmdFormat string
	if autoLogin == "yes" {
		syncCmdFormat = `rsync --out-format="%%n" -rltDz -H '-e ssh -i ` + privateKey +
			` -p %s' %s %s/ %s@%s:%s/`
	} else {
		syncCmdFormat = `rsync --out-format="%%n" -rltDz -H '-e ssh -p %s' %s %s/ %s@%s:%s/`
	}
	preSyncCmd := fmt.Sprintf(syncCmdFormat, port, syncExclude, localPath, preUsername, preHost, prePath)
	proHost := strings.Split(beego.AppConfig.String("proHost"), ",")
	//proSync := make(map[string]string)
	//for _, host := range proHost {
	//	proSync[host] = fmt.Sprintf(syncFormat, port, excludeStr, prePath, preUsername, host, prePath)
	//}
	sshPre := libs.NewSsh(autoLogin, privateKey, preUsername, prePassword, preHost, port)

	return &Release{
		localUsername: localUsername,
		autoLogin:     autoLogin,
		privateKey:    privateKey,
		preUsername:   preUsername,
		prePassword:   prePassword,
		preHost:       preHost,
		port:          port,
		localPath:     localPath,
		prePath:       prePath,
		preBackupPath: preBackupPath,
		preBackupNum:  preBackupNum,
		preFileTmp:    preFileTmp,
		proLink:       proLink,
		proPathOne:    proPathOne,
		proPathTwo:    proPathTwo,
		proHost:       proHost,
		syncCmdFormat: syncCmdFormat,
		syncExclude:   syncExclude,
		preSyncCmd:    preSyncCmd,
		sshPre:        sshPre,
		email:         email,
		//proSync:       proSync,
	}
}

func (r *Release) branchStatus() (status []string, err error) {
	res, err := libs.Cmd("git status", r.localPath)
	if err != nil {
		return status, err
	}

	status = strings.Split(res, "\n")
	return status, nil
}

func (r *Release) branchLog() (log []string, err error) {
	res, err := libs.Cmd("git log --stat -20", r.localPath)
	if err != nil {
		return log, err
	}

	logs := strings.Split(res, "\n")
	delimiter := "-------------------------------------------------------------------"
	for i := 0; i < len(logs); i++ {
		if strings.Contains(logs[i], "commit ") == true {
			log = append(log, delimiter)
		}
		log = append(log, logs[i])
	}
	return log, nil
}

func (r *Release) branchMaster() error {
	res, err := libs.Cmd("git branch", r.localPath)
	if err != nil {
		return err
	}

	if strings.Contains(res, "* master") == false {
		return errors.New("当前分支不是主分支！")
	}

	return nil
}

func (r *Release) branchPull() (pullRes []string, err error) {
	pullCmd := fmt.Sprintf(`su - %s -c 'cd %s && git pull origin master'`,
		r.localUsername, r.localPath)
	res, err := libs.Cmd(pullCmd)
	if err != nil {
		return pullRes, err
	}

	pullRes = strings.Split(res, "\n")
	return pullRes, nil
}

func (r *Release) Pre(qs orm.QuerySeter, order *models.ReleaseOrder) {
	beego.Info("[SYNC PRE START]")
	beego.Info("[SYNC PRE HOST]:", r.preHost)
	// 验证目录文件是否存在
	dirs := []string{r.preFileTmp, r.prePath, r.preBackupPath}
	for _, dir := range dirs {
		_, err := r.sshPre.Cmd(fmt.Sprintf("ls %s", dir))
		if err != nil {
			_, err = r.sshPre.Cmd(fmt.Sprintf("mkdir -p %s", dir))
			if err != nil {
				beego.Error("目录不存在，自动创建失败：", dir)
				return
			} else {
				beego.Info("目录不存在，自动创建成功：", dir)
			}
		}
	}

	beego.Info("[SYNC PRE CMD]:", r.preSyncCmd)
	fileLog, err := libs.Cmd(r.preSyncCmd)
	if err != nil {
		beego.Error("[SYNC PRE FAIL]:", err)
		return
	}

	fileLog = libs.StringTrim(fileLog, `"'`) //去除首尾单引号.双引号
	beego.Debug("[SYNC PRE RES]:", fileLog)

	// 给予www权限,非root用户要有sudo权限
	chownCmd := fmt.Sprintf("sudo chown -R %s.www %s && sudo chmod -R 775 %s",
		r.preUsername, r.prePath, r.prePath)
	beego.Info("[PRE CHOWN CMD]:", chownCmd)
	if _, err = r.sshPre.Cmd(chownCmd); err != nil {
		beego.Error("[PRECHOWN WWW FAIL]:", err)
		return
	}

	// 备份用于回滚
	backupDir := fmt.Sprintf("%s/%d", r.preBackupPath, order.Timestamp)
	mkdirCmd := fmt.Sprintf("mkdir %s", backupDir)
	copyCmd := fmt.Sprintf("cp -pr %s/* %s/", r.prePath, backupDir)
	beego.Info("[PREBACKUP MKDIR CMD]:", mkdirCmd)
	if _, err = r.sshPre.Cmd(mkdirCmd); err != nil {
		beego.Error("[PREBACKUP MKDIR FAIL]:", err)
		return
	}

	beego.Info("[PREBACKUP COPY CMD]:", copyCmd)
	if _, err = r.sshPre.Cmd(copyCmd); err != nil {
		beego.Error("[PREBACKUP COPY FAIL]:", err)
		return
	}

	// 删除大于preBackupNum的备份文件
	preBackup, err := r.sshPre.Cmd(fmt.Sprintf("ls %s/", r.preBackupPath))
	if err != nil {
		beego.Error(err)
		return
	}

	preBackupList := strings.Split(preBackup, "\n")
	if len(preBackupList) > r.preBackupNum {
		for _, preBackupDir := range preBackupList[:5] {
			rmCmd := fmt.Sprintf("rm -rf %s/%s", r.preBackupPath, preBackupDir)
			beego.Info("[SYNC PRE RM CMD]:", rmCmd)
			if _, err = r.sshPre.Cmd(rmCmd); err != nil {
				beego.Error(err)
				return
			}
		}
	}

	_, err = qs.Filter("Timestamp", order.Timestamp).Update(orm.Params{"Status": true, "FileLog": fileLog})
	if err != nil {
		beego.Error("更改工单状态失败: ", err)
		return
	}

	beego.Info("[SYNC PRE FINSH]")
}

func (r *Release) ProCmd(cmd, host string) (res string, err error) {
	var command string
	if r.autoLogin == "yes" {
		command = fmt.Sprintf(`ssh -i %s -p%s %s '%s'`,
			r.privateKey, r.port, host, cmd)
	} else {
		command = fmt.Sprintf(`ssh -p%s %s '%s'`, r.port, host, cmd)
	}
	res, err = r.sshPre.Cmd(command)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *Release) Pro(qs orm.QuerySeter, order *models.ReleaseOrder) {
	var (
		num      int
		fileLog  string
		syncPath string
		err      error
	)

	beego.Info("[SYNC PRO START]")
	for _, host := range r.proHost {
		num++
		beego.Info("[SYNC PRO HOST]:", num, host)
		// 根据当前软连接指向获取需要同步到哪个目录(_one,_two)
		proLinks, err := r.ProCmd(fmt.Sprintf(`ls -l %s`, r.proLink), host)
		if err != nil {
			beego.Error(err)
			return
		}

		linkDestPath := strings.TrimSpace(strings.Split(proLinks, "->")[1])
		if linkDestPath == r.proPathOne {
			syncPath = r.proPathTwo
		} else if linkDestPath == r.proPathTwo {
			syncPath = r.proPathOne
		}
		if syncPath != r.proPathOne && syncPath != r.proPathTwo {
			beego.Error("[SYNC PRO SYNCPATH IS ERROR]:", syncPath)
			return
		}

		// 执行同步
		proSyncCmd := fmt.Sprintf(r.syncCmdFormat, r.port, r.syncExclude, r.prePath, r.preUsername, host, syncPath)
		beego.Info("[SYNC PRO CMD]:", proSyncCmd)
		fileLog, err = r.sshPre.Cmd(proSyncCmd)
		if err != nil {
			beego.Error("[SYNC PRO FAIL]:", err)
			return
		}

		// 修改软连接到新的版本目录
		proNewLinkCmd := fmt.Sprintf(`rm -f %s && ln -s %s %s`,
			r.proLink, syncPath, r.proLink)
		beego.Info("[SYNC PRO NEW LINK CMD]:", proNewLinkCmd)
		if _, err = r.ProCmd(proNewLinkCmd, host); err != nil {
			beego.Error(err)
			return
		}

		// 两个工作目录里保持内容相同（_one,_two）
		proSyncCmdOneTwo := fmt.Sprintf(`rsync --out-format="%%n" -az -H %s %s/ %s/`,
			r.syncExclude, syncPath, linkDestPath)
		beego.Info("[SYNC PRO ONE TWO CMD]:", proSyncCmdOneTwo)
		proFileLog, err := r.ProCmd(proSyncCmdOneTwo, host)
		if err != nil {
			beego.Error(err)
			return
		}

		beego.Debug("[SYNC PRO ONE_TWO RES]:", libs.StringTrim(proFileLog, `"'`))
	}

	fileLog = libs.StringTrim(fileLog, `"'`) //去除首尾单引号.双引号
	beego.Debug("[SYNC PRO RES]:", fileLog)
	_, err = qs.Filter("Timestamp", order.Timestamp).Update(orm.Params{"Status": true, "FileLog": fileLog})
	if err != nil {
		beego.Error("更改工单状态失败:", err)
		return
	}

	//email
	messages := "发布结果<br/><br/><br/>"
	messages += "发布内容：" + order.Title + "<br/><br/>"
	messages += "发布文件：" + strings.Replace(fileLog, "\n", "<br/>", -1) + "<br/><br/>"
	messages += "发布时间：" + libs.TimestampToStr(order.Timestamp) + "<br/><br/>"
	messages += "发布人：" + order.OpName + "<br/><br/>"
	if r.email == "yes" {
		mail := libs.NewMail()
		err = mail.Send(messages)
		if err != nil {
			beego.Error("发送邮件失败：", err)
		} else {
			beego.Info("[SYNC PRO EMAIL]: 发送邮件成功")
		}
	}

	beego.Info("[SYNC PRO FINSH]")
}

func (r *Release) PreBack(backDir string, qs orm.QuerySeter, orderTimeStamp int) {
	beego.Info("[PREBACK SYNC START]")
	backSyncCmd := fmt.Sprintf(`rsync --out-format="%%n" -az -H %s/%s/ %s/`, r.preBackupPath, backDir, r.prePath)
	beego.Info("[PREBACK SYNC CMD]:", backSyncCmd)
	fileLog, err := r.sshPre.Cmd(backSyncCmd)
	if err != nil {
		beego.Error("[PREBACK SYNC FAIL]:", err)
		return
	}

	fileLog = libs.StringTrim(fileLog, `"'`) //去除首尾单引号.双引号
	beego.Debug("[PREBACK SYNC RES]:", fileLog)
	beego.Info("[PREBACK SYNC DIR]:", backDir)
	_, err = qs.Filter("Timestamp", orderTimeStamp).Update(orm.Params{"Status": true, "FileLog": fileLog})
	if err != nil {
		beego.Error("回滚后更改工单状态失败:", err)
		return
	}

	beego.Info("[PREBACK SYNC FINSH]")
}

func (r *Release) File(localFileTmp, preFileTmp, proFile string) (err error) {
	cpFormat := "/usr/bin/scp -P %s %s %s@%s:%s"
	preScpCmd := fmt.Sprintf(cpFormat, r.port, localFileTmp, r.preUsername, r.preHost, preFileTmp)
	beego.Info("[PRE SCP FILE]:", preScpCmd)
	if _, err = libs.Cmd(preScpCmd); err != nil {
		beego.Error(err)
		return err
	}

	preCpCmd := fmt.Sprintf("/usr/bin/cp %s %s", preFileTmp, proFile)
	beego.Info("[PRE CP FILE]:", preCpCmd)
	if _, err = r.sshPre.Cmd(preCpCmd); err != nil {
		beego.Error(err)
		return err
	}

	for _, host := range r.proHost {
		proScpCmd := fmt.Sprintf(cpFormat, r.port, preFileTmp, r.preUsername, host, proFile)
		beego.Info("[PRO SCP FILE]:", proScpCmd)
		if _, err = r.sshPre.Cmd(proScpCmd); err != nil {
			beego.Error(err)
			return err
		}
	}

	return nil
}
