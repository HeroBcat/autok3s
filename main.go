package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {

	var (
		username          string
		password          string
		masterIPs         []string
		workerIPs         []string
		authKeysPath      string
		preMasterCommands []string
		preWorkerCommands []string
		masterCommands    []string
		workerCommands    []string
		masterExtraArgs   string
		workerExtraArgs   string
	)

	rootCmd := &cobra.Command{
		Use:   "autok3s",
		Short: "auto install docker.io k3s nfs-common",
		Run: func(cmd *cobra.Command, args []string) {

			if password == "" {
				log.Println("password can't be empty")
				os.Exit(1)
			}

			authContent := ""

			_, err := os.Stat(authKeysPath)
			if err == nil {
				result, err := ioutil.ReadFile(authKeysPath)
				if err == nil {
					authContent = string(result)
				}
			}

			var k3sToken string
			var k3sURL string
			var isHA = len(masterIPs) > 1

			for _, ip := range masterIPs {

				if authContent != "" {
					splits := strings.Split(authContent, "\n")
					for _, split := range splits {
						runCommands(ip, username, password, "echo "+split+" >> /home/"+username+"/.ssh/authorized_keys")
					}
				}

				runCommands(ip, username, password, preMasterCommands...)

				runCommands(ip, username, password, "sudo apt install nfs-common -y")
				k3sCmd := "curl -sfL http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -s - --tls-san " + ip
				if strings.Contains(masterExtraArgs, "--docker") {
					runCommands(ip, username, password, "sudo apt install docker.io -y")
				}
				if masterExtraArgs != "" {
					k3sCmd += " " + masterExtraArgs
				}

				if isHA {
					if k3sURL != "" && k3sToken != "" {
						k3sCmd += " --server " + k3sURL
						k3sCmd += " --token " + k3sToken
					} else {
						k3sCmd += "  --cluster-init"
					}
				}

				runCommands(ip, username, password, k3sCmd)

				if k3sURL == "" {
					k3sURL = "https://" + ip + ":6443"
				}

				if k3sToken == "" {
					k3sToken = getCommandsOutput(ip, username, password, "sudo cat -s /var/lib/rancher/k3s/server/node-token")
					k3sConfig := getCommandsOutput(ip, username, password, "sudo cat -s /etc/rancher/k3s/k3s.yaml")
					k3sConfig = strings.Replace(k3sConfig, "127.0.0.1", ip, 1)
					ioutil.WriteFile("k3s.yaml", []byte(k3sConfig), os.ModePerm)
				}

				runCommands(ip, username, password, masterCommands...)
			}

			for _, ip := range workerIPs {

				splits := strings.Split(authContent, "\n")
				for _, split := range splits {
					runCommands(ip, username, password, "echo "+split+" >> /home/"+username+"/.ssh/authorized_keys")
				}

				runCommands(ip, username, password, preWorkerCommands...)

				runCommands(ip, username, password, "sudo apt install nfs-common -y")
				k3sCmd := "curl -sfL http://rancher-mirror.cnrancher.com/k3s/k3s-install.sh | INSTALL_K3S_MIRROR=cn sh -s - agent"
				if strings.Contains(workerExtraArgs, "--docker") {
					runCommands(ip, username, password, "sudo apt install docker.io -y")
				}
				if workerExtraArgs != "" {
					k3sCmd += " " + workerExtraArgs
				}
				k3sCmd += " --server " + k3sURL + " --token " + k3sToken
				runCommands(ip, username, password, k3sCmd)

				runCommands(ip, username, password, workerCommands...)
			}

		},
	}

	rootCmd.Flags().StringVarP(&username, "user", "u", "ubuntu", "Specify the username of the server")
	rootCmd.Flags().StringVarP(&password, "password", "p", "", "Specify the password of the server")
	rootCmd.Flags().StringSliceVarP(&masterIPs, "master_ips", "m", nil, "Specify the ip of the master servers")
	rootCmd.Flags().StringSliceVarP(&workerIPs, "worker_ips", "w", nil, "Specify the ip of the worker servers")
	rootCmd.Flags().StringVarP(&authKeysPath, "auth_keys_path", "k", "", "Specify the path of id rsa content")
	rootCmd.Flags().StringSliceVar(&preMasterCommands, "pre_master_commands", nil, "Specify the commands that need to be pre-installed on the master server")
	rootCmd.Flags().StringSliceVar(&preWorkerCommands, "pre_worker_commands", nil, "Specify the commands that need to be pre-installed on the worker server")
	rootCmd.Flags().StringSliceVar(&masterCommands, "master_commands", nil, "Specify the commands that need to be installed on the master server")
	rootCmd.Flags().StringSliceVar(&workerCommands, "worker_commands", nil, "Specify the commands that need to be installed on the worker server")
	rootCmd.Flags().StringVar(&masterExtraArgs, "master_extra_args", "", "Specify the extra args on the master server")
	rootCmd.Flags().StringVar(&workerExtraArgs, "worker_extra_args", "", "Specify the extra args on the worker server")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
