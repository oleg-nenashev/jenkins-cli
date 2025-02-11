package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/jenkins-zh/jenkins-cli/client"
	"github.com/jenkins-zh/jenkins-cli/util"
	"github.com/spf13/cobra"
)

// JobTypeOption is the job type cmd option
type JobTypeOption struct {
	OutputOption

	RoundTripper http.RoundTripper
}

var jobTypeOption JobTypeOption

func init() {
	jobCmd.AddCommand(jobTypeCmd)
	jobTypeCmd.Flags().StringVarP(&jobTypeOption.Format, "output", "o", "table", "Format the output")
}

var jobTypeCmd = &cobra.Command{
	Use:   "type",
	Short: "Print the types of job which in your Jenkins",
	Long:  `Print the types of job which in your Jenkins`,
	Run: func(cmd *cobra.Command, _ []string) {
		jclient := &client.JobClient{
			JenkinsCore: client.JenkinsCore{
				RoundTripper: centerUpgradeOption.RoundTripper,
			},
		}
		getCurrentJenkinsAndClient(&(jclient.JenkinsCore))

		if status, err := jclient.GetJobTypeCategories(); err == nil {
			var data []byte
			if data, err = jobTypeOption.Output(status); err == nil {
				if len(data) > 0 {
					cmd.Print(string(data))
				}
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	},
}

// Output renders data into a table
func (o *JobTypeOption) Output(obj interface{}) (data []byte, err error) {
	if data, err = o.OutputOption.Output(obj); err != nil {
		buf := new(bytes.Buffer)

		jobCategories := obj.([]client.JobCategory)
		table := util.CreateTable(buf)
		table.AddRow("number", "name", "type")
		for _, jobCategory := range jobCategories {
			for i, item := range jobCategory.Items {
				table.AddRow(fmt.Sprintf("%d", i), item.DisplayName,
					jobCategory.Name)
			}
		}
		table.Render()
		err = nil
		data = buf.Bytes()
	}
	return
}
